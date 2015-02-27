package main

import (
	"errors"
	"fmt"
	zmq "github.com/pebbe/zmq4"
	"log"
	"time"
)

type ServiceDirectory struct {
	services map[string]*ServiceType
	index    map[string][]string
	workers  map[string]*Worker
}

type ServiceType struct {
	name      string
	instances map[string]*Worker
}

type Worker struct {
	name     string // human readable id
	identity string // routing frame
}

func (serviceDirectory *ServiceDirectory) AddWorker(identity string, serviceTypeName string, serviceInstanceName string) {
	// create or return worker
	worker, exists := serviceDirectory.workers[identity]

	if exists == false {
		name := fmt.Sprintf("%q", identity)
		worker = &Worker{
			identity: identity,
			name:     name,
		}

		serviceDirectory.workers[identity] = worker
	}

	log.Printf("?: Worker - %q", worker)

	// create or return serviceType
	serviceType, exists := serviceDirectory.services[serviceTypeName]

	if exists == false {
		serviceType = &ServiceType{
			name:      serviceTypeName,
			instances: make(map[string]*Worker),
		}

		serviceDirectory.services[serviceTypeName] = serviceType
	}

	// register serviceInstance, or err if already found
	_, exists = serviceType.instances[serviceInstanceName]

	if exists == false {
		log.Printf("Adding worker %q to %s.%s", worker, serviceTypeName, serviceInstanceName)
		serviceType.instances[serviceInstanceName] = worker
	}

	// add to index
	serviceDirectory.index[serviceTypeName] = append(serviceDirectory.index[serviceTypeName], serviceInstanceName)
}

func (serviceDirectory *ServiceDirectory) WorkerForService(serviceTypeName string, serviceInstanceName string) (serviceWorker *Worker, err error) {
	// check for serviceType
	serviceType, exists := serviceDirectory.index[serviceTypeName]

	if !exists {
		err = errors.New("Unknown serviceType")
		return
	}

	log.Printf("Found serviceType %q", serviceType)

	for _, serviceInstance := range serviceType {
		if serviceInstance == serviceInstanceName {
			log.Printf("Match: %q", serviceDirectory.services[serviceTypeName])
			serviceWorker, exists = serviceDirectory.services[serviceTypeName].instances[serviceInstanceName]
			if exists {
				log.Printf("Found serviceWorker: %q", serviceWorker)
				return
			}
		}
	}

	err = errors.New("Unknown serviceInstance")
	return
}

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

	serviceDirectory := ServiceDirectory{
		services: make(map[string]*ServiceType),
		workers:  make(map[string]*Worker),
		index:    make(map[string][]string),
	}

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
				serviceType := data[0]
				for _, serviceInstance := range data[1:] {
					log.Printf("?: %q\n", serviceInstance)
					serviceDirectory.AddWorker(sender, serviceType, serviceInstance)
				}
				log.Printf("I: services - %q\n", serviceDirectory.index)
			}

			if protocol == "MDPW02" && command == "2" {
				log.Printf("I: %s replying\n", sender)
				correlationID := data[1]
				response := data[2:]
				socket.SendMessage(data[0], correlationID, "SUCCESS", response)
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
				log.Printf("I: sending data to worker %s", worker.name)
				msgCount, err := socket.SendMessage(worker.identity, sender, data)
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
