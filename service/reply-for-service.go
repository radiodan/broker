package service

import (
	"encoding/json"
	"errors"
	"log"
)

func (b *Broker) ReplyForService(msg *Message) (err error) {
	switch msg.ServiceInstance {
	case "discovery":
		indexJSON, err := json.Marshal(b.Service.index)

		if err != nil {
			log.Printf("%q", err)
			return errors.New("Error generating service directory")
		}

		b.Socket.SendMessage(msg.Sender, msg.CorrelationId, "SUCCESS", indexJSON)
	case "service":
		var responseType string
		var exists bool

		if len(msg.Payload) >= 2 {
			log.Printf("%q", msg.Payload)
			serviceType := msg.Payload[0]
			serviceInstance := msg.Payload[1]

			exists = b.Service.ServiceExists(serviceType, serviceInstance)
		} else {
			exists = false
		}

		if exists == true {
			responseType = "SUCCESS"
		} else {
			responseType = "FAIL"
		}

		b.Socket.SendMessage(msg.Sender, msg.CorrelationId, responseType)
	default:
		err = errors.New("broker." + msg.ServiceInstance + " not found")
	}
	return
}
