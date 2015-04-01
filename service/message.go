package service

import (
	"errors"
	log "github.com/Sirupsen/logrus"
)

type Message struct {
	Sender          string
	Protocol        string
	Command         string
	CorrelationId   string
	ServiceType     string
	ServiceInstance string
	Payload         []string
}

func NewMessage(params []string) (*Message, error) {
	var err error

	if len(params) < 3 {
		err = errors.New("Invalid message format")
		return nil, err
	}

	m := &Message{
		Sender:   params[0],
		Protocol: params[1],
		Command:  params[2],
	}

	if len(params) > 3 {
		m.Payload = params[3:]
	}

	log.Printf("%q", m)

	return m, err
}
