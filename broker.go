package main

import (
	"github.com/radiodan/broker/broker"
	"log"
)

func main() {
	location := "tcp://127.0.0.1:7171"
	brokerServer := broker.New(location)

	serviceDirectory := broker.NewServiceDirectory()
	messageHandler := broker.NewMessageHandler(serviceDirectory)

	log.Printf("Listening on %s\n", location)
	brokerServer.Poll(messageHandler)
}
