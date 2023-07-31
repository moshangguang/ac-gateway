package fasthttputils

import (
	"ac-gateway/constant"
	"encoding/json"
	"github.com/valyala/fasthttp"
	"net/http"
)

func RespErrorWithStatusCode(ctx *fasthttp.RequestCtx, httpStatus, code int, err error) {
	if ctx == nil {
		return
	}
	if err == nil {
		return
	}
	ctx.SetStatusCode(httpStatus)
	data, _ := json.Marshal(map[string]interface{}{
		constant.ResponseMessage: err.Error(),
		constant.ResponseCode:    code,
	})
	ctx.Response.Header.Add("Content-Type", "application/json; charset=utf-8")
	_, err = ctx.Write(data)
}
func RespError(ctx *fasthttp.RequestCtx, err error) {
	RespErrorWithStatusCode(ctx, http.StatusBadRequest, constant.ErrGatewayCode, err)
}

func HasChangeResponse(response *fasthttp.Response) bool {
	if response == nil {
		return false
	}
	return !(len(response.Body()) == 0 && response.StatusCode() == 0)
}
