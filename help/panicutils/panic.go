package panicutils

import (
	"fmt"
	"runtime"
)

func GetPanicStack() string {
	buf := make([]byte, 4096)
	n := runtime.Stack(buf, false)
	return fmt.Sprintf("%s", string(buf[:n]))
}

func Must(err error) {
	if err != nil {
		panic(err)
	}
}
