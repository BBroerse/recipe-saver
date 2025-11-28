package logger

// NoopLogger is a logger that does nothing
// Perfect for testing where you don't want log output
type NoopLogger struct{}

// NewNoopLogger creates a new no-op logger
func NewNoopLogger() Logger {
	return &NoopLogger{}
}

func (l *NoopLogger) Debug(msg string, fields ...Field) {}
func (l *NoopLogger) Info(msg string, fields ...Field)  {}
func (l *NoopLogger) Warn(msg string, fields ...Field)  {}
func (l *NoopLogger) Error(msg string, fields ...Field) {}
func (l *NoopLogger) Fatal(msg string, fields ...Field) {}

func (l *NoopLogger) With(fields ...Field) Logger {
	return l
}
