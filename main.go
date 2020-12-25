package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/julienschmidt/sse"
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
	router := httprouter.New()
	timer := sse.New()

	router.ServeFiles("/js/*filepath", http.Dir("js"))
	router.ServeFiles("/css/*filepath", http.Dir("css"))

	router.GET("/", serveHomepage)

	router.Handler("GET", "/time", timer)
	go streamTime(timer)

	err := http.ListenAndServe(":81", router)
	if err != nil {
		fmt.Println("Problem starting web server: " + err.Error())
		os.Exit(-1)
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
