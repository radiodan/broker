package broker

import (
	zmq "github.com/pebbe/zmq4"
	"log"
	"runtime"
	"time"
)

func (b *Broker) Poll(messageHandler *MessageHandler) {
	b.Connect()
	runtime.SetFinalizer(b, (*Broker).Close)

	poller := zmq.NewPoller()
	poller.Add(b.socket, zmq.POLLIN)

	for {
		polled, err := poller.Poll(time.Second * 10)

		if err != nil {
			log.Println("E: Interrupted")
			log.Printf("%q\n", err)
			break //  Interrupted
		}

		if len(polled) <= 0 {
			// TODO: heartbeat
			continue
		}

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

		responses, err := messageHandler.Respond(message)

		if err != nil {
			log.Println("!: Could not respond")
			continue
		}

		if len(responses) == 0 {
			log.Println("!: Will not respond")
			continue
		}

		for _, response := range responses {
			msgCount, err := b.socket.SendMessage(response)

			if err != nil {
				log.Printf("! %x\n", err)
			} else {
				log.Printf("I: sent %i bytes\n", msgCount)
			}
		}
	}
}
