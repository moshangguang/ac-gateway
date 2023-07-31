package ginutils

import (
	"ac-gateway/constant"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RespSuccess(c *gin.Context) {
	RespSuccessWithData(c, struct{}{})
}
func RespSuccessWithData(c *gin.Context, data interface{}) {
	if c == nil {
		return
	}
	c.JSON(http.StatusOK, gin.H{
		constant.ResponseCode:    constant.ResponseSuccessCode,
		constant.ResponseMessage: constant.ResponseMsgSuccess,
		constant.ResponseData:    data,
	})
}
func RespFail(c *gin.Context, httpStatus, code int, msg string) {
	if c == nil {
		return
	}

	c.JSON(httpStatus, gin.H{
		constant.ResponseCode:    code,
		constant.ResponseMessage: msg,
	})
}
func RespError(c *gin.Context, err error) {
	if err == nil {
		return
	}
	RespFail(c, http.StatusBadRequest, constant.ResponseErrorCode, err.Error())
}
