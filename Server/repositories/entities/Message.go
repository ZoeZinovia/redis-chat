package entities

import "redisChat/Server/services/viewmodels"

type Message struct {
	ID      int    `json:"id"`
	Message string `json:"msg"`
	User    string `json:"user"`
}

type Messages struct {
	Messages []Message `json:"messages"`
}

func ToEntity(v *viewmodels.UDPMessage) *Message {
	newEntity := Message{
		ID:      v.ID,
		Message: v.Message,
		User:    v.User,
	}
	return &newEntity
}

func ToViewmodel(m *Message) *viewmodels.UDPMessage {
	newViewmodel := viewmodels.UDPMessage{
		ID:      m.ID,
		Message: m.Message,
		User:    m.User,
	}

	return &newViewmodel
}
