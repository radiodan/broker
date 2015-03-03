package service

import (
	zmq "github.com/pebbe/zmq4"
	"time"
)

type Broker struct {
	endpoint    string
	socket      *zmq.Socket
	heartbeatAt time.Time
}

func New(endpoint string) (b Broker) {
	b.endpoint = endpoint
	b.heartbeatAt = time.Now().Add(HEARTBEAT_INTERVAL)
	return
}
