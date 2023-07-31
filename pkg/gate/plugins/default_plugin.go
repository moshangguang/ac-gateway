package plugins

import "github.com/valyala/fasthttp"

type DefaultPlugin struct {
}

func (*DefaultPlugin) RequestFilter(_ *fasthttp.RequestCtx, _ string, _ PluginChain) {

}

func (*DefaultPlugin) ResponseFilter(_ *fasthttp.RequestCtx, _ string, _ PluginChain) {

}

var _ Plugin = &DefaultPlugin{}
