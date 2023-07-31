package plugins

import (
	"github.com/spf13/cast"
	"github.com/valyala/fasthttp"
)

type StripPrefixPlugin struct {
	DefaultPlugin
}

func (plugin *StripPrefixPlugin) StripPrefix(ctx *fasthttp.RequestCtx, conf string) {
	stripPrefix := cast.ToInt(conf)
	if stripPrefix == 0 {
		return
	}
	index := -1
	path := string(ctx.Path())
	for i, c := range path {
		if i == 0 && c == '/' {
			continue
		}
		if c == '/' {
			stripPrefix--
		}
		if stripPrefix == 0 {
			index = i
			break
		}
	}
	if index == -1 {
		return
	}
	ctx.Request.URI().SetPath(path[index:])
}
func (plugin *StripPrefixPlugin) RequestFilter(ctx *fasthttp.RequestCtx, conf string, chain PluginChain) {
	plugin.StripPrefix(ctx, conf)
	chain.Next()
}

func init() {
	RegisterPlugin("strip_prefix", &StripPrefixPlugin{})
}
