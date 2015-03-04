package service

import (
	"log"
	"time"
)

func (b *Broker) Purge() {
	now := time.Now()

	for name, worker := range b.Service.workers {
		log.Printf("I: Looking at worker %s", name)

		if worker.Expiry.Before(now) {
			log.Printf("I: Removing worker %s", name)

			b.Service.RemoveWorker(worker)
			b.DisconnectWorker(worker.Identity, "Heartbeat Timeout")
		}
	}
}
