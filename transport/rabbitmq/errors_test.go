package rabbitmq_test

import (
	"context"
	"testing"

	"github.com/kolbis/corego/transport/rabbitmq"
)

func TestRabbitError(t *testing.T) {
	err := rabbitmq.NewRabbitErrorf(123, "err")

	if rabbitmq.IsRabbitError(err) == false {
		t.Fail()
	}
	rerr := err.(*rabbitmq.Error)
	if rerr.ErrorCode != 123 {
		t.Fail()
	}
}

func TestNack(t *testing.T) {
	ctx := context.Background()
	nack := rabbitmq.ShouldNack(ctx)

	if nack != false {
		t.Fail()
	}

	ctx = rabbitmq.SetNack(ctx)

	nack = rabbitmq.ShouldNack(ctx)

	if nack != true {
		t.Fail()
	}
}
