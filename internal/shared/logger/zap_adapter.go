package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ZapLogger wraps zap.Logger to implement our Logger interface
type ZapLogger struct {
	logger *zap.Logger
}

// NewZapLogger creates a new Zap-based logger based on environment
func NewZapLogger(env string) (Logger, error) {
	var zapLogger *zap.Logger
	var err error

	switch env {
	case "production", "prod":
		// Production: JSON, info level, sampling
		config := zap.NewProductionConfig()
		config.EncoderConfig.TimeKey = "timestamp"
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		zapLogger, err = config.Build()

	case "development", "dev":
		// Development: Console, debug level, colorized
		config := zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		zapLogger, err = config.Build()

	default:
		// Default to development
		zapLogger, err = zap.NewDevelopment()
	}

	if err != nil {
		return nil, err
	}

	return &ZapLogger{logger: zapLogger}, nil
}

func NewZapLoggerFromZap(zapLogger *zap.Logger) Logger {
	return &ZapLogger{logger: zapLogger}
}

func (l *ZapLogger) Debug(msg string, fields ...Field) {
	l.logger.Debug(msg, l.convertFields(fields)...)
}

func (l *ZapLogger) Info(msg string, fields ...Field) {
	l.logger.Info(msg, l.convertFields(fields)...)
}

func (l *ZapLogger) Warn(msg string, fields ...Field) {
	l.logger.Warn(msg, l.convertFields(fields)...)
}

func (l *ZapLogger) Error(msg string, fields ...Field) {
	l.logger.Error(msg, l.convertFields(fields)...)
}

// Fatal logs a fatal message and exits
func (l *ZapLogger) Fatal(msg string, fields ...Field) {
	l.logger.Fatal(msg, l.convertFields(fields)...)
}

// With returns a new logger with the given fields always attached
func (l *ZapLogger) With(fields ...Field) Logger {
	return &ZapLogger{
		logger: l.logger.With(l.convertFields(fields)...),
	}
}

// Sync flushes any buffered log entries (call on shutdown)
func (l *ZapLogger) Sync() error {
	return l.logger.Sync()
}

// convertFields converts our Field type to zap.Field
func (l *ZapLogger) convertFields(fields []Field) []zap.Field {
	zapFields := make([]zap.Field, len(fields))
	for i, f := range fields {
		zapFields[i] = zap.Any(f.Key, f.Value)
	}
	return zapFields
}
