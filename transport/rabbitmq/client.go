package rabbitmq

import (
	"context"
	"fmt"

	"github.com/kolbis/corego/errors"
	tlelogger "github.com/kolbis/corego/logger"
)

// Client is a rabbit contract to publish and consume messages
// It provide an easy access to the publish anc consume APIs
type Client interface {
	Consume(context.Context)
	Publish(context.Context, *Message, string) error
	Close(context.Context) error
}

type client struct {
	connectionManager *ConnectionManager
	logger            tlelogger.Logger
	subscribers       *[]Subscriber
	publisher         *Publisher
}

// NewClient will create a new instance of a client
// Best practice is to have a single one per application and reuse it
func NewClient(connMgr *ConnectionManager, logManager tlelogger.Logger, publisher *Publisher, subscribers *[]Subscriber) Client {
	return &client{
		logger:            logManager,
		subscribers:       subscribers,
		publisher:         publisher,
		connectionManager: connMgr,
	}
}

// Consume will call consume on all the subscribers
// For each subscriber it will create a new go routine and will wait on it for incoming messages
func (c *client) Consume(ctx context.Context) {
	conn := *c.connectionManager
	for _, sub := range *c.subscribers {
		ch, err := conn.GetChannel()
		messages, err := sub.Consume(ch)

		if err == nil {
			go func() {
				for msg := range messages {
					//TODO: replace with logger
					fmt.Printf("Received message: %s", msg.Body)
					sub.KitSubscriber.ServeDelivery(sub.Channel)(&msg)
				}
			}()
		}
	}
}

// Publish will call the publisher to publish the message
func (c *client) Publish(ctx context.Context, message *Message, exchangeName string) error {
	p := *c.publisher
	return p.Publish(ctx, message, exchangeName)
}

// Close will call the close functions on the publisher and subscribers
// Basically it will close all the open connections and channels
// This must be called before the application terminate to prevent connection or channel leaks
func (c *client) Close(ctx context.Context) error {
	var err error

	// closing the publisher channel
	p := *c.publisher
	perr := p.Close(ctx)
	if perr != nil {
		err = errors.NewApplicationError(perr, "failed to close rabbit publisher")
	}

	// closing the subscribers channels
	subs := *c.subscribers
	if subs != nil && len(subs) > 0 {
		for _, sub := range subs {
			suberr := sub.Close(ctx)
			if suberr != nil {
				err = errors.Annotate(err, suberr.Error())
			}
		}
	}

	// closing the connection
	conn := *c.connectionManager
	cerr := conn.CloseConnection(ctx)
	if cerr != nil {
		err = errors.Annotate(err, cerr.Error())
	}
	return err
}
