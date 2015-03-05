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
			for _, m := range worker.Queue {
				fmt.Printf("I: Cancelling queued message %v", m)
				b.Socket.SendMessage(m[0], m[1], "FAIL")
			}
			b.DisconnectWorker(worker.Identity, "Heartbeat Timeout")
		}
	}
}
