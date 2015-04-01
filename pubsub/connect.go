package pubsub

import (
	log "github.com/Sirupsen/logrus"
	zmq "github.com/pebbe/zmq4"
)

func (b *Broker) connect() {
	// create connection Socket
	context, err := zmq.NewContext()

	b.PubSocket, err = context.NewSocket(zmq.XPUB)

	if err != nil {
		log.Fatal("Could not start broker: %v\n", err)
		return
	}

	b.SubSocket, err = context.NewSocket(zmq.SUB)

	if err != nil {
		log.Fatal("Could not start broker: %v\n", err)
		return
	}

	b.PubSocket.Bind(b.PubEndpoint)
	b.SubSocket.Bind(b.SubEndpoint)
}

func (b *Broker) Close() {
	b.CloseSocket(b.SubSocket)
	b.CloseSocket(b.PubSocket)
}

func (b *Broker) CloseSocket(socket *zmq.Socket) {
	if socket != nil {
		socket.Close()
		socket = nil
	}
}
