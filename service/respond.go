package service

import (
	"log"
)

func (b *Broker) Respond(msg *Message) {
	switch msg.Protocol {
	case "MDPW02":
		b.respondToWorker(msg)
	case "MDPC02":
		b.respondToClient(msg)
	default:
		log.Printf("Unknown protocol %s", msg.Protocol)
	}
}
