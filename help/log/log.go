package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var (
	// 打印所有级别的日志
	lowPriority = zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.DebugLevel
	})

	// High-priority output should also go to standard error, and low-priority
	// output should also go to standard out.
	consoleDebugging = zapcore.Lock(os.Stdout)

	consoleEncoder = zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig())

	// Join the outputs, encoders, and level-handling functions into
	// zapcore.Cores, then tee the four cores together.
	core = zapcore.NewTee(
		// 打印在控制台
		zapcore.NewCore(consoleEncoder, consoleDebugging, lowPriority),
	)

	Logger = zap.New(core)
)

func Info(msg string, fields ...zapcore.Field) {
	Logger.Info(msg, fields...)
}
func Fatal(msg string, fields ...zapcore.Field) {
	Logger.Fatal(msg, fields...)
}
func Error(msg string, fields ...zapcore.Field) {
	Logger.Error(msg, fields...)
}
