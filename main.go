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
	http.HandleFunc("/wshandler", wsRoute)
	http.HandleFunc("/start", func(w http.ResponseWriter, r *http.Request) {
		go queueProcessor()
		w.Write([]byte("starting"))
	})
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

func startRunning() {
	running = true
}
func stopRunning() {
	running = false
}
func isRunning() bool {
	return running
}
