package broker

import (
	zmq "github.com/pebbe/zmq4"
	"log"
)

func (b *Broker) Connect() {
	// create connection socket
	context, err := zmq.NewContext()

	b.socket, err = context.NewSocket(zmq.ROUTER)

	if err != nil {
		log.Printf("Could not start broker: %v\n", err)
		return
	}

	b.socket.Bind(b.endpoint)
}

func (b *Broker) Close() {
	b.socket.Close()
}
