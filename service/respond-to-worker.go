package service

import (
	"log"
)

func (b *Broker) respondToWorker(msg *Message) {
	worker, exists := b.Service.workers[msg.Sender]

	switch msg.Command {
	case COMMAND_READY:
		if exists == true {
			// TODO: send worker rejection
			// reset heartbeat expiry
			worker.Refresh()
			return
		}

		log.Printf("I: %s is a worker", msg.Sender)
		services, err := NewServiceMessage(msg.Payload)

		if err != nil {
			log.Printf(
				"!: Services invalid for worker %s: %v", msg.Sender, msg.Payload,
			)
			// TODO: send worker rejection
			return
		}

		err = b.Service.AddWorker(msg.Sender, services)

		if err != nil {
			log.Printf(
				"!: Failed to add worker %s: %s", msg.Sender, err,
			)
			// TODO: send worker rejection
		}
	case COMMAND_REQUEST:
		log.Printf("I: %s replying\n", msg.Sender)
		correlationID := msg.Payload[1]
		response := msg.Payload[2:]
		r := []string{msg.Payload[0], correlationID, "SUCCESS"}
		r = append(r, response...)

		b.Socket.SendMessage(r)

		worker, exists := b.Service.workers[msg.Sender]

		if exists != true {
			return
		}

		// reset heartbeat expiry
		worker.Refresh()

		r, messageWaiting := worker.NextMsg()

		if messageWaiting == true {
			b.Socket.SendMessage(r)
		} else {
			worker.Ready = true
		}
	case COMMAND_HEARTBEAT:
		// reset heartbeat expiry
		if exists == true {
			worker.Refresh()
		}
	default:
		log.Printf("!: Unknown command %s", msg.Command)
	}
}
