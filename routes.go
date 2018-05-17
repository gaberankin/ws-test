package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

type wsResponse struct {
	Type    string      `json:"Type"`
	Message interface{} `json:"Message"`
}

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
		mes, err := json.Marshal(pollDataMember)
		if err != nil {
			log.Printf("Unable to properly format poll data %v\n", pollDataMember)
			break
		}
		if pollDataMember.Message == "stop" {
			c.WriteMessage(websocket.CloseMessage, mes)
			break
		}
		err = c.WriteMessage(websocket.TextMessage, mes)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}
