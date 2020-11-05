package rabbitmq

import (
	"context"

	"github.com/streadway/amqp"
	"github.com/kolbis/corego/errors"
)

// ConnectionManager responsibility is to manage rabbit connecations
type ConnectionManager interface {
	GetConnection() (*amqp.Connection, error)
	GetChannel() (*amqp.Channel, error)
	CloseConnection(context.Context) error
	CloseChannel(context.Context, *amqp.Channel) error
}

// NewConnectionManager will create a new ConnectionManager
// Best practice is to have a single instance of it and reuse it, in order to save connections
func NewConnectionManager(connInfo ConnectionInfo) ConnectionManager {
	c := connmgr{
		connectionInfo: connInfo,
		isConnected:    false,
	}

	return &c
}

type connmgr struct {
	connectionInfo ConnectionInfo
	connection     *amqp.Connection
	isConnected    bool
}

// GetConnection will return the single connection to rabbit
// If a connection does not exist it will create a new one and return it
func (c *connmgr) GetConnection() (*amqp.Connection, error) {
	var err error
	if c.isConnected == false {
		err = c.connect()
	}
	return c.connection, err
}

// GetChannel will use the single connection to create a channel and return it
// It is up to the caller to make sure that the channel is closed to prevent channels leak
func (c *connmgr) GetChannel() (*amqp.Channel, error) {
	var ch *amqp.Channel

	conn, err := c.GetConnection()
	if err == nil {
		ch, err = conn.Channel()
		if err != nil {
			return ch, errors.NewApplicationErrorf("faild to open a channel", err.Error())
		}
	}
	return ch, err
}

// CloseChannel will saifly close the channel
func (c *connmgr) CloseChannel(ctx context.Context, ch *amqp.Channel) error {
	if ch != nil {
		err := ch.Close()
		if err != nil {
			return errors.NewApplicationErrorf("failed to close rabbit channel", err.Error())
		}
	}
	return nil
}

// CloseConnection will saifly close the single connection
func (c *connmgr) CloseConnection(ctx context.Context) error {
	var err error

	if c.isConnected {
		connerr := c.connection.Close()
		if connerr != nil {
			err = errors.NewApplicationErrorf("failed to close rabbit connection %s", connerr.Error())
		} else {
			c.isConnected = false
		}
	}

	return err
}

func (c *connmgr) connect() error {
	conn, err := amqp.Dial(c.connectionInfo.ConnectionString)
	if err != nil {
		return errors.NewApplicationErrorf("failed to connect to rabbit", err.Error())
	}
	c.connection = conn
	c.isConnected = true
	return nil
}
