package rabbitmq_test

import (
	"testing"

	tlectx "github.com/kolbis/corego/context"
	tlelogger "github.com/kolbis/corego/logger"
	"github.com/kolbis/corego/transport/rabbitmq"
)

type loggedInCommandData struct {
	ID   int
	Name string
}

const (
	exchangeName string = "exchange1"
	username     string = "thelotter"
	pwd          string = "Dhvbuo1"
	host         string = "int-k8s1"
	vhost        string = "thelotter"
	port         int    = 32672
)

func TestPublishMessage(t *testing.T) {
	ctx := tlectx.Root()
	req := loggedInCommandData{ID: 1, Name: "guy kolbis"}
	message := rabbitmq.NewMessage(req, "thelotter.userloggedin")

	logManager := tlelogger.NewNopLogger()
	connInfo := rabbitmq.NewConnectionInfo(host, port, username, pwd, vhost)
	conn := rabbitmq.NewConnectionManager(connInfo)
	publisher := rabbitmq.NewPublisher(&conn)
	client := rabbitmq.NewClient(&conn, logManager, &publisher, nil)
	err := client.Publish(ctx, &message, exchangeName)

	if err != nil {
		t.Error(err)
	}
}
