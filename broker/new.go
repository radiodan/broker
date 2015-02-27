package broker

import (
	zmq "github.com/pebbe/zmq4"
)

type Broker struct {
	endpoint string
	socket   *zmq.Socket
}

func New(endpoint string) (b Broker) {
	b.endpoint = endpoint
	return
}
