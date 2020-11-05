package context

import (
	"context"
	"time"

	"github.com/kolbis/corego/utils"
)

const (
	// MaxTimeout is 15 seconds
	MaxTimeout time.Duration = time.Second * 15
)

// TimeoutCalculator ...
type TimeoutCalculator interface {
	NextTimeoutFromContext(context.Context) (time.Duration, time.Time)
	NextTimeout(time.Duration, time.Time) (time.Duration, time.Time)
	NewTimeout() (time.Duration, time.Time)
}

type timeoutcalc struct{}

// NewTimeoutCalculator creates a new instance of timeout calculator
func NewTimeoutCalculator() TimeoutCalculator {
	return timeoutcalc{}
}

// NextTimeoutFromContext calculate the next deadline and duration from context which should be used by downstream services
func (c timeoutcalc) NextTimeoutFromContext(ctx context.Context) (time.Duration, time.Time) {
	duration, deadline := GetTimeout(ctx)
	return c.NextTimeout(duration, deadline)
}

// NextTimeout calculate the next deadline and duration which should be used by downstream services
func (c timeoutcalc) NextTimeout(duration time.Duration, deadline time.Time) (time.Duration, time.Time) {
	if deadline.IsZero() {
		return c.NewTimeout()
	}

	dt := utils.DateTime{}
	deadline = deadline.Add(time.Second * -2)
	duration = dt.DiffFromNow(deadline)

	return duration, deadline
}

// NewTimeout will return new timeout including the timeout duration as time.Duration and deadline as time.Time
func (c timeoutcalc) NewTimeout() (time.Duration, time.Time) {
	duration := MaxTimeout
	deadline := utils.NewDateTime().AddDuration(duration)

	return duration, deadline
}
