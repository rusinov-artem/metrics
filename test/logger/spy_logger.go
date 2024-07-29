package logger

import (
	"bytes"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func SpyLogger() (*zap.Logger, *bytes.Buffer) {
	logger, _ := zap.NewDevelopment()
	logBuffer := bytes.NewBufferString("")
	logger = logger.WithOptions(zap.WrapCore(func(core zapcore.Core) zapcore.Core {
		return zapcore.NewCore(
			zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig()),
			zapcore.AddSync(logBuffer),
			zapcore.DebugLevel,
		)
	}))
	return logger, logBuffer
}
