package main

import (
	zmq "github.com/pebbe/zmq4"
	"log"
	"time"
)

func main() {
	// create connection socket
	context, err := zmq.NewContext()

	socket, err := context.NewSocket(zmq.ROUTER)
	defer socket.Close()

	location := "tcp://127.0.0.1:7171"

	socket.Bind(location)

	poller := zmq.NewPoller()
	poller.Add(socket, zmq.POLLIN)

	if err != nil {
		log.Printf("Could not start broker: %v\n", err)
		return
	}

	log.Println("Now listening at: " + location)

	workers := map[string]string{}

	for {
		polled, err := poller.Poll(time.Second * 10)

		if err != nil {
			log.Println("E: Interrupted")
			break //  Interrupted
		}

		if len(polled) > 0 {
			msg, err := socket.RecvMessage(0)

			if err != nil {
				log.Println("E: Interrupted")
				break //  Interrupted
			}

			sender := msg[0]
			protocol := msg[1]
			command := msg[2]
			data := msg[3:]

			if protocol == "MDPW02" && command == "1" {
				log.Printf("I: %s is a worker\n", sender)
				workers[protocol] = sender
				log.Printf("I: workers - %q\n", workers)
			}

			if protocol == "MDPW02" && command == "2" {
				log.Printf("I: %s replying\n", sender)
				correlationID := data[1]
				response := data[2:]
				socket.SendMessage(data[0], correlationID, "SUCCESS", response)
			}

			if protocol == "MDPC02" && command == "1" {
				log.Printf("I: %s is a client\n", sender)
				if len(workers) > 0 {
					// send message to worker, don't care about services yet
					log.Printf("I: sending data to worker %s", workers["MDPW02"])
					msgCount, err := socket.SendMessage(workers["MDPW02"], sender, data)
					if err != nil {
						log.Printf("! %x", err)
					} else {
						log.Printf("I: sent %i bytes", msgCount)
					}
				}
			}

			log.Printf("I: received message: %q\n", msg)
			log.Printf("I: sender: %q\n", sender)
			log.Printf("I: protocol: %q\n", protocol)
			log.Printf("I: command: %q\n", command)
			log.Printf("I: data: %q\n", data)
		}
	}
}
