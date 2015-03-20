package service

import (
	"time"
)

func (b *Broker) Heartbeat() {
	for _, worker := range b.Service.workers {
		b.Socket.SendMessage(worker.Identity, PROTOCOL_BROKER, COMMAND_HEARTBEAT)
	}
	b.HeartbeatAt = time.Now().Add(HEARTBEAT_INTERVAL)
}
