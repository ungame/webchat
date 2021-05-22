package chat

import (
	"encoding/json"
)

type Message struct {
	Text string `json:"text"`
}

func (m *Message) Binary() []byte {
	b, _ := json.Marshal(m)
	return b
}
