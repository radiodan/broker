package service

import (
	"fmt"
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

			for _, m := range worker.Queue {
				log.Printf("I: Cancelling queued message %v\n", m)
				errMsg := fmt.Sprintf("Worker %s is no longer reachable", name)

				//TODO: stop guessing message indexes
				b.Socket.SendMessage(m[2], m[3], "FAIL", errMsg)
			}
		}
	}
}
