package services

import (
	"context"
	"errors"
)

var (
	ErrNotConnected = errors.New("you cannot send a message if you are not yet connected")
	ErrBadRequest   = errors.New("request does not follow correct syntax")
	ErrNoPermission = errors.New("action is not permitted")
)

type UDPService interface {
	// ReceiveMessage unmarshals a message received via udp and calls the handleMessage function
	ReceiveMessage(buf []byte, address string, n int) (string, error)

	// SendCloseMessages sends a close message to all clients that are still connected
	SendCloseMessages()

	// getLastMessagesForUDP gets the messages stored in the db
	GetLastMessagesForUDP(ctx context.Context) (response string, err error)
}
