package service

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"time"
)

func (b *Broker) Purge() {
	log := log.WithFields(
		log.Fields{"file": "service/purge.go"},
	)

	now := time.Now()

	for name, worker := range b.Service.workers {
		log.Debug(fmt.Sprintf("Looking at worker %s", name))

		if worker.Expiry.Before(now) {
			log.Info(fmt.Sprintf("Removing worker %s", name))

			b.Service.RemoveWorker(worker)
			b.DisconnectWorker(worker.Identity, "Heartbeat Timeout")

			for _, req := range worker.Queue {
				log.Debug(fmt.Sprintf("Cancelling queued message %v\n", req))
				errMsg := fmt.Sprintf("Worker %s is no longer reachable", name)

				b.Socket.SendMessage(
					req.Message.Sender,
					PROTOCOL_BROKER,
					req.Message.CorrelationId,
					"FAIL",
					errMsg,
				)
			}
		}
	}
}
