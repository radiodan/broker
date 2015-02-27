package main

import (
	"github.com/radiodan/broker/broker"
	"log"
)

func main() {
	location := "tcp://127.0.0.1:7171"
	brokerServer := broker.New(location)

	log.Printf("Listening on %s\n", location)
	brokerServer.Connect()
}
