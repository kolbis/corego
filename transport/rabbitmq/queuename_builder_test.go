package rabbitmq_test

import (
	"strings"
	"testing"

	"github.com/kolbis/corego/transport/rabbitmq"
)

func TestCommandQueueName(t *testing.T) {
	name := "UserLoggedIn"
	want := "UserLoggedIn-command"
	is := rabbitmq.BuildCommandQueueName(name)

	if is != want {
		t.Fail()
	}
}

func TestCommandQueueNameSuffixExist(t *testing.T) {
	name := "UserLoggedIn-command"
	want := "UserLoggedIn-command"
	is := rabbitmq.BuildCommandQueueName(name)

	if is != want {
		t.Fail()
	}
}

func TestPrivateQueueName(t *testing.T) {
	name := "UserLoggedIn"
	want := "UserLoggedIn-private-"
	is := rabbitmq.BuildPrivateQueueName(name)

	if strings.HasPrefix(is, want) == false {
		t.Fail()
	}
}
