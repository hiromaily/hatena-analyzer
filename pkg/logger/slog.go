package logger

import (
	"log/slog"
	"os"
)

type SlogLogger struct {
	log  *slog.Logger
	args []any
}

func NewSlogLogger(
	level slog.Level,
	hostname string,
	appCode string,
	commitID string,
) *SlogLogger {
	args := []any{
		slog.String("hostname", hostname),
		slog.String("appCode", appCode),
		slog.String("commitID", commitID),
	}
	// logger option
	options := &slog.HandlerOptions{Level: level}

	return &SlogLogger{
		log:  slog.New(slog.NewJSONHandler(os.Stdout, options)),
		args: args,
	}
}

// Debug
func (s *SlogLogger) Debug(msg string, args ...any) {
	s.log.Debug(msg, s.appendArgs(args...)...)
}

// Info
func (s *SlogLogger) Info(msg string, args ...any) {
	s.log.Info(msg, s.appendArgs(args...)...)
}

// Warn
func (s *SlogLogger) Warn(msg string, args ...any) {
	s.log.Warn(msg, s.appendArgs(args...)...)
}

// Error
func (s *SlogLogger) Error(msg string, args ...any) {
	s.log.Error(msg, s.appendArgs(args...)...)
}

// appends the args to the default args
func (s *SlogLogger) appendArgs(args ...any) []any {
	return append(s.args, args...)
}
