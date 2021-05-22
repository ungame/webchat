package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"webchat/chat"

	"github.com/gorilla/websocket"
)

var templates *template.Template

func main() {

	ctrlC := make(chan os.Signal)
	signal.Notify(ctrlC, os.Interrupt)

	templates = template.Must(template.ParseGlob("./views/*.html"))

	upgrader := &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	c := chat.New(upgrader)

	http.HandleFunc("/", index)
	http.HandleFunc("/chat", c.Handler)

	go c.Start()

	go func() {
		for {
			select {
			case <-ctrlC:
				fmt.Println("Closing Webchat...")
				c.Close()
				os.Exit(0)
			}
		}
	}()

	log.Println("Running on :8081...")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func index(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "index.html", "ws://localhost:8081/chat")
}
