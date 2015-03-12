package pubsub

import (
	zmq "github.com/pebbe/zmq4"
	"log"
	"runtime"
	"time"
)

func (b *Broker) Poll() {
	b.connect()
	runtime.SetFinalizer(b, (*Broker).Close)

	poller := zmq.NewPoller()
	poller.Add(b.PubSocket, zmq.POLLIN)
	poller.Add(b.SubSocket, zmq.POLLIN)

	b.SubSocket.SetSubscribe("")

	for {
		polled, err := poller.Poll(1000 * time.Millisecond)

		if err != nil {
			log.Printf("!: %q", err)
			break
		}

		for _, item := range polled {
			switch socket := item.Socket; socket {
			case b.SubSocket:
				msg, _ := b.SubSocket.RecvMessage(0)
				log.Printf("Topic: %s, Msg: %s", msg[0], msg[1])
				b.PubSocket.SendMessage(msg)
			case b.PubSocket:
				msg, _ := b.PubSocket.RecvMessage(0)

				frame := msg[0]
				topic := frame[1:]

				switch frame[0] {
				case 1:
					log.Printf("Subscribed: %s", topic)
				case 0:
					log.Printf("UnSubscribed: %s", topic)
				}
			}
		}
	}
}
