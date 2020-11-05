package logger

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	tlecontext "github.com/kolbis/corego/context"
)

// ErrorWithContext middleware function print error log with context
func ErrorWithContext(ctx context.Context, logger log.Logger, message string, args ...interface{}) error {
	return logWithContext(ctx, level.Error(logger), message, args)
}

// WarnWithContext middleware function print warn log with context
func WarnWithContext(ctx context.Context, logger log.Logger, message string, args ...interface{}) error {
	return logWithContext(ctx, level.Warn(logger), message, args)
}

// InfoWithContext middleware function print info log with context
func InfoWithContext(ctx context.Context, logger log.Logger, message string, args ...interface{}) error {
	return logWithContext(ctx, level.Info(logger), message, args)
}

// DebugWithContext middleware function print debug log with context
func DebugWithContext(ctx context.Context, logger log.Logger, message string, args ...interface{}) error {
	return logWithContext(ctx, level.Debug(logger), message, args)
}

// Error middleware function print error log without context
func Error(logger log.Logger, message string, args ...interface{}) error {
	return l(level.Error(logger), message, args)
}

// Warn middleware function print warn log without context
func Warn(logger log.Logger, message string, args ...interface{}) error {
	return l(level.Warn(logger), message, args)
}

// Info middleware function print info log without context
func Info(logger log.Logger, message string, args ...interface{}) error {
	return l(level.Info(logger), message, args)
}

// Debug middleware function print debug log without context
func Debug(logger log.Logger, message string, args ...interface{}) error {
	return l(level.Debug(logger), message, args)
}

func logWithContext(ctx context.Context, logger log.Logger, message string, args ...interface{}) error {
	correlationID := tlecontext.GetCorrelation(ctx)
	duration, deadline := tlecontext.GetTimeout(ctx)
	return logger.Log(
		"message", message,
		"correaltionId", correlationID,
		"duration", duration,
		"deadline", deadline,
		"args", args)
}

func l(logger log.Logger, message string, args ...interface{}) error {
	return logger.Log(
		"message", message,
		"args", args)
}
