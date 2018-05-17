package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func wsRoute(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		pollData := processingQueue.Poll()
		if pollData == nil {
			continue
		}
		pollDataMember := pollData.(queueMember)
		if pollDataMember.Message == "stop" {
			c.WriteJSON(pollDataMember)
			c.WriteMessage(websocket.CloseMessage, []byte{})
			break
		}
		err = c.WriteJSON(pollDataMember)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}
