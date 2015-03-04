package service

import (
	"log"
)

func (b *Broker) Respond(msg *Message) {
	switch msg.Protocol {
	case PROTOCOL_WORKER:
		b.respondToWorker(msg)
	case PROTOCOL_CLIENT:
		b.respondToClient(msg)
	default:
		log.Printf("Unknown protocol %s", msg.Protocol)
	}
}
