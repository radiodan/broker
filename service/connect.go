package service

import (
	log "github.com/Sirupsen/logrus"
	zmq "github.com/pebbe/zmq4"
)

func (b *Broker) connect() {
	// create connection Socket
	context, err := zmq.NewContext()

	b.Socket, err = context.NewSocket(zmq.ROUTER)

	if err != nil {
		log.Fatal("Could not start broker: %v\n", err)
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
