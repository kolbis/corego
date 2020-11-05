package ratelimit

import (
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/ratelimit"
	"golang.org/x/time/rate"
)

// NewDefaultErrorLimitterMiddleware returns an endpoint.Middleware that acts as a rate limiter.
// Requests that would exceed the maximum request rate are simply rejected with an error.
// NewDefaultLimitterMiddleware will create a middleware which retry every second and allows 100 concurrent requests per second.
func NewDefaultErrorLimitterMiddleware() endpoint.Middleware {
	config := NewDefaultRateLimiterConfig()
	return NewErrorLimitterMiddleware(config)
}

// NewErrorLimitterMiddleware returns an endpoint.Middleware that acts as a rate limiter.
// Requests that would exceed the maximum request rate are simply rejected with an error.
func NewErrorLimitterMiddleware(config RateLimiterConfig) endpoint.Middleware {
	l := rate.NewLimiter(config.Limit, config.MaxQueriesPerSecond)
	middleware := ratelimit.NewErroringLimiter(l)
	return middleware
}

// NewDefaultDelayingLimitterMiddleware returns an endpoint.Middleware that acts as a request throttler.
// Requests that would exceed the maximum request rate are delayed via the Waiter function
// NewDefaultLimitterMiddleware will create a middleware which retry every second and allows 100 concurrent requests per second.
func NewDefaultDelayingLimitterMiddleware() endpoint.Middleware {
	config := NewDefaultRateLimiterConfig()
	return NewDelayingLimitterMiddleware(config)
}

// NewDelayingLimitterMiddleware returns an endpoint.Middleware that acts as a request throttler.
// Requests that would exceed the maximum request rate are delayed via the Waiter function
func NewDelayingLimitterMiddleware(config RateLimiterConfig) endpoint.Middleware {
	l := rate.NewLimiter(config.Limit, config.MaxQueriesPerSecond)
	middleware := ratelimit.NewDelayingLimiter(l)
	return middleware
}

var (
	rateLimitError string = ratelimit.ErrLimited.Error()
)

// IsRateLimitted will check if the error is caused by rate error rate limit
func IsRateLimitted(err error) bool {
	if err != nil && err.Error() == rateLimitError {
		return true
	}
	return false
}
