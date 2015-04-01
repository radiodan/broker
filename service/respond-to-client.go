package service

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
)

func (b *Broker) respondToClient(msg *Message) {
	var err error
	var worker *Worker

	switch msg.Command {
	case COMMAND_READY:
		log.Printf("I: %s is a client\n", msg.Sender)
		log.Printf("I: msg.Payload - %q", msg.Payload)

		if len(msg.Payload) < 3 {
			errMsg := "Malformed command"
			log.Printf("!: %s", errMsg)

			b.Socket.SendMessage(msg.Sender, PROTOCOL_BROKER, msg.Payload[0], "FAIL", errMsg)
		}

		msg.CorrelationId = msg.Payload[0]
		msg.ServiceType = msg.Payload[1]
		msg.ServiceInstance = msg.Payload[2]

		if len(msg.Payload) > 3 {
			msg.Payload = msg.Payload[3:]
		} else {
			msg.Payload = []string{}
		}

		if msg.ServiceType == "broker" {
			// reply
			err = b.ReplyForService(msg)
		} else {
			worker, err = b.Service.WorkerForService(msg.ServiceType, msg.ServiceInstance)
		}

		if err != nil {
			errMsg := fmt.Sprintf("No worker for %s.%s", msg.ServiceType, msg.ServiceInstance)
			log.Printf("!: %s", errMsg)

			b.Socket.SendMessage(msg.Sender, PROTOCOL_BROKER, msg.CorrelationId, "FAIL", errMsg)
			return
		}

		if worker == nil {
			return
		}

		log.Printf("I: sending data to worker %s", worker.Name)

		req := NewRequest(worker.Identity, msg)

		if worker.Ready == true {
			log.Printf("I: Send REQ %q", req)
			b.Socket.SendMessage(req.Serialize())
			worker.Ready = false
		} else {
			worker.AppendToQueue(req)
			log.Println("I: Appended msg for later processing")
		}
		return
	default:
		log.Printf("!: Unknown command %s", msg.Command)
	}
}
