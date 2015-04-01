package service

import (
	log "github.com/Sirupsen/logrus"
	zmq "github.com/pebbe/zmq4"
	"runtime"
	"time"
)

func (b *Broker) Poll() {
	b.connect()
	runtime.SetFinalizer(b, (*Broker).Close)

	poller := zmq.NewPoller()
	poller.Add(b.Socket, zmq.POLLIN)

	for {
		polled, err := poller.Poll(HEARTBEAT_INTERVAL)

		if err != nil {
			log.Println("E: Interrupted")
			log.Printf("E: %q", err)
			break //  Interrupted
		}

		if len(polled) > 0 {
			msg, err := b.Socket.RecvMessage(0)

			if err != nil {
				log.Println("E: Interrupted")
				log.Printf("%q\n", err)
				break //  Interrupted
			}

			message, err := NewMessage(msg)

			if err != nil {
				log.Println("!: Message malformed")
				continue
			}

			b.Respond(message)
		}

		if time.Now().After(b.HeartbeatAt) {
			log.Println("I: Heartbeat")
			b.Purge()
			b.Heartbeat()
		}
	}
}
