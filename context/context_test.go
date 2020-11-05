package context_test

import (
	"context"
	"testing"
	"time"

	tlectx "github.com/kolbis/corego/context"
	"github.com/kolbis/corego/errors"
)

func TestNewFrom(t *testing.T) {
	ctx := context.Background()
	corrid := "12345"
	duration := time.Second * 10
	deadline := time.Now().UTC().Add(duration)

	ctx = tlectx.SetCorrealtion(ctx, corrid)
	ctx = tlectx.SetTimeout(ctx, duration, deadline)
	ctx, cancel := context.WithDeadline(ctx, deadline)

	if cancel == nil {
		t.Fail()
	}
}

func TestSetContext(t *testing.T) {
	ctx := context.Background()
	corrid := "12345"
	duration := time.Second * 10
	deadline := time.Now().UTC().Add(duration)

	ctx = tlectx.SetCorrealtion(ctx, corrid)
	corridResult := tlectx.GetCorrelation(ctx)

	if corrid != corridResult {
		t.Error("correlation does not match")
	}

	ctx = tlectx.SetTimeout(ctx, duration, deadline)
	durationResult, deadlineResult := tlectx.GetTimeout(ctx)

	if durationResult != duration {
		t.Error("duration does not match")
	}

	if deadlineResult != deadline {
		t.Error("deadline does not match")
	}
}

func TestGetOrCreateCorrelationID_Create(t *testing.T) {
	ctx := context.Background()
	var corrid string
	corrid, ctx = tlectx.GetOrCreateCorrelationFromContext(ctx, true)
	actualCorrelationID := tlectx.GetCorrelation(ctx)

	if actualCorrelationID == "" {
		t.Error("correlation was not set to context")
	}

	if corrid != actualCorrelationID {
		t.Error("correlation does not match")
	}
}

func TestGetOrCreateTimeout_Create(t *testing.T) {
	ctx := context.Background()

	duration, deadline, ctx := tlectx.GetOrCreateTimeoutFromContext(ctx, true)
	actualDuration, actualDeadlne := tlectx.GetTimeout(ctx)

	if actualDuration != duration {
		t.Error("duration does not match")
	}

	if deadline != actualDeadlne {
		t.Error("deadline does not match")
	}
}

func TestGetOrCreateError_ErrorExist(t *testing.T) {
	ctx := context.Background()
	err := errors.NewApplicationErrorf("error!")
	errCode := 99
	ctx = tlectx.SetError(ctx, err, errCode)
	hasError, isCode, isError := tlectx.GetError(ctx)

	if hasError == false {
		t.Error("expected an error")
	}

	if errCode != isCode {
		t.Error("err code do not match")
	}

	if errors.IsApplicationError(isError) == false {
		t.Error("expected an application error")
	}

	if isError.Error() != err.Error() {
		t.Error("errors do not match in content")
	}
}

func TestGetOrCreateError_ErrorNotExist(t *testing.T) {
	ctx := context.Background()
	hasError, isCode, isError := tlectx.GetError(ctx)

	if hasError == true {
		t.Error("did not expect an error")
	}

	if isCode != 0 {
		t.Error("did not expected an err code")
	}

	if isError != nil {
		t.Error("expected a nil error")
	}
}
