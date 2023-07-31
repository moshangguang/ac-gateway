package plugins

import "github.com/valyala/fasthttp"

type GatewayPlugin struct {
	DefaultPlugin
}

func (*GatewayPlugin) ResponseFilter(ctx *fasthttp.RequestCtx, _ string, chain PluginChain) {
	ctx.Response.Header.Add("proxy-server", "ac-gateway")
	chain.Next()
}

func init() {
	RegisterPlugin("gateway", &GatewayPlugin{})
}
