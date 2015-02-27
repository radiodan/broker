package broker

import (
	zmq "github.com/pebbe/zmq4"
	"github.com/radiodan/broker/services"
	"log"
	"time"
)

func (b *Broker) poll() {
	poller := zmq.NewPoller()
	poller.Add(b.socket, zmq.POLLIN)

	serviceDirectory := services.NewServiceDirectory()

	for {
		polled, err := poller.Poll(time.Second * 10)

		if err != nil {
			log.Println("E: Interrupted")
			break //  Interrupted
		}

		if len(polled) > 0 {
			msg, err := b.socket.RecvMessage(1)

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
				serviceType := data[0]
				for _, serviceInstance := range data[1:] {
					log.Printf("?: %q\n", serviceInstance)
					serviceDirectory.AddWorker(sender, serviceType, serviceInstance)
				}
			}

			if protocol == "MDPW02" && command == "2" {
				log.Printf("I: %s replying\n", sender)
				correlationID := data[1]
				response := data[2:]
				b.socket.SendMessage(data[0], correlationID, "SUCCESS", response)
			}

			if protocol == "MDPC02" && command == "1" {
				log.Printf("I: %s is a client\n", sender)
				log.Printf("I: data - %q", data)

				//correlationId := data[0]
				serviceType := data[1]
				serviceInstance := data[2]
				//msg := data[3:]

				worker, err := serviceDirectory.WorkerForService(serviceType, serviceInstance)

				if err != nil {
					log.Printf("I: No worker for %s.%s", serviceType, serviceInstance)
					return
				}
				log.Printf("I: sending data to worker %s", worker.Name)
				msgCount, err := b.socket.SendMessage(worker.Identity, sender, data)
				if err != nil {
					log.Printf("! %x", err)
				} else {
					log.Printf("I: sent %i bytes", msgCount)
				}
			}

			//log.Printf("I: received message: %q\n", msg)
			//log.Printf("I: sender: %q\n", sender)
			//log.Printf("I: protocol: %q\n", protocol)
			//log.Printf("I: command: %q\n", command)
			//log.Printf("I: data: %q\n", data)
		}
	}
}
