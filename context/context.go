package context

import (
	"context"
	"time"
)

// Root is the context which is created for a service / application when it initialized
// The context created has no expiration or deadline
// When we need to call external service, we need to create a new context with deadlines, durtion and correlation ID
func Root() context.Context {
	ctx := context.Background()
	corrid := NewCorrelation()
	ctx = SetRootCorrealtion(ctx, corrid)

	return ctx
}

// SetCorrealtion will set into the given context a corraltion ID value
func SetCorrealtion(ctx context.Context, correlationID string) context.Context {
	ctx = context.WithValue(ctx, CorrelationIDKey, correlationID)
	return ctx
}

// SetRootCorrealtion will set into the given context a corraltion ID value
func SetRootCorrealtion(ctx context.Context, correlationID string) context.Context {
	ctx = context.WithValue(ctx, CorrelationIDRootKey, correlationID)
	return ctx
}

// SetTimeout will set in to the given context the duration and the deadline
func SetTimeout(ctx context.Context, duration time.Duration, deadline time.Time) context.Context {
	ctx = context.WithValue(ctx, DurationKey, duration)
	ctx = context.WithValue(ctx, DeadlineKey, deadline)
	return ctx
}

// SetError will set in to the given context an error
// For example, in rabbit, before returning the response, we will check for an error and based on it will decide if to ack
func SetError(ctx context.Context, err error, errCode int) context.Context {
	ctx = context.WithValue(ctx, ErrorKey, err)
	ctx = context.WithValue(ctx, ErrorCodeKey, errCode)
	return ctx
}

// GetError will return the error and error code which were set into the context
func GetError(ctx context.Context) (bool, int, error) {
	errAsInterface := ctx.Value(ErrorKey)
	errCodeAsInterface := ctx.Value(ErrorCodeKey)

	var errCode int
	var err error
	hasError := false
	if errAsInterface != nil {
		err = errAsInterface.(error)
		hasError = true
	}
	if errCodeAsInterface != nil {
		errCode = errCodeAsInterface.(int)
		hasError = true
	}

	return hasError, errCode, err
}

// GetTimeout will return the duration and the deadline from the given context
// If it cannot find it, it will respectively return nil
func GetTimeout(ctx context.Context) (time.Duration, time.Time) {
	durationAsInterface := ctx.Value(DurationKey)
	deadlineAsInterface := ctx.Value(DeadlineKey)

	var duration time.Duration
	var deadline time.Time

	if durationAsInterface != nil {
		duration = durationAsInterface.(time.Duration)
	}
	if deadlineAsInterface != nil {
		deadline = deadlineAsInterface.(time.Time)
	}

	return duration, deadline
}

// GetCorrelation will return the correlation ID from the context
// if it cannot find it it will try to get the root correlation ID
// If it cannot find it, it will return nil
func GetCorrelation(ctx context.Context) string {
	val := ctx.Value(CorrelationIDKey)
	var corrid string

	if val == nil {
		val = GetRootCorrelation(ctx)
	}

	if val != nil {
		corrid = val.(string)
	}

	return corrid
}

// GetRootCorrelation will return the correlation ID from the context
// If it cannot find it, it will return nil
func GetRootCorrelation(ctx context.Context) string {
	val := ctx.Value(CorrelationIDRootKey)

	var corrid string
	if val != nil {
		corrid = val.(string)
	}

	return corrid
}

// GetOrCreateCorrelation will get the correlation ID from the context
// If it cannot find, it will try to use the root correlation ID
// If it does not exist, it will create a new correlation ID
func GetOrCreateCorrelation(ctx context.Context) string {
	val := ctx.Value(CorrelationIDKey)
	var corrid string

	if val == nil {
		val = GetRootCorrelation(ctx)
		if val == nil {
			corrid = NewCorrelation()
		}
	} else {
		corrid = val.(string)
	}

	return corrid
}

// GetOrCreateTimeout will get duration and deadline from the context
// if it does not exist it will create new duration and deadline
// If appendToContext, it will update the input context with the duration and deadline
func GetOrCreateTimeout(ctx context.Context) (time.Duration, time.Time, context.Context) {
	return GetOrCreateTimeoutFromContext(ctx, false)
}

// GetOrCreateTimeoutFromContext will get duration and deadline from the context
// if it does not exist it will create new duration and deadline
// If appendToContext, it will update the input context with the duration and deadline
func GetOrCreateTimeoutFromContext(ctx context.Context, appendToContext bool) (time.Duration, time.Time, context.Context) {
	duration, deadline := GetTimeout(ctx)

	if deadline.IsZero() {
		t := NewTimeoutCalculator()
		duration, deadline = t.NewTimeout()

		if appendToContext {
			ctx = SetTimeout(ctx, duration, deadline)
		}
	}

	return duration, deadline, ctx
}

// GetOrCreateCorrelationFromContext will get correlation ID from the context
// if it does not exist it will create new correlation ID
// If appendToContext, it will update the input context with the correlation ID
func GetOrCreateCorrelationFromContext(ctx context.Context, appendToContext bool) (string, context.Context) {
	corrid := GetCorrelation(ctx)

	if corrid == "" {
		corrid = NewCorrelation()
		if appendToContext {
			ctx = SetCorrealtion(ctx, corrid)
		}
	}

	return corrid, ctx
}

// NewCorrelation ...
func NewCorrelation() string {
	c := NewCorrelationID()
	return c.New()
}
