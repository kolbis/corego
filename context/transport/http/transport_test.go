package http_test

import (
	"context"
	"net/http/httptest"
	"testing"

	tlectx "github.com/kolbis/corego/context"
	"github.com/kolbis/corego/context/transport"
	"github.com/kolbis/corego/context/transport/http"
)

func TestReadWriteWhenCorrelationAndTimeoutExistOnCaller(t *testing.T) {
	httpTransport := http.NewTransport()

	// root context
	ctxRoot := tlectx.Root()
	ctxA, _ := transport.CreateTransportContext(ctxRoot)
	req := httptest.NewRequest("GET", "http://example.com/foo", nil)

	// writing all the headers
	newCtxA, _ := httpTransport.Write(ctxA, req)

	// reading all the headers back
	ctxB, _ := httpTransport.Read(context.Background(), req)

	corridA := tlectx.GetCorrelation(ctxA)
	corridNewA := tlectx.GetCorrelation(newCtxA)
	corridB := tlectx.GetCorrelation(ctxB)
	corridRoot := tlectx.GetCorrelation(ctxRoot)

	// all services should have the same correlation ID
	if corridA != corridB || corridA != corridRoot || corridA != corridNewA {
		t.Error("correlation id should be same")
	}

	durationA, deadlineA := tlectx.GetTimeout(newCtxA)
	durationB, deadlineB := tlectx.GetTimeout(ctxB)

	// all services should have the same correlation ID
	if deadlineA.Before(deadlineB) {
		t.Error("A deadline should be after B deadline")
	}

	if durationA <= durationB {
		t.Error("B allowed duration should be smaller than A duration")
	}
}

func TestReadWriteWhenCorrelationAndTimeoutNotExistOnCaller(t *testing.T) {
	httpTransport := http.NewTransport()

	ctx := context.Background()
	req := httptest.NewRequest("GET", "http://example.com/foo", nil)

	// writing all the headers
	ctxA, _ := httpTransport.Write(ctx, req)

	// reading all the headers back
	ctxB, _ := httpTransport.Read(context.Background(), req)

	corridA := tlectx.GetCorrelation(ctxA)
	corridB := tlectx.GetCorrelation(ctxB)

	// all services should have the same correlation ID
	if corridA != corridB {
		t.Error("correlation id should be same")
	}

	durationA, deadlineA := tlectx.GetTimeout(ctxA)
	durationB, deadlineB := tlectx.GetTimeout(ctxB)

	if deadlineA.Before(deadlineB) == true {
		t.Error("deadlines should be the same")
	}

	if durationA != durationB {
		t.Error("durations should be same")
	}
}
