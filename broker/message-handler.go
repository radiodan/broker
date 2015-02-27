package broker

import (
	"errors"
	"fmt"
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

func (m *MessageHandler) Respond(msg *Message) (r []string, err error) {
	if msg.Protocol == "MDPW02" && msg.Command == "1" {
		log.Printf("I: %s is a worker\n", msg.Sender)
		serviceType := msg.Payload[0]
		for _, serviceInstance := range msg.Payload[1:] {
			log.Printf("?: %q\n", serviceInstance)
			m.Service.AddWorker(msg.Sender, serviceType, serviceInstance)
		}

		// acknowledge request?
		return
	}

	if msg.Protocol == "MDPW02" && msg.Command == "2" {
		log.Printf("I: %s replying\n", msg.Sender)
		log.Printf("I: msg.Payload - %q", msg.Payload)
		correlationID := msg.Payload[1]
		response := msg.Payload[2:]
		r = []string{msg.Payload[0], correlationID, "SUCCESS"}
		r = append(r, response...)
		return
	}

	if msg.Protocol == "MDPC02" && msg.Command == "1" {
		log.Printf("I: %s is a client\n", msg.Sender)
		log.Printf("I: msg.Payload - %q", msg.Payload)

		//correlationId := msg.Payload[0]
		serviceType := msg.Payload[1]
		serviceInstance := msg.Payload[2]
		//payload := msg.Payload[3:]

		worker, err := m.Service.WorkerForService(serviceType, serviceInstance)

		if err != nil {
			errMsg := fmt.Sprintf("I: No worker for %s.%s", serviceType, serviceInstance)
			log.Printf(errMsg)
			err := errors.New(errMsg)
			return r, err
		}

		log.Printf("I: sending data to worker %s", worker.Name)
		r = []string{
			worker.Identity,
			msg.Sender,
		}

		r = append(r, msg.Payload...)
	}
	return
}
