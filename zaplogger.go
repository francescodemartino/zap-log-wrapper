package log_wrapper

import "go.uber.org/zap"

type ZapLogger struct {
	logger *zap.Logger
}

func (z *ZapLogger) Info(msg string, fields ...zap.Field) {
	z.logger.Info(msg, fields...)
}

func (z *ZapLogger) Error(msg string, fields ...zap.Field) {
	z.logger.Error(msg, fields...)
}
