package rabbitmq_test

import (
	"testing"

	"github.com/kolbis/corego/transport/rabbitmq"
)

type payload struct {
	id int
}

func TestNewMessage(t *testing.T) {
	message := rabbitmq.NewMessage(payload{id: 1}, "urn1", "urn2")
	data := message.Payload.Data.(payload)

	if data.id != 1 {
		t.Fail()
	}

	if message.MessageType[0] != "urn:message:urn1" {
		t.Fail()
	}

	if message.MessageType[1] != "urn:message:urn2" {
		t.Fail()
	}
}
