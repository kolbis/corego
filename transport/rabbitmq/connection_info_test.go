package rabbitmq_test

import (
	"testing"

	"github.com/kolbis/corego/transport/rabbitmq"
)

func TestConnectionInfoURL(t *testing.T) {
	username := "user"
	pwd := "pwd"
	host := "localhost"
	vhost := "thelotter"
	port := 5672

	connectionMeta := rabbitmq.NewConnectionInfo(host, port, username, pwd, vhost)

	want := "amqp://user:pwd@localhost:5672/thelotter"
	is := connectionMeta.ConnectionString
	if is != want {
		t.Fail()
	}
}
