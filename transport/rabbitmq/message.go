package rabbitmq

import (
	"fmt"
)

// MessagePayload is the basic unit used to send and received data thought rabbitmq
// or more information about the URN http://masstransit-project.com/MassTransit/architecture/interoperability.html
type MessagePayload struct {
	Data           interface{}            `json:"data,omitempty"`
	AdditionalData map[string]interface{} `json:"additionalData"`
}

// Message is the Message classes used by masstransit for every message
type Message struct {
	MessageType []string        `json:"messageType"`
	Payload     *MessagePayload `json:"message"`
}

// NewMessage will create a rabbit transport message
// It is expected for all messages to be published and consumed from and to rabbit, to be a Message
func NewMessage(payloadData interface{}, urn ...string) Message {
	var urnSlice = buildURN(urn...)
	payload := MessagePayload{}
	payload.Data = payloadData

	message := Message{
		Payload:     &payload,
		MessageType: urnSlice,
	}

	return message
}

func buildURN(urn ...string) []string {
	var urnSlice = make([]string, 0)
	for _, u := range urn {
		urnSlice = append(urnSlice, fmt.Sprintf("urn:message:%v", u))
	}
	return urnSlice
}
