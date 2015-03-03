package main

import (
	"github.com/radiodan/broker/service"
	"log"
)

func main() {
	location := "tcp://127.0.0.1:7171"
	serviceBroker := service.New(location)

	serviceDirectory := service.NewServiceDirectory()
	messageHandler := service.NewMessageHandler(serviceDirectory)

	go serviceBroker.Poll(messageHandler)

	log.Printf("Listening on %s\n", location)

	// cheap trick to keep the main thread running
	forever := make(chan bool)
	<-forever
}
