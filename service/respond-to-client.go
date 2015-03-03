package service

import (
	"fmt"
	"log"
)

func (b *Broker) respondToClient(msg *Message) {
	switch msg.Command {
	case "1":
		log.Printf("I: %s is a client\n", msg.Sender)
		log.Printf("I: msg.Payload - %q", msg.Payload)

		//correlationId := msg.Payload[0]
		serviceType := msg.Payload[1]
		serviceInstance := msg.Payload[2]
		//payload := msg.Payload[3:]

		worker, err := b.Service.WorkerForService(serviceType, serviceInstance)

		if err != nil {
			errMsg := fmt.Sprintf("!: No worker for %s.%s", serviceType, serviceInstance)
			log.Printf(errMsg)
			return
		}

		log.Printf("I: sending data to worker %s", worker.Name)

		res := []string{
			worker.Identity,
			msg.Sender,
		}
		res = append(res, msg.Payload...)

		if worker.Ready {
			worker.Ready = false
			b.Socket.SendMessage(res)
		} else {
			worker.Queue = append([][]string{res}, worker.Queue...)
			log.Println("I: Appended msg for later processing")
		}
		return
	default:
		log.Printf("!: Unknown command %s", msg.Command)
	}
}