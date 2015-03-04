package service

import (
	"time"
)

func (b *Broker) Heartbeat() {
	for _, worker := range b.Service.workers {
		b.Socket.SendMessage(worker.Identity, COMMAND_HEARTBEAT, []string{})
	}
	b.HeartbeatAt = time.Now().Add(HEARTBEAT_INTERVAL)
}
