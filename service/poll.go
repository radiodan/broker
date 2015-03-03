package service

import (
	zmq "github.com/pebbe/zmq4"
	"log"
	"runtime"
	"time"
)

func (b *Broker) Poll(messageHandler *MessageHandler) {
	b.connect()
	runtime.SetFinalizer(b, (*Broker).Close)

	poller := zmq.NewPoller()
	poller.Add(b.socket, zmq.POLLIN)

	for {
		polled, err := poller.Poll(HEARTBEAT_INTERVAL)

		if err != nil {
			log.Println("E: Interrupted")
			log.Printf("%q\n", err)
			break //  Interrupted
		}

		if len(polled) > 0 {
			msg, err := b.socket.RecvMessage(0)

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

			messageHandler.Respond(message, b.socket)
		}

		// TODO: heartbeat
		if time.Now().After(b.heartbeatAt) {
			log.Println("I: Heartbeat")
			//broker.Purge()
			//for _, worker := range b.waiting {
			//	worker.Send(WORKER_HEARTBEAT, []string{})
			//}
			b.heartbeatAt = time.Now().Add(HEARTBEAT_INTERVAL)
		}
	}
}
