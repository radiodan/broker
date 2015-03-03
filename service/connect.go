package service

import (
	zmq "github.com/pebbe/zmq4"
	"log"
)

func (b *Broker) connect() {
	// create connection Socket
	context, err := zmq.NewContext()

	b.Socket, err = context.NewSocket(zmq.ROUTER)

	if err != nil {
		log.Printf("Could not start broker: %v\n", err)
		return
	}

	b.Socket.Bind(b.Endpoint)
}

func (b *Broker) Close() {
	if b.Socket != nil {
		b.Socket.Close()
		b.Socket = nil
	}
}
