package router

import (
	"ac-gateway/pkg/gate/plugins"
	"github.com/valyala/fasthttp"
)

type pluginChain struct {
	ctx        *fasthttp.RequestCtx
	index      int
	pluginConf []pluginConf
	next       func()
}

func (chain *pluginChain) applyContext(ctx *fasthttp.RequestCtx) *pluginChain {
	chain.ctx = ctx
	return chain
}
func (chain *pluginChain) applyNext(next func()) *pluginChain {
	chain.next = next
	return chain
}
func (chain *pluginChain) applyPlugin(pluginConf []pluginConf) *pluginChain {
	chain.pluginConf = pluginConf
	return chain
}

func (chain *pluginChain) Next() {
	chain.index++
	if chain.index < len(chain.pluginConf) {
		conf := chain.pluginConf[chain.index]
		conf.filter(chain.ctx, conf.conf, chain)
	}
	chain.next()
}

func (chain *pluginChain) Reset() *pluginChain {
	chain.ctx = nil
	chain.index = -1
	chain.pluginConf = nil
	chain.next = emptyNext
	return chain
}
func (chain *pluginChain) Abort() bool {
	return len(chain.pluginConf) < chain.index
}

var _ plugins.PluginChain = &pluginChain{}
var emptyNext = func() {}
