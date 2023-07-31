package plugins

import (
	"ac-gateway/constant"
	"ac-gateway/help/contextutils"
	"ac-gateway/help/strutils"
	"github.com/valyala/fasthttp"
)

type PostNodePlugin struct {
	DefaultPlugin
}

func (p *PostNodePlugin) ResponseFilter(ctx *fasthttp.RequestCtx, _ string, chain PluginChain) {
	node, _ := contextutils.GetStringValue(ctx, constant.CtxRequestNode)
	if strutils.IsNotEmpty(node) {
		ctx.Response.Header.Add("ac-gateway-node", node)
	}
	chain.Next()
}

func init() {
	RegisterPlugin("post_node", &PostNodePlugin{})
}
