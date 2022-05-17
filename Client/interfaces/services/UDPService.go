package services

import (
	"errors"
	"redisChat/Client/services/viewmodels"
)

var ErrPermissionDenied = errors.New("permission denied")

type UDPService interface {
	// ReceiveMessage unmarshals a message received via udp and calls the handleMessage function
	ReceiveMessage(buf []byte, n int)

	// SendMessages marshals the provided message and sends it via udp
	SendMessage(message *viewmodels.UDPMessage) (err error)

	// DeleteMessage marshals the delete message and sends it via udp
	DeleteMessage(messageID int, message *viewmodels.UDPMessage) (err error)

	// GetAllMessages retrieves all messages that are stored in memory
	// Set on connect and updated with every broadcast
	GetAllMessages() (messages *viewmodels.UDPMessages)
}
