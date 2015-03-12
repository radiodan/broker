package pubsub

import (
	zmq "github.com/pebbe/zmq4"
)

type Broker struct {
	PubEndpoint string
	PubSocket   *zmq.Socket
	SubEndpoint string
	SubSocket   *zmq.Socket
}

func New(pubEndpoint string, subEndpoint string) (b Broker) {
	b.PubEndpoint = pubEndpoint
	b.SubEndpoint = subEndpoint

	return
}
