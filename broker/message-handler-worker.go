package broker

import (
	"fmt"
	"log"
)

func (m *MessageHandler) workerHandler(msg *Message, channel chan []string) {
	switch msg.Command {
	case "1":
		log.Printf("I: %s is a worker\n", sender)
		serviceType := data[0]
		for _, serviceInstance := range data[1:] {
			log.Printf("?: %q\n", serviceInstance)
			serviceDirectory.AddWorker(sender, serviceType, serviceInstance)
		}
		return
	case "2":
		log.Printf("I: %s replying\n", sender)
		correlationID := data[1]
		response := data[2:]
		r = []string{data[0], correlationID, "SUCCESS"}
		r = append(r, response...)
	default:
		log.Printf("!: Unknown command %s", msg.Command)
	}
}
