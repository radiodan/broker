package main

import (
	"github.com/radiodan/broker/service"
	"log"
)

func main() {
	serviceLocation := "tcp://127.0.0.1:7171"
	serviceBroker := service.New(serviceLocation)

	go serviceBroker.Poll()

	log.Printf("Listening on %s", serviceLocation)

	// cheap trick to keep the main thread running
	forever := make(chan bool)
	<-forever
}
