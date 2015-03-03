package service

import (
	zmq "github.com/pebbe/zmq4"
	"log"
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

func (m *MessageHandler) Respond(msg *Message, socket *zmq.Socket) {
	switch msg.Protocol {
	case "MDPW02":
		m.workerHandler(msg, socket)
	case "MDPC02":
		m.clientHandler(msg, socket)
	default:
		log.Printf("Unknown protocol %s", msg.Protocol)
	}
}
