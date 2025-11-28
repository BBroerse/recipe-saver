package logger

// Logger is the application's logging interface
// This abstracts the underlying logging implementation (zap, slog, etc.)
type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)

	// With returns a new logger with the given fields always attached
	With(fields ...Field) Logger
}

// Field represents a structured logging field
type Field struct {
	Key   string
	Value interface{}
}

// Helper functions for creating fields
func String(key, value string) Field {
	return Field{Key: key, Value: value}
}

func Int(key string, value int) Field {
	return Field{Key: key, Value: value}
}

func Error(err error) Field {
	return Field{Key: "error", Value: err}
}

func Duration(key string, value interface{}) Field {
	return Field{Key: key, Value: value}
}

func Any(key string, value interface{}) Field {
	return Field{Key: key, Value: value}
}
