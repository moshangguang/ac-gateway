package router

import (
	"ac-gateway/constant"
	"ac-gateway/help/contextutils"
	"ac-gateway/help/strutils"
	"ac-gateway/pkg/gate/plugins"
	"ac-gateway/pkg/models/ddl"
	"context"
	_ "github.com/juju/errors"
	"github.com/valyala/fasthttp"
	"github.com/vibrantbyte/go-antpath/antpath"
	"sync"
)

var chainPool = sync.Pool{
	New: func() any {
		return new(pluginChain)
	},
}

type RouterMatcher struct {
	methodRouter []MethodRouter
}
type MatchValue struct {
	lb         LoadBalancer
	prePlugin  []pluginConf
	postPlugin []pluginConf
}

func (router *RouterMatcher) Match(ctx context.Context) (RouterHandler, error) {
	path, _ := contextutils.GetStringValue(ctx, constant.CtxRequestPath)
	if strutils.IsEmpty(path) {
		return RouterHandler{}, constant.ErrCtxNotFoundRequestPath
	}
	method, _ := contextutils.GetStringValue(ctx, constant.CtxRequestMethod)
	if strutils.IsEmpty(method) {
		return RouterHandler{}, constant.ErrCtxNotFoundRequestMethod
	}
	host, _ := contextutils.GetStringValue(ctx, constant.CtxHeaderHost)
	for _, m := range router.methodRouter {
		if m.method == method {
			r, ok := m.match(host, path)
			if ok {
				return r, nil
			}
		}
	}
	return RouterHandler{}, constant.ErrCtxNotFoundMatchRouter
}

func (router *RouterMatcher) initMethod(method string) int {
	for i, m := range router.methodRouter {
		if m.method == method {
			return i
		}
	}
	router.methodRouter = append(router.methodRouter, MethodRouter{
		method: method,
		host:   make([]HostRouter, 0),
	})
	return len(router.methodRouter) - 1
}
func (router *RouterMatcher) addRoute(r ddl.Router) {
	if router.methodRouter == nil {
		router.methodRouter = make([]MethodRouter, 0)
	}
	for _, m := range r.Methods {
		index := router.initMethod(m)
		router.methodRouter[index] = router.methodRouter[index].addRoute(r)
	}
}

func NewRouter(routerList []ddl.Router) *RouterMatcher {
	router := new(RouterMatcher)
	for _, r := range routerList {
		router.addRoute(r)
	}
	return router
}

type MethodRouter struct {
	method string
	host   []HostRouter
}

func (methodRouter MethodRouter) match(host string, path string) (RouterHandler, bool) {
	if strutils.IsEmpty(host) {
		for _, h := range methodRouter.host {
			if h.host == host {
				for _, r := range h.router {
					match := r.antPathMatcher.MatchStrings(path, nil)
					if match {
						return r, true
					}
				}
			}
		}
	} else {
		for _, h := range methodRouter.host {
			for _, r := range h.router {
				match := r.antPathMatcher.MatchStrings(path, nil)
				if match {
					return r, true
				}
			}
		}
	}
	return RouterHandler{}, false
}
func (methodRouter MethodRouter) addRoute(router ddl.Router) MethodRouter {
	if methodRouter.host == nil {
		methodRouter.host = make([]HostRouter, 0)
	}
	index := -1
	for i, host := range methodRouter.host {
		if router.Host == host.host {
			index = i
			break
		}
	}
	if index == -1 {
		methodRouter.host = append(methodRouter.host, HostRouter{
			host:   router.Host,
			router: make([]RouterHandler, 0),
		})
		index = len(methodRouter.host) - 1
	}
	methodRouter.host[index] = methodRouter.host[index].addRoute(router)
	return methodRouter
}

type HostRouter struct {
	host   string
	router []RouterHandler
}

func (hostRouter HostRouter) addRoute(router ddl.Router) HostRouter {
	if hostRouter.router == nil {
		hostRouter.router = make([]RouterHandler, 0, 1)
	}
	for _, r := range hostRouter.router {
		if r.uri == router.Uri {
			return hostRouter
		}
	}

	rh := RouterHandler{
		uri:            router.Uri,
		antPathMatcher: antpath.NewMatchesStringMatcher(router.Uri, true),
	}
	if router.Upstream.Type == constant.RouteUpstreamTypeIPHash {
		nodes := make([]string, 0, len(router.Upstream.Nodes))
		for s := range router.Upstream.Nodes {
			nodes = append(nodes, s)
		}
		rh.lb = NewIPHashLoadBalancer(nodes)

	} else {
		nodes := make([]string, 0, len(router.Upstream.Nodes))
		for s, n := range router.Upstream.Nodes {
			for i := 0; i < n; i++ {
				nodes = append(nodes, s)
			}
		}
		rh.lb = NewRoundRobinLoadBalancer(nodes)
	}
	if len(router.Plugins.ExtPluginPreReq.Conf) != 0 {
		rh.prePlugin = make([]pluginConf, 0, len(router.Plugins.ExtPluginPreReq.Conf))
		for _, c := range router.Plugins.ExtPluginPreReq.Conf {
			p, ok := plugins.GetPlugin(c.Name)
			if ok {
				rh.prePlugin = append(rh.prePlugin, pluginConf{
					filter: p.RequestFilter,
					conf:   c.Value,
				})
			}
		}
	}
	if len(router.Plugins.ExtPluginPostResp.Conf) != 0 {
		rh.postPlugin = make([]pluginConf, 0, len(router.Plugins.ExtPluginPostResp.Conf))
		for _, c := range router.Plugins.ExtPluginPostResp.Conf {
			p, ok := plugins.GetPlugin(c.Name)
			if ok {
				rh.postPlugin = append(rh.postPlugin, pluginConf{
					filter: p.ResponseFilter,
					conf:   c.Value,
				})
			}
		}
	}
	hostRouter.router = append(hostRouter.router, rh)
	return hostRouter
}

type RouterHandler struct {
	uri            string
	antPathMatcher *antpath.AntPathStringMatcher
	lb             LoadBalancer
	prePlugin      []pluginConf
	postPlugin     []pluginConf
}

func (routerHandler RouterHandler) RequestFilter(ctx *fasthttp.RequestCtx) bool {
	if len(routerHandler.prePlugin) == 0 {
		return true
	}
	chain := chainPool.Get().(*pluginChain)
	defer chainPool.Put(chain)
	defer chain.Reset()
	chain.Reset().
		applyContext(ctx).
		applyPlugin(routerHandler.prePlugin).
		Next()
	return !chain.Abort()
}

func (routerHandler RouterHandler) ResponseFilter(ctx *fasthttp.RequestCtx) {
	if len(routerHandler.postPlugin) == 0 {
		return
	}
	chain := chainPool.Get().(*pluginChain)
	defer chainPool.Put(chain)
	defer chain.Reset()
	chain.Reset().
		applyContext(ctx).
		applyPlugin(routerHandler.postPlugin).
		Next()
}

func (routerHandler RouterHandler) GetNode(ctx context.Context) (string, error) {
	return routerHandler.lb.GetNode(ctx)
}

type pluginConf struct {
	filter plugins.PluginFilter
	conf   string
}

func CopyPluginConfList(list []pluginConf) []pluginConf {
	result := make([]pluginConf, 0, len(list))
	for _, item := range list {
		result = append(result, pluginConf{
			filter: item.filter,
			conf:   item.conf,
		})
	}
	return result
}

type Upstream struct {
	typ      string
	nodes    map[string]int
	nodeList []string
	uri      string
	lb       LoadBalancer
}
