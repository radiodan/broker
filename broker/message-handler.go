package broker

import (
	"errors"
	"fmt"
)

type MessageHandler struct {
	Service *ServiceDirectory
}

func NewMessageHandler(service *ServiceDirectory) *MessageHandler {
	m := &MessageHandler{
		Service: service,
	}

	return m
}

func (m *MessageHandler) Respond(msg *Message) ([][]string, error) {
	switch msg.Protocol {
	case "MDPW02":
		return m.workerHandler(msg)
	case "MDPC02":
		return m.clientHandler(msg)
	default:
		errString := fmt.Sprintf("Unknown protocol %s", msg.Protocol)
		return make([][]string, 0), errors.New(errString)
	}
}
