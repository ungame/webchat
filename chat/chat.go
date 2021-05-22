package chat

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Chat interface {
	Handler(w http.ResponseWriter, r *http.Request)
	Start()
	Close()
}

type ChatActions interface {
	Broadcast(message *Message)
	Leave(user *User)
}

type chatImpl struct {
	upgrader    *websocket.Upgrader
	locker      *sync.Mutex
	connections map[string]*User
	join        chan *User
	leave       chan *User
	broadcast   chan *Message
	close       chan bool
}

func New(upgrader *websocket.Upgrader) Chat {
	return &chatImpl{
		upgrader:    upgrader,
		locker:      &sync.Mutex{},
		connections: make(map[string]*User),
		join:        make(chan *User),
		leave:       make(chan *User),
		broadcast:   make(chan *Message),
		close:       make(chan bool),
	}
}

func (c *chatImpl) Handler(w http.ResponseWriter, r *http.Request) {

	conn, err := c.upgrader.Upgrade(w, r, nil)
	if err != nil {
		w.Write([]byte("upgrade ws failed: " + err.Error()))
		return
	}

	c.join <- NewUser(conn, c)
}

func (c *chatImpl) Start() {
	for {
		select {
		case user := <-c.join:
			c.addConnection(user)
			c.dispatch(&Message{
				Text: fmt.Sprintf("%s joined", user.ID),
			})

		case user := <-c.leave:
			c.closeConnection(user)
			c.dispatch(&Message{
				Text: fmt.Sprintf("%s left", user.ID),
			})

		case message := <-c.broadcast:
			c.dispatch(message)
			
		case _ = <-c.close:
			c.Close()
		}
	}
}

func (c *chatImpl) Close() {
	c.locker.Lock()
	for _, conn := range c.connections {
		c.Leave(conn)
	}
	defer c.locker.Unlock()
}

func (c *chatImpl) Leave(user *User) {
	c.leave <- user
}

func (c *chatImpl) Broadcast(message *Message) {
	c.broadcast <- message
}

func (c *chatImpl) addConnection(user *User) {
	c.connections[user.ID] = user
	go user.Read()
}

func (c *chatImpl) closeConnection(user *User) {
	conn, exists := c.connections[user.ID]
	if exists {
		err := conn.Connection.Close()
		if err != nil {
			log.Println("Error on close Ws Connection: ", err.Error())
		}
		delete(c.connections, user.ID)
		log.Println("User left: ", user.ID)
	}
}

func (c *chatImpl) dispatch(message *Message) {
	log.Println("Broadcast: ", message.Text)
	for _, conn := range c.connections {
		err := conn.Connection.WriteJSON(message)
		if err != nil {
			log.Println("Error on write json: ", err.Error())
			c.Leave(conn)
		}
	}
}
