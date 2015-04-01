package main

import (
	"flag"
	"fmt"
	"github.com/radiodan/broker/pubsub"
	"github.com/radiodan/broker/service"
	"log"
	"path"
)

func main() {
	serviceLocation, pubLocation, subLocation := parseFlags()

	serviceBroker := service.New(serviceLocation)
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

func parseFlags() (service string, pub string, sub string) {
	servicePort := flag.Int("service-port", 7171, "Port for service")
	pubPort := flag.Int("pub-port", 7172, "Port for publishing")
	subPort := flag.Int("sub-port", 7173, "Port for subscribing")

	serviceSocket := flag.String("service-socket", "", "Socket path for service")
	pubSocket := flag.String("pub-socket", "", "Socket path for publishing")
	subSocket := flag.String("sub-socket", "", "Socket path for subscribing")

	flag.Parse()

	service = connectionPath(servicePort, serviceSocket)
	pub = connectionPath(pubPort, pubSocket)
	sub = connectionPath(subPort, subSocket)

	return
}

func connectionPath(port *int, socket *string) (fullPath string) {
	switch {
	case *socket == "":
		fullPath = fmt.Sprintf("tcp://0.0.0.0:%v", *port)
	default:
		socketPath := path.Clean(*socket)
		fullPath = fmt.Sprintf("ipc://%s", socketPath)
	}

	return
}
