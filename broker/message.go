package broker

import (
	"errors"
	"log"
)

type Message struct {
	Sender   string
	Protocol string
	Command  string
	Payload  []string
}

func NewMessage(params []string) (*Message, error) {
	var err error

	m := &Message{
		Sender:   params[0],
		Protocol: params[1],
		Command:  params[2],
		Payload:  params[3:],
	}

	if len(m.Sender) == 0 {
		err = errors.New("FAIL")
	}

	log.Printf("%q", m)

	return m, err
}
