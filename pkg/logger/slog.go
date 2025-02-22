package logger

import (
	"log/slog"
	"os"

	"github.com/phsym/console-slog"
)

//
// SlogJSONLogger
//

type SlogJSONLogger struct {
	log  *slog.Logger
	args []any
}

func NewSlogJSONLogger(
	level slog.Level,
	appCode string,
	commitID string,
) *SlogJSONLogger {
	args := []any{
		slog.String("appCode", appCode),
	}
	if commitID != "" {
		args = append(args, slog.String("commitID", commitID))
	}

	// logger option
	options := &slog.HandlerOptions{Level: level}

	return &SlogJSONLogger{
		log:  slog.New(slog.NewJSONHandler(os.Stdout, options)),
		args: args,
	}
}

// Debug
func (s *SlogJSONLogger) Debug(msg string, args ...any) {
	s.log.Debug(msg, s.appendArgs(args...)...)
}

// Info
func (s *SlogJSONLogger) Info(msg string, args ...any) {
	s.log.Info(msg, s.appendArgs(args...)...)
}

// Warn
func (s *SlogJSONLogger) Warn(msg string, args ...any) {
	s.log.Warn(msg, s.appendArgs(args...)...)
}

// Error
func (s *SlogJSONLogger) Error(msg string, args ...any) {
	s.log.Error(msg, s.appendArgs(args...)...)
}

// appends the args to the default args
func (s *SlogJSONLogger) appendArgs(args ...any) []any {
	return append(s.args, args...)
}

//
// SlogConsoleLogger
// use https://github.com/phsym/console-slog
//

type SlogConsoleLogger struct {
	log  *slog.Logger
	args []any
}

func NewSlogConsoleLogger(
	level slog.Level,
) *SlogConsoleLogger {
	options := &console.HandlerOptions{Level: level}
	return &SlogConsoleLogger{
		log:  slog.New(console.NewHandler(os.Stderr, options)),
		args: []any{},
	}
}

// Debug
func (s *SlogConsoleLogger) Debug(msg string, args ...any) {
	s.log.Debug(msg, args...)
}

// Info
func (s *SlogConsoleLogger) Info(msg string, args ...any) {
	s.log.Info(msg, args...)
}

// Warn
func (s *SlogConsoleLogger) Warn(msg string, args ...any) {
	s.log.Warn(msg, args...)
}

// Error
func (s *SlogConsoleLogger) Error(msg string, args ...any) {
	s.log.Error(msg, args...)
}
