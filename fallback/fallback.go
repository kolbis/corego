package fallback

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/kolbis/corego/circuitbreaker"
	"github.com/kolbis/corego/ratelimit"
)

// NewFallbackMiddleware will create a fallback middleware
// the middleware will wait for Circuit to be opened. Once it is opened, it will execute the input endppoint
// and will return its result. The result can be default value or a result from the fallback cache
func NewFallbackMiddleware(fallbackEndpoint endpoint.Endpoint) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			resp, err := next(ctx, request)
			if circuitbreaker.IsCircuitOpen(err) || ratelimit.IsRateLimitted(err) {
				return fallbackEndpoint(ctx, request)
			}

			return resp, err
		}
	}
}
