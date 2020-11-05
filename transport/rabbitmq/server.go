package rabbitmq

import (
	"context"

	tlelogger "github.com/kolbis/corego/logger"
	tletracer "github.com/kolbis/corego/tracer"
)

// Server responsibility to initiate the ability to consume messages
type Server interface {
	Run(context.Context) error
	Shutdown(context.Context)
}

type server struct {
	logger            tlelogger.Logger
	tracer            tletracer.Tracer
	connectionManager *ConnectionManager
	client            *Client
}

// NewServer will create a new instance of Server which can be executed to start and recieving messages
func NewServer(logger tlelogger.Logger, tracer tletracer.Tracer, rabbit *Client, conn *ConnectionManager) Server {
	return &server{
		client:            rabbit,
		logger:            logger,
		tracer:            tracer,
		connectionManager: conn,
	}
}

// Run will start all the listening on all the consumers
func (s server) Run(ctx context.Context) error {
	defer s.Shutdown(ctx)

	forever := make(chan bool)
	c := *s.client
	err := c.Consume(ctx)

	if err == nil {
		<-forever
	}

	return err
}

// Shutdown will close the server and call client to close resources
func (s server) Shutdown(ctx context.Context) {
	tlelogger.DebugWithContext(ctx, s.logger, "shutdown amqp server")
	c := *s.client
	c.Close(ctx)
}
