package http

import (
	"context"
	"net/http"
	"time"

	httpkit "github.com/go-kit/kit/transport/http"
	tlectx "github.com/kolbis/corego/context"
	"github.com/kolbis/corego/context/transport"
	"github.com/kolbis/corego/utils"
)

type httptransport struct {
}

// NewTransport will create a new HTTP transport
func NewTransport() transport.Transport {
	return httptransport{}
}

// Read will read from http.Request into context
func (httptrans httptransport) Read(ctx context.Context, req interface{}) (context.Context, context.CancelFunc) {
	r := req.(*http.Request)

	headerCorrelationID := r.Header.Get(string(tlectx.CorrelationIDKey))
	headerDuration := r.Header.Get(string(tlectx.DurationKey))
	headerDeadline := r.Header.Get(string(tlectx.DeadlineKey))

	correlationID := headerCorrelationID
	if headerCorrelationID == "" {
		correlationID = tlectx.NewCorrelation()
	}

	var duration time.Duration
	var deadline time.Time
	if headerDuration == "" || headerDeadline == "" {
		t := tlectx.NewTimeoutCalculator()
		duration, deadline = t.NewTimeout()
	} else {
		conv := utils.NewConvertor()
		duration = conv.MilisecondsToDuration(conv.FromStringToInt64(headerDuration))
		deadline = conv.FromUnixToTime(conv.FromStringToInt64(headerDeadline))
	}

	ctx = tlectx.SetCorrealtion(ctx, correlationID)
	ctx = tlectx.SetTimeout(ctx, duration, deadline)
	ctx, cancel := context.WithDeadline(ctx, deadline)

	return ctx, cancel
}

// Write will write from context into http.Request
func (httptrans httptransport) Write(ctx context.Context, req interface{}) (context.Context, context.CancelFunc) {
	r := req.(*http.Request)
	conv := utils.NewConvertor()

	newCtx, cancel := transport.CreateTransportContext(ctx)
	corrid := tlectx.GetCorrelation(newCtx)
	duration, deadline := tlectx.GetTimeout(newCtx)

	durationHeader := conv.FromInt64ToString(conv.DurationToMiliseconds(duration))
	deadlineHeader := conv.FromInt64ToString(conv.FromTimeToUnix(deadline))

	r.Header.Add(string(tlectx.CorrelationIDKey), corrid)
	r.Header.Add(string(tlectx.DurationKey), durationHeader)
	r.Header.Add(string(tlectx.DeadlineKey), deadlineHeader)

	return newCtx, cancel
}

func write(ctx context.Context, r *http.Request) context.Context {
	t := NewTransport()
	newCtx, _ := t.Write(ctx, r)
	return newCtx
}

func read(ctx context.Context, r *http.Request) context.Context {
	t := NewTransport()
	newCtx, _ := t.Read(ctx, r)
	return newCtx
}

// WriteBefore ...
func WriteBefore() httpkit.ClientOption {
	return httpkit.ClientBefore(write)
}

// ReadBefore ...
func ReadBefore() httpkit.ServerOption {
	return httpkit.ServerBefore(read)
}
