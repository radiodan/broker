package service

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"strings"
	"time"
)

func (b *Broker) respondToWorker(msg *Message) {
	log := log.WithFields(
		log.Fields{"file": "service/respond-to-broker.go"},
	)

	worker, exists := b.Service.workers[msg.Sender]

	switch msg.Command {
	case COMMAND_READY:
		if exists == true {
			// TODO: what to respond with?
			// reset heartbeat expiry
			worker.Refresh()
			return
		}

		log.Debug(fmt.Sprintf("%s is a worker", msg.Sender))
		services, err := NewServiceMessage(msg.Payload)

		if err != nil {
			errString := fmt.Sprintf(
				"Services invalid for worker %s: %s",
				msg.Sender, strings.Join(msg.Payload, ","),
			)
			log.Warn(errString)
			b.DisconnectWorker(msg.Sender, errString)
			return
		}

		err = b.Service.AddWorker(msg.Sender, services)

		if err != nil {
			errString := fmt.Sprintf(
				"Failed to add worker %s: %s", msg.Sender, err,
			)
			log.Warn(errString)
			b.DisconnectWorker(msg.Sender, errString)
		}
	case COMMAND_REQUEST:
		log.Debug(fmt.Sprintf("%s replying", msg.Sender))
		correlationID := msg.Payload[1]
		response := msg.Payload[2:]
		r := []string{msg.Payload[0], PROTOCOL_WORKER, correlationID, "SUCCESS"}
		r = append(r, response...)

		b.Socket.SendMessage(r)

		worker, exists := b.Service.workers[msg.Sender]

		if exists != true {
			return
		}

		// reset heartbeat expiry
		worker.Refresh()

		req, messageWaiting := worker.NextMsg()

		if messageWaiting == true {
			b.Socket.SendMessage(req.Serialize())
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
		log.Warn(fmt.Sprintf("Unknown command %s", msg.Command))
	}
}
