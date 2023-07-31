package fasthttputils

import (
	"ac-gateway/help/strutils"
	"github.com/valyala/fasthttp"
	"strings"
)

func GetHeaderHost(ctx *fasthttp.RequestCtx) (string, bool) {
	host := string(ctx.Request.Header.Host())
	return strings.TrimSpace(host), strutils.IsNotEmpty(host)
}
func GetRequestPath(ctx *fasthttp.RequestCtx) string {
	path := string(ctx.URI().Path())
	return path
}
