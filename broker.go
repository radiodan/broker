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

	go brokerServer.Poll(messageHandler)

	log.Printf("Listening on %s\n", location)

	// cheap trick to keep the main thread running
	forever := make(chan bool)
	<-forever
}
