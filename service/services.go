package service

import (
	"errors"
	"fmt"
	"log"
	"strings"
)

type Services map[string][]string

type ServiceDirectory struct {
	services map[string]*ServiceType
	index    map[string][]string
	workers  map[string]*Worker
}

type ServiceType struct {
	name      string
	instances map[string]*Worker
}

func NewServiceMessage(serviceArray []string) (services Services, err error) {
	services = Services{}

	for _, service := range serviceArray {
		sSplit := strings.Split(service, ".")

		if len(sSplit) != 2 {
			err := errors.New("Invalid service list")
			return services, err
		}

		serviceType := sSplit[0]
		serviceInstance := sSplit[1]

		services[serviceType] = append(services[serviceType], serviceInstance)
	}

	return
}

func NewServiceDirectory() *ServiceDirectory {
	return &ServiceDirectory{
		services: make(map[string]*ServiceType),
		workers:  make(map[string]*Worker),
		index:    make(map[string][]string),
	}
}

func (serviceDirectory *ServiceDirectory) AddWorker(identity string, services Services) (err error) {
	// validate services
	name := fmt.Sprintf("%q", identity)
	invalidServices := serviceDirectory.validateServicesForWorker(name, services)

	if len(invalidServices) > 0 {
		errString := fmt.Sprintf(
			"Worker %s cannot register services %s",
			identity, strings.Join(invalidServices, ", "),
		)
		err = errors.New(errString)
		return err
	}

	// create or return worker
	worker, exists := serviceDirectory.workers[identity]

	if exists == false {
		worker = NewWorker(identity, services)

		serviceDirectory.workers[identity] = worker
	}

	for serviceTypeName, serviceInstances := range services {
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
		for _, serviceInstanceName := range serviceInstances {
			log.Printf("Adding worker %s to %s.%s",
				worker.Name, serviceTypeName, serviceInstanceName)
			serviceType.instances[serviceInstanceName] = worker
			// add to index
			serviceDirectory.index[serviceTypeName] = append(serviceDirectory.index[serviceTypeName], serviceInstanceName)
		}

	}
	log.Printf("?: Worker - %q", worker)

	return
}

func (serviceDirectory *ServiceDirectory) RemoveWorker(worker *Worker) {
	for sType, sInstances := range worker.Services {
		for _, sInstance := range sInstances {
			// copy the index instances array
			instances := append([]string(nil), serviceDirectory.index[sType]...)

			for i, inst := range serviceDirectory.index[sType] {
				if inst == sInstance {
					// remove matched instance from array
					instances = append(
						instances[:i],
						instances[i+1:]...,
					)
				}
			}

			if len(instances) == 0 {
				// remove entry from index
				delete(serviceDirectory.index, sType)
			} else {
				// replace index with new array
				serviceDirectory.index[sType] = instances
			}

			//remove worker reference from services
			delete(serviceDirectory.services[sType].instances, sInstance)
		}
	}

	// remove worker
	delete(serviceDirectory.workers, worker.Identity)
}

func (serviceDirectory *ServiceDirectory) validateServicesForWorker(name string, services Services) (invalidServices []string) {
	for sType, sInstances := range services {
		// if the serviceType doesnt exist, all the instances can be registered
		_, exists := serviceDirectory.index[sType]
		if exists == false {
			continue
		}

		for _, sInstance := range sInstances {
			for _, instance := range serviceDirectory.index[sType] {
				if instance == sInstance {
					invalidServices = append(
						invalidServices, fmt.Sprintf("%s.%s", sType, sInstance),
					)
				}
			}
		}
	}

	return
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
			serviceWorker, exists = serviceDirectory.services[serviceTypeName].instances[serviceInstanceName]
			if exists {
				log.Printf("Found serviceWorker: %q", serviceWorker.Name)
				return
			}
		}
	}

	err = errors.New("Unknown serviceInstance")
	return
}

func (serviceDirectory *ServiceDirectory) ServiceExists(serviceType string, serviceInstance string) (exists bool) {
	// check for serviceType
	services, exists := serviceDirectory.index[serviceType]

	if exists == false {
		return false
	}

	for _, si := range services {
		if si == serviceInstance {
			return true
		}
	}

	return false
}
