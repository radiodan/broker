package service

import (
	log "github.com/Sirupsen/logrus"
)

func (b *Broker) Respond(msg *Message) {
	log := log.WithFields(
		log.Fields{"file": "service/respond.go"},
	)

	switch msg.Protocol {
	case PROTOCOL_WORKER:
		b.respondToWorker(msg)
	case PROTOCOL_CLIENT:
		b.respondToClient(msg)
	default:
		log.Error("Unknown protocol %s", msg.Protocol)
	}
}
