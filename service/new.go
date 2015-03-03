package service

import (
	zmq "github.com/pebbe/zmq4"
	"time"
)

type Broker struct {
	Endpoint    string
	Socket      *zmq.Socket
	HeartbeatAt time.Time
	Service     *ServiceDirectory
}

func New(endpoint string) (b Broker) {
	b.Endpoint = endpoint
	b.HeartbeatAt = time.Now().Add(HEARTBEAT_INTERVAL)
	b.Service = NewServiceDirectory()
	return
}
