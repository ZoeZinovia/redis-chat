package entities

import (
	"redisChat/Client/services/viewmodels"
)

type Message struct {
	ID      int    `json:"id"`
	Message string `json:"message"`
	User    string `json:"user"`
}

type Messages struct {
	Messages *[]Message `json:"messages"`
}

// ToEntity converts from viewmodel to entity
func (m *Messages) ToEntity(v *viewmodels.UDPMessages) {
	tempMessages := []Message{}
	if v != nil {
		for _, msg := range v.Messages {
			tempMessages = append(tempMessages, Message{
				ID:      msg.ID,
				Message: msg.Message,
				User:    msg.User,
			})
		}
	}

	m.Messages = &tempMessages
}

// ToViewmodel converts from entity to viewmodel
func (m *Messages) ToViewmodel(v *viewmodels.UDPMessages) {

	tempMessages := []viewmodels.UDPMessage{}
	if m != nil {
		for _, msg := range *m.Messages {
			tempMessages = append(tempMessages, viewmodels.UDPMessage{
				ID:      msg.ID,
				Message: msg.Message,
				User:    msg.User,
			})
		}
	}
	v.Messages = tempMessages
}
