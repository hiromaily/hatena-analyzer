package logger

//
// Noop Logger
//

type NoopLogger struct{}

func NewNoopLogger() *NoopLogger {
	return &NoopLogger{}
}

// Debug
//
//nolint:revive
func (*NoopLogger) Debug(msg string, args ...any) {
	// DummyLogger disables logging
}

// Info
//
//nolint:revive
func (*NoopLogger) Info(msg string, args ...any) {
	// DummyLogger disables logging
}

// Warn
//
//nolint:revive
func (*NoopLogger) Warn(msg string, args ...any) {
	// DummyLogger disables logging
}

// Error
//
//nolint:revive
func (*NoopLogger) Error(msg string, args ...any) {
	// DummyLogger disables logging
}
