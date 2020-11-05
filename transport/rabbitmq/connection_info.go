package rabbitmq

import (
	"fmt"
)

// ConnectionInfo is a container for rabbit connection information
type ConnectionInfo struct {
	// ConnectionString like amqp://guest:guest@localhost:5672/
	// It will be build when calling the NewConnectionInfo
	ConnectionString string

	// Usewrname to connect to RabbitMQ
	Username string

	// Pwd to connect to RabbitMQ
	Pwd string

	// VirtualHost to connect to RabbitMQ
	VirtualHost string

	// Port to connect to RabbitMQ
	Port int

	// Host to connect to RabbitMQ
	Host string
}

// NewConnectionInfo will create a new instance of the connection information and will build the rabbitmq connection string
func NewConnectionInfo(host string, port int, username string, password string, vhost string) ConnectionInfo {
	url := fmt.Sprintf("amqp://%s:%s@%s:%d/%s", username, password, host, port, vhost)
	return ConnectionInfo{
		ConnectionString: url,
		Host:             host,
		VirtualHost:      vhost,
		Pwd:              password,
		Username:         username,
		Port:             port,
	}
}
