package cache

import (
	"math/rand"
	"time"
)

const (
	// DefaultJitterFactor factor to apply by default on the jitter
	DefaultJitterFactor float64 = 0.10
)

// Jitter returns a time.Duration between duration and duration + maxFactor * duration,
// to allow clients to avoid converging on periodic behavior.  If maxFactor is 0.0, a
// suggested default value will be chosen.
func Jitter(duration time.Duration, maxFactor float64) time.Duration {
	if maxFactor <= 0.0 {
		maxFactor = DefaultJitterFactor
	}

	randomRange := 2*rand.Float64() - 1 // Return a rendom range between [-1...1]
	jitter := time.Duration(randomRange * maxFactor * float64(duration))

	wait := duration + jitter
	return wait
}
