package service

import (
	"errors"
	"fmt"
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
	log := log.WithFields(
		log.Fields{"file": "service/message.go"},
	)

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

	log.Debug(fmt.Sprintf("New Message: %v", m))

	return m, err
}
