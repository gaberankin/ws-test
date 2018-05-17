package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/hishboy/gocommons/lang"
)

var processingQueue = lang.NewQueue()
var running = false

type queueMember struct {
	Message string
}

func main() {
	addr := "localhost:3000"
	flag.Parse()
	log.SetFlags(0)
	// where the websocket connection is handled
	http.HandleFunc("/wshandler", wsRoute)
	// Simple endpoint that lets me start the queue at my leisure
	http.HandleFunc("/start", func(w http.ResponseWriter, r *http.Request) {
		go queueProcessor()
		w.Write([]byte("starting"))
	})
	// maybe one day i'll put this into a template or something stupid?  idk
	http.Handle("/", http.FileServer(http.Dir("./static")))
	log.Printf("Starting server on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func queueProcessor() {
	if isRunning() {
		return
	}
	startRunning()
	defer stopRunning()
	defer processingQueue.Push(queueMember{Message: "stop"})
	for i := 0; i < 20; i++ {
		time.Sleep(time.Second * 5)
		processingQueue.Push(queueMember{
			Message: fmt.Sprintf("%d", i),
		})
	}
}

// helper functions to prevent multiple queues.  could probably be enclosed better (container struct for queue, with processing
// function as a method of that struct), but good enough for the experiment.
func startRunning() {
	running = true
}
func stopRunning() {
	running = false
}
func isRunning() bool {
	return running
}
