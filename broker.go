package main

import (
	"github.com/radiodan/broker/pubsub"
	"github.com/radiodan/broker/service"
	"log"
)

func main() {
	serviceLocation := "tcp://127.0.0.1:7171"
	serviceBroker := service.New(serviceLocation)

	pubLocation := "tcp://127.0.0.1:7172"
	subLocation := "tcp://127.0.0.1:7173"

	pubSubBroker := pubsub.New(pubLocation, subLocation)

	go serviceBroker.Poll()
	go pubSubBroker.Poll()

	log.Printf("Broker services on %s", serviceLocation)
	log.Printf("Broker publishes on %s", pubLocation)
	log.Printf("Broker subscribes on %s", subLocation)

	// cheap trick to keep the main thread running
	forever := make(chan bool)
	<-forever
}
