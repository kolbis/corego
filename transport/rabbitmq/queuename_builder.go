package rabbitmq

import (
	"fmt"
	"strings"

	"github.com/kolbis/corego/utils"
)

const (
	commandSuffix string = "-command"
)

// BuildPrivateQueueName will take a queue name and will adjust it to the private queue standard
// Private queues are temporary queues which will be auto deleted
func BuildPrivateQueueName(queueName string) string {
	name := fmt.Sprintf("%s-private-%s", queueName, utils.NewUUID())
	return name
}

// BuildCommandQueueName will take a queue name and will adjust it to the command queue standard
// Command queues are durable persisted queues
func BuildCommandQueueName(queueName string) string {
	if strings.HasSuffix(queueName, commandSuffix) == false {

		name := queueName + commandSuffix
		return name
	}
	return queueName
}
