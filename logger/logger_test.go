package logger_test

import (
	"context"
	"testing"

	logger "github.com/kolbis/corego/logger"
)

func TestLoggerWithContextReturnNil(t *testing.T) {
	ctx := context.Background()
	nl := logger.NewNopLogger()

	logRes := logger.ErrorWithContext(ctx, nl, "this is an error with context")
	if logRes != nil {
		t.Errorf("TestLoggerWithContextReturnNil return error %v", logRes)
	}
}
