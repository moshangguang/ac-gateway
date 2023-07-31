package plugins

import (
	"ac-gateway/help/strutils"
	"github.com/valyala/fasthttp"
)

type PluginChain interface {
	Next()
}
type Plugin interface {
	RequestFilter(ctx *fasthttp.RequestCtx, conf string, chain PluginChain)
	ResponseFilter(ctx *fasthttp.RequestCtx, conf string, chain PluginChain)
}

type PluginFilter func(ctx *fasthttp.RequestCtx, conf string, chain PluginChain)

var pluginMap = make(map[string]Plugin)

func RegisterPlugin(name string, plugin Plugin) {
	if plugin == nil {
		panic("plugin not allow nil")
	}
	if strutils.IsEmpty(name) {
		panic("plugin name not allow empty")
	}
	if _, ok := pluginMap[name]; ok {
		panic("plugin name not allow repeat")
	}
	pluginMap[name] = plugin
}

func GetPlugin(name string) (Plugin, bool) {
	plugin, ok := pluginMap[name]
	return plugin, ok
}
