package gosafe

import (
	"ac-gateway/help/log"
	"ac-gateway/help/panicutils"
	"go.uber.org/zap"
)

func GoSafe(fn func()) {
	go func() {
		defer func() {
			if p := recover(); p != nil {
				panicStack := panicutils.GetPanicStack()
				log.Error("出现异常", zap.Any("p", p), zap.String("panic_stack", panicStack))
			}
		}()
		fn()
	}()
}
