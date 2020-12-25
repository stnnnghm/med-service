package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/kardianos/service"
)

var (
	serviceIsRunning bool
	programIsRunning bool
	writeSync        sync.Mutex
)

const serviceName = "Simple Service"
const serviceDesc = "Simple Service Description"

type program struct{}

func (p program) Start(s service.Service) error {
	fmt.Println(s.String() + " started")

	writeSync.Lock()
	serviceIsRunning = true
	writeSync.Unlock()

	go p.run()
	return nil
}

func (p program) Stop(s service.Service) error {
	writeSync.Lock()
	serviceIsRunning = false
	writeSync.Unlock()

	for programIsRunning {
		fmt.Println(s.String() + " stopping...")
		time.Sleep(1 * time.Second)
	}

	fmt.Println(s.String() + " stopped")
	return nil
}

func (p program) run() {
	for serviceIsRunning {
		fmt.Println("Service is running")
		time.Sleep(1 * time.Second)
	}
}

func main() {
	serviceConfig := &service.Config{
		Name:        serviceName,
		DisplayName: serviceName,
		Description: serviceDesc,
	}

	prg := &program{}
	s, err := service.New(prg, serviceConfig)
	if err != nil {
		fmt.Println("Cannot create the service: " + err.Error())
	}

	err = s.Run()
	if err != nil {
		fmt.Println("Cannot start the service: " + err.Error())
	}

}
