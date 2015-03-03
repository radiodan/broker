package service

import (
	"log"
	"time"
)

func (b *Broker) Purge() (deletedWorkers []string) {
	now := time.Now()

	for name, worker := range b.Service.workers {
		log.Printf("I: Looking at worker %s", name)

		if worker.Expiry.Before(now) {
			log.Printf("I: Removing worker %s", name)
			// TODO: send real disconnect message
			b.Socket.SendMessage(
				worker.Identity, "MDPW_DISCONNECT", "", []string{},
			)
			deletedWorkers = append(deletedWorkers, name)
			delete(b.Service.workers, name)
			// TODO: remove from serviceDirectory
		}
	}

	return
}
