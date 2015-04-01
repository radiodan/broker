package service

import (
	log "github.com/Sirupsen/logrus"
	zmq "github.com/pebbe/zmq4"
	"runtime"
	"time"
)

func (b *Broker) Poll() {
	log := log.WithFields(
		log.Fields{"file": "service/poll.go"},
	)

	b.connect()
	runtime.SetFinalizer(b, (*Broker).Close)

	poller := zmq.NewPoller()
	poller.Add(b.Socket, zmq.POLLIN)

	for {
		polled, err := poller.Poll(HEARTBEAT_INTERVAL)

		if err != nil {
			log.Error("Interrupted")
			log.Error(err)
			break //  Interrupted
		}

		if len(polled) > 0 {
			msg, err := b.Socket.RecvMessage(0)

			if err != nil {
				log.Error("Interrupted")
				log.Error(err)
				break //  Interrupted
			}

			message, err := NewMessage(msg)

			if err != nil {
				log.Warn("Message malformed")
				continue
			}

			b.Respond(message)
		}

		if time.Now().After(b.HeartbeatAt) {
			log.Debug("Heartbeat")
			b.Purge()
			b.Heartbeat()
		}
	}
}
