package service

import (
	"fmt"
	"log"
	"strings"
	"time"
)

func (b *Broker) respondToWorker(msg *Message) {
	worker, exists := b.Service.workers[msg.Sender]

	switch msg.Command {
	case COMMAND_READY:
		if exists == true {
			// TODO: what to respond with?
			// reset heartbeat expiry
			worker.Refresh()
			return
		}

		log.Printf("I: %s is a worker", msg.Sender)
		services, err := NewServiceMessage(msg.Payload)

		if err != nil {
			errString := fmt.Sprintf(
				"Services invalid for worker %s: %s",
				msg.Sender, strings.Join(msg.Payload, ","),
			)
			log.Println("!: " + errString)
			b.DisconnectWorker(msg.Sender, errString)
			return
		}

		err = b.Service.AddWorker(msg.Sender, services)

		if err != nil {
			errString := fmt.Sprintf(
				"Failed to add worker %s: %s", msg.Sender, err,
			)
			log.Println("!: " + errString)
			b.DisconnectWorker(msg.Sender, errString)
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
	case COMMAND_DISCONNECT:
		worker, exists := b.Service.workers[msg.Sender]

		if exists != true {
			return
		}

		// set expiry to now, will be cleaned up during purge
		worker.Expiry = time.Now()
	default:
		log.Printf("!: Unknown command %s", msg.Command)
	}
}
