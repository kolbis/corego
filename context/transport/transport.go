package transport

import (
	"context"

	tlectx "github.com/kolbis/corego/context"
)

// Transport contract is used to read and write context information when tranporting between services
type Transport interface {
	// Read will read from transport into context
	Read(context.Context, interface{}) (context.Context, context.CancelFunc)

	// Write will write from context into transport
	Write(context.Context, interface{}) (context.Context, context.CancelFunc)
}

// CreateTransportContext will create a new transport context which will be used to write into request transport
// For example, before we send http request, we will CreateTransportContext, and the context which was created will be used
// to build http request headers
func CreateTransportContext(ctx context.Context) (context.Context, context.CancelFunc) {
	calc := tlectx.NewTimeoutCalculator()
	var cancel context.CancelFunc

	_, newCtx := tlectx.GetOrCreateCorrelationFromContext(ctx, true)
	duration, deadline := calc.NextTimeoutFromContext(newCtx)

	newCtx = tlectx.SetTimeout(newCtx, duration, deadline)
	newCtx, cancel = context.WithDeadline(newCtx, deadline)

	return newCtx, cancel
}
