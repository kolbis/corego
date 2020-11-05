package transport_test

import (
	"testing"

	tlectx "github.com/kolbis/corego/context"
	"github.com/kolbis/corego/context/transport"
)

func TestCreateTransportContext(t *testing.T) {
	// root context
	ctxRoot := tlectx.Root()

	// first service in flow
	ctxA, _ := transport.CreateTransportContext(ctxRoot)

	// second service in flow
	ctxB, _ := transport.CreateTransportContext(ctxA)

	corridA := tlectx.GetCorrelation(ctxA)
	corridB := tlectx.GetCorrelation(ctxB)
	corridRoot := tlectx.GetCorrelation(ctxRoot)

	// all services should have the same correlation ID
	if corridA != corridB || corridA != corridRoot {
		t.Fail()
	}

	durationA, deadlineA := tlectx.GetTimeout(ctxA)
	durationB, deadlineB := tlectx.GetTimeout(ctxB)

	// all services should have the same correlation ID
	if deadlineA.Before(deadlineB) {
		t.Fail()
	}

	if durationA <= durationB {
		t.Fail()
	}
}
