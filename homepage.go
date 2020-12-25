package main

import (
	"html/template"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

// HomePage holds data about the homepage for the Go templates
type HomePage struct {
	Time string
}

func serveHomepage(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	writeSync.Lock()
	programIsRunning = true
	writeSync.Unlock()

	var homepage HomePage
	homepage.Time = time.Now().String()

	tmpl := template.Must(template.ParseFiles("html/homepage.html"))
	_ = tmpl.Execute(w, homepage)

	writeSync.Lock()
	programIsRunning = false
	writeSync.Unlock()
}
