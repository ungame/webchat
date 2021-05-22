package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"webchat/chat"

	"github.com/gorilla/websocket"
)

func main() {
	wsURL := url.URL{
		Scheme: "ws",
		Host:   "localhost:8081",
		Path:   "/chat",
		
	}

	dialer := websocket.DefaultDialer

	conn, _, err := dialer.Dial(wsURL.String(), nil)
	if err != nil {
		log.Panicln(err)
	}
	defer conn.Close()

	for {
		msgType, payload, err := conn.ReadMessage()
		if err != nil {
			log.Panicln(err)
		}
		fmt.Println("Received Message Type: ", msgType)

		var message chat.Message

		err = json.Unmarshal(payload, &message)
		if err != nil {
			log.Println("Error on decode message:", err.Error())
		} else {
			fmt.Println("Message: ", message.Text)
		}
	}

}
