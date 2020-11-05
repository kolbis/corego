package http

import (
	"context"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	tlesb "github.com/kolbis/corego/circuitbreaker"
	tlectx "github.com/kolbis/corego/context"
	"github.com/kolbis/corego/errors"
	tlelimit "github.com/kolbis/corego/ratelimit"
)

// Wrap will wrap the data in a Request while copying the transport correlation id, duration and timeout
func (r Request) Wrap(ctx context.Context, data interface{}) Request {
	corrid := tlectx.GetOrCreateCorrelation(ctx)

	req := Request{
		Data:          data,
		CorrelationID: corrid,
	}
	return req
}

// Execute the endpoint and returns an anlyzed Response object
// It include information about the status code, circuit, rate limit
// Use this function from the client in order to execute rest endpoint on the server
func Execute(ctx context.Context, req interface{}, ep endpoint.Endpoint) Response {
	var res interface{}
	var err error
	var statusCode int
	circuitOpen := false
	rateLimitted := false

	if res, err = ep(ctx, req); err != nil {
		rateLimitted = tlelimit.IsRateLimitted(err)
		circuitOpen = tlesb.IsCircuitOpen(err)
	}

	response, ok := res.(*http.Response)
	if ok {
		statusCode = response.StatusCode
	} else if err != nil {
		statusCode = 500
	} else {
		statusCode = 200
	}

	if statusCode >= 400 && err == nil {
		err = errors.NewApplicationErrorf("failed to execute endpoint")
	}

	return Response{
		Data:               res,
		Error:              err,
		IsCircuitOpen:      circuitOpen,
		IsReachedRateLimit: rateLimitted,
		StatusCode:         statusCode,
	}
}
