package rabbitmq

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	amqpkit "github.com/go-kit/kit/transport/amqp"
	amqptransport "github.com/go-kit/kit/transport/amqp"
	"github.com/streadway/amqp"
	tlectxamqp "github.com/kolbis/corego/context/transport/rabbitmq"
)

// SubscribeManager responsible to create rabbit subscribers
type SubscribeManager interface {

	// NewPrivateSubscriber will create a private subscriber
	// Private subscriber is non durable subscriber which do not share queues
	// We use this kind of subscribers when we want each consumer to recieve a copy of a message
	// for example for invalidating data event, where we want each instance to clear its data from the cache
	NewPrivateSubscriber(exchangeName string, queueName string, endpoint endpoint.Endpoint, dec amqptransport.DecodeRequestFunc, enc amqptransport.EncodeResponseFunc, options ...amqptransport.SubscriberOption) Subscriber

	// NewCommandSubscriber will create a command subscriber
	// Command subscriber is a durable subscriber which share the given queue with other subscribers using a competing queue pattern
	NewCommandSubscriber(exchangeName string, queueName string, endpoint endpoint.Endpoint, dec amqptransport.DecodeRequestFunc, enc amqptransport.EncodeResponseFunc, options ...amqptransport.SubscriberOption) Subscriber
}

type submgr struct {
	connMgr  *ConnectionManager
	topology Topology
}

// NewSubscriberManager will create a new instance of the subscriber manager
// Recommended to have a single one per application
func NewSubscriberManager(connMgr *ConnectionManager) SubscribeManager {
	s := submgr{
		connMgr:  connMgr,
		topology: NewTopology(),
	}

	return s
}

func (s submgr) NewCommandSubscriber(
	exchangeName string,
	queueName string,
	endpoint endpoint.Endpoint,
	dec amqptransport.DecodeRequestFunc,
	enc amqptransport.EncodeResponseFunc,
	options ...amqptransport.SubscriberOption,
) Subscriber {

	queueName = BuildCommandQueueName(queueName)
	sub := newKitSubscriber(endpoint, exchangeName, dec, enc, options...)
	return Subscriber{
		ConnectionManager:     s.connMgr,
		KitSubscriber:         sub,
		QueueName:             queueName,
		ExchangeName:          exchangeName,
		BuildQueueTopology:    s.topology.BuildDurableQueue,
		BuildExchangeTopology: s.topology.BuildDurableExchange,
		BindQueueTopology:     s.topology.QueueBind,
		ConsumeTopology:       s.topology.Consume,
		QosTopology:           s.topology.Qos,
	}
}

func (s submgr) NewPrivateSubscriber(
	exchangeName string,
	queueName string,
	endpoint endpoint.Endpoint,
	dec amqptransport.DecodeRequestFunc,
	enc amqptransport.EncodeResponseFunc,
	options ...amqptransport.SubscriberOption,
) Subscriber {

	queueName = BuildPrivateQueueName(queueName)
	sub := newKitSubscriber(endpoint, exchangeName, dec, enc, options...)
	return Subscriber{
		ConnectionManager:     s.connMgr,
		KitSubscriber:         sub,
		QueueName:             queueName,
		ExchangeName:          exchangeName,
		BuildQueueTopology:    s.topology.BuildNonDurableQueue,
		BuildExchangeTopology: s.topology.BuildNonDurableExchange,
		BindQueueTopology:     s.topology.QueueBind,
		ConsumeTopology:       s.topology.Consume,
		QosTopology:           s.topology.Qos,
	}
}

func newKitSubscriber(
	endpoint endpoint.Endpoint,
	exchangeName string,
	dec amqptransport.DecodeRequestFunc,
	enc amqptransport.EncodeResponseFunc,
	options ...amqptransport.SubscriberOption) *amqptransport.Subscriber {

	ops := make([]amqpkit.SubscriberOption, 0)
	ops = append(ops, options...)
	ops = append(ops, amqptransport.SubscriberResponsePublisher(nopResponsePublisher))
	ops = append(ops, amqptransport.SubscriberErrorEncoder(amqptransport.ReplyErrorEncoder))
	ops = append(
		ops,
		amqptransport.SubscriberBefore(
			amqptransport.SetPublishExchange(exchangeName),
			tlectxamqp.ReadMessageRequestFunc(),
			amqptransport.SetPublishDeliveryMode(2),
		))

	sub := amqptransport.NewSubscriber(endpoint, dec, enc, ops...)

	return sub
}

func nopResponsePublisher(ctx context.Context, deliv *amqp.Delivery, ch amqpkit.Channel, pub *amqp.Publishing) error {
	// if there was an error we dont want to ack
	shoundNack := ShouldNack(ctx)
	if shoundNack {
		deliv.Nack(false, true)
	}
	return deliv.Ack(false)
}
