package rabbitmq

import (
	"context"

	amqpkit "github.com/go-kit/kit/transport/amqp"
	"github.com/streadway/amqp"
)

// Subscriber stored the configurations and services required to build a subsripber to incoming messages
// It defines the topology required for the subscriber including how to constract a queue and exchange
type Subscriber struct {
	KitSubscriber         *amqpkit.Subscriber
	ConnectionManager     *ConnectionManager
	Channel               *amqp.Channel
	IsConnected           bool
	ExchangeName          string
	QueueName             string
	BuildQueueTopology    func(channel *amqp.Channel, queueName string) (amqp.Queue, error)
	BuildExchangeTopology func(*amqp.Channel, string) error
	BindQueueTopology     func(*amqp.Channel, string, string) error
	ConsumeTopology       func(*amqp.Channel, string) (<-chan amqp.Delivery, error)
	QosTopology           func(ch *amqp.Channel) error
}

// Consume will create a delivery channel which messages will be send to as the are consumed from the queue
// It will ensure that the queue, exchange are configured before opening the delivery channel
func (sub *Subscriber) Consume(ch *amqp.Channel) (<-chan amqp.Delivery, error) {
	sub.Channel = ch
	sub.QosTopology(sub.Channel)
	sub.BuildQueueTopology(sub.Channel, sub.QueueName)
	sub.BuildExchangeTopology(sub.Channel, sub.ExchangeName)
	sub.BindQueueTopology(sub.Channel, sub.QueueName, sub.ExchangeName)
	c, err := sub.ConsumeTopology(sub.Channel, sub.QueueName)

	return c, err
}

// Close will shutdown the associated subscriber resources (channel)
func (sub *Subscriber) Close(ctx context.Context) error {
	conn := *sub.ConnectionManager
	err := conn.CloseChannel(ctx, sub.Channel)
	if err == nil {
		sub.IsConnected = false
	}
	return err
}

func (sub *Subscriber) connect() error {
	conn := *sub.ConnectionManager
	ch, err := conn.GetChannel()
	if err == nil {
		sub.Channel = ch
		sub.IsConnected = true
	}
	//p.changeConnection(ctx, conn, ch)
	return err
}
