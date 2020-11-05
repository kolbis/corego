package rabbitmq

import (
	"context"
	"encoding/json"

	"github.com/go-kit/kit/endpoint"
	amqptransport "github.com/go-kit/kit/transport/amqp"
	"github.com/streadway/amqp"
	tlectx "github.com/kolbis/corego/context"
	tlectxrabbit "github.com/kolbis/corego/context/transport/rabbitmq"
	"github.com/kolbis/corego/errors"
)

// Publisher is used to publish messages to rabbit
type Publisher interface {
	Close(context.Context) error
	Publish(context.Context, *Message, string) error
}

type publisher struct {
	connectionManager *ConnectionManager
	ch                *amqp.Channel
	isConnected       bool
	topology          Topology
}

// NewPublisher will create a new publisher and will establish a connection to rabbit
func NewPublisher(conn *ConnectionManager) Publisher {
	p := publisher{
		connectionManager: conn,
		topology:          NewTopology(),
	}
	p.connect()
	return &p
}

// Publish will publish a message into the requested exchange
// if the exchange do not exist it will create it
func (p publisher) Publish(ctx context.Context, message *Message, exchangeName string) error {
	if p.isConnected == false {
		return errors.NewApplicationErrorf("before publishing, you must connect to rabbitMQ")
	}

	p.buildExchange(exchangeName)
	ep, _ := p.publishEndpoint(ctx, message, exchangeName, defaultRequestEncoder(exchangeName))
	_, err := ep(ctx, message)

	return err
}

// PublishEndpoint creates an endpoint that by calling on it, it will publish a message
// You dont need to call it directly. Use the publish method
func (p publisher) publishEndpoint(ctx context.Context, message *Message, exchangeName string, encodeFunc amqptransport.EncodeRequestFunc, options ...amqptransport.PublisherOption) (endpoint.Endpoint, error) {
	duration, _ := tlectx.GetTimeout(ctx)

	// building the publisher options
	before := amqptransport.PublisherBefore(
		// must come first since it creates the transport context!
		tlectxrabbit.WriteMessageRequestFunc(),
		amqptransport.SetPublishDeliveryMode(2),
		amqptransport.SetPublishExchange(exchangeName))

	ops := make([]amqptransport.PublisherOption, 0)
	ops = append(ops, options...)
	ops = append(ops, amqptransport.PublisherTimeout(duration), amqptransport.PublisherDeliverer(amqptransport.SendAndForgetDeliverer))
	ops = append(ops, before)

	publisher := amqptransport.NewPublisher(p.ch, &amqp.Queue{Name: ""}, encodeFunc, noopResponseDecoder, ops...)

	return publisher.Endpoint(), nil
}

func (p publisher) buildExchange(exchanegName string) {
	conn := *p.connectionManager
	ch, err := conn.GetChannel()
	if err == nil {
		defer ch.Close()
		p.topology.BuildDurableExchange(ch, exchanegName)
	}
}

// noopResponseDecoder is a response decoder which does nothing
// It is used for a fire and forget messages
func noopResponseDecoder(ctx context.Context, d *amqp.Delivery) (response interface{}, err error) {
	return struct{}{}, nil
}

func defaultRequestEncoder(exchangeName string) func(context.Context, *amqp.Publishing, interface{}) error {
	f := func(ctx context.Context, p *amqp.Publishing, request interface{}) error {
		message := request.(*Message)
		body, err := json.Marshal(message)
		p.Body = body
		return err
	}
	return f
}

// Close will close the publisher rabbit channel
func (p *publisher) Close(ctx context.Context) error {
	var err error

	if p.isConnected && p.ch != nil {
		cherr := p.ch.Close()

		if cherr != nil {
			err = errors.NewApplicationError(err, cherr.Error())
		} else {
			p.isConnected = false
		}
	}

	return err
}

func (p *publisher) connect() error {
	connMgr := *p.connectionManager
	ch, err := connMgr.GetChannel()
	if err == nil {
		p.ch = ch
		//p.changeConnection(ctx, conn, ch)
		p.isConnected = true
	}
	return err
}
