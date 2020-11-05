package rabbitmq

import (
	"context"
	"time"

	kitamqptransport "github.com/go-kit/kit/transport/amqp"
	"github.com/streadway/amqp"
	tlectx "github.com/kolbis/corego/context"
	"github.com/kolbis/corego/context/transport"
	"github.com/kolbis/corego/utils"
)

const (
	deadlineHeaderKey       string = "tle-deadline-unix"
	durationHeaderKey       string = "tle-duration-ms"
	callerProcessHeaderKey  string = "tle-caller-process"
	callerHostNameHeaderKey string = "tle-caller-hostname"
	callerPIDHeaderKey      string = "tle-caller-processid"
	callerOSHeaderKey       string = "tle-caller-os"
)

type amqptransport struct {
}

// NewTransport will create a new AMQP transport
func NewTransport() transport.Transport {
	return amqptransport{}
}

func (amqptrans amqptransport) Read(ctx context.Context, req interface{}) (context.Context, context.CancelFunc) {
	del := req.(*amqp.Delivery)

	headerCorrelationID := del.CorrelationId
	headerDuration := del.Headers[durationHeaderKey].(string)
	headerDeadline := del.Headers[deadlineHeaderKey].(string)

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

func (amqptrans amqptransport) Write(ctx context.Context, req interface{}) (context.Context, context.CancelFunc) {
	newCtx, cancel := transport.CreateTransportContext(ctx)
	corrid := tlectx.GetCorrelation(newCtx)
	duration, deadline := tlectx.GetTimeout(newCtx)

	pub := req.(*amqp.Publishing)
	pub.MessageId = utils.NewUUID()
	pub.Timestamp = utils.NewDateTime().Now()
	pub.CorrelationId = corrid
	pub.ContentType = "application/vnd.masstransit+json"
	headers := setHeaders(pub.Headers, duration, deadline)
	pub.Headers = headers

	return newCtx, cancel
}

// ReadMessageRequestFunc will be executed once a message is consumed
// it will read from the delivery and will create a context
func ReadMessageRequestFunc() kitamqptransport.RequestFunc {
	return func(ctx context.Context, _ *amqp.Publishing, del *amqp.Delivery) context.Context {
		t := NewTransport()
		newCtx, _ := t.Read(ctx, del)
		return newCtx
	}
}

// WriteMessageRequestFunc ...
func WriteMessageRequestFunc() kitamqptransport.RequestFunc {
	return func(ctx context.Context, pub *amqp.Publishing, _ *amqp.Delivery) context.Context {
		t := NewTransport()
		newCtx, _ := t.Write(ctx, pub)
		return newCtx
	}
}

func setHeaders(headers amqp.Table, duration time.Duration, deadline time.Time) amqp.Table {
	if headers == nil {
		headers = amqp.Table{}
	}

	conv := utils.NewConvertor()
	durationHeader := conv.FromInt64ToString(conv.DurationToMiliseconds(duration))
	deadlineHeader := conv.FromInt64ToString(conv.FromTimeToUnix(deadline))

	headers[deadlineHeaderKey] = deadlineHeader
	headers[durationHeaderKey] = durationHeader
	headers[callerProcessHeaderKey] = utils.ProcessName()
	headers[callerHostNameHeaderKey] = utils.HostName()
	headers[callerPIDHeaderKey] = utils.ProcessID()
	headers[callerOSHeaderKey] = utils.OperatingSystem()

	return headers
}
