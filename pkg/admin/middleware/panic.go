package middleware

import (
	"ac-gateway/help/log"
	"ac-gateway/help/panicutils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func PanicHandle(c *gin.Context) {
	defer func() {
		v := recover()
		if v == nil {
			return
		}
		panicStack := panicutils.GetPanicStack()
		log.Logger.Error("后台路由出现panic", zap.Any("v", v), zap.String("panic_stack", panicStack))
	}()
	c.Next()
}
