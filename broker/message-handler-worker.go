package broker

import (
	zmq "github.com/pebbe/zmq4"
	"log"
)

func (m *MessageHandler) workerHandler(msg *Message, socket *zmq.Socket) {
	switch msg.Command {
	case "1":
		log.Printf("I: %s is a worker", msg.Sender)
		services, err := NewServiceMessage(msg.Payload)

		if err != nil {
			log.Printf(
				"!: Services invalid for worker %s: %v", msg.Sender, msg.Payload,
			)
			// TODO: send worker rejection
			return
		}

		err = m.Service.AddWorker(msg.Sender, services)

		if err != nil {
			log.Printf(
				"!: Failed to add worker %s: %s", msg.Sender, err,
			)
			// TODO: send worker rejection
		}
		return
	case "2":
		log.Printf("I: %s replying\n", msg.Sender)
		correlationID := msg.Payload[1]
		response := msg.Payload[2:]
		r := []string{msg.Payload[0], correlationID, "SUCCESS"}
		r = append(r, response...)

		socket.SendMessage(r)

		worker, exists := m.Service.workers[msg.Sender]

		if exists != true {
			return
		}

		worker.Ready = true

		r, exists = worker.NextMsg()

		if exists == true {
			socket.SendMessage(r)
			worker.Ready = false
		}
	default:
		log.Printf("!: Unknown command %s", msg.Command)
	}
}
