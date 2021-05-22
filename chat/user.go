package chat

import (
	"encoding/json"
	"log"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type User struct {
	ID         string
	Connection *websocket.Conn
	Actions    ChatActions
}

func NewUser(conn *websocket.Conn, actions ChatActions) *User {
	return &User{
		ID:         uuid.NewString(),
		Connection: conn,
		Actions:    actions,
	}
}

func (u *User) Read() {
	for {
		_, b, err := u.Connection.ReadMessage()
		if err != nil {
			log.Printf("Error on read message: %s. ID=%s\n", err.Error(), u.ID)
			break
		} else {
			var msg Message
			err := json.Unmarshal(b, &msg)
			if err != nil {
				msg.Text = err.Error()
				_ = u.Connection.WriteMessage(websocket.TextMessage, msg.Binary())
			} else {
				u.Actions.Broadcast(&msg)
			}
		}
	}

	u.Actions.Leave(u)
}
