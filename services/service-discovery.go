package services

import (
	"errors"
	"fmt"
	"log"
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
	Name     string // human readable id
	Identity string // routing frame
}

func NewServiceDirectory() ServiceDirectory {
	return ServiceDirectory{
		services: make(map[string]*ServiceType),
		workers:  make(map[string]*Worker),
		index:    make(map[string][]string),
	}
}

func (serviceDirectory *ServiceDirectory) AddWorker(identity string, serviceTypeName string, serviceInstanceName string) {
	// create or return worker
	worker, exists := serviceDirectory.workers[identity]

	if exists == false {
		name := fmt.Sprintf("%q", identity)
		worker = &Worker{
			Identity: identity,
			Name:     name,
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
