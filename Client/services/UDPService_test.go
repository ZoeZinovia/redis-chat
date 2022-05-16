package services

import (
	"redisChat/Client/repositories/entities"
	"redisChat/Client/repositories/mocks"
	"redisChat/Client/services/viewmodels"

	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetAllMessages(t *testing.T) {
	// Define mocks
	messageRepository := new(mocks.MessageRepository)

	// Define service
	udpService := NewUDPService(messageRepository)

	// Define desired response from repository layer
	repoMessages := &entities.Messages{
		Messages: &[]entities.Message{
			{
				ID:      1,
				Message: "Message1",
				User:    "John",
			},
			{
				ID:      2,
				Message: "Message2",
				User:    "Maria",
			},
			{
				ID:      3,
				Message: "Message3",
				User:    "John",
			},
		},
	}

	// Define desired response from service
	messages := &viewmodels.UDPMessages{
		Messages: []viewmodels.UDPMessage{
			{
				ID:      1,
				Message: "Message1",
				User:    "John",
			},
			{
				ID:      2,
				Message: "Message2",
				User:    "Maria",
			},
			{
				ID:      3,
				Message: "Message3",
				User:    "John",
			},
		},
	}

	tests := []struct {
		testName              string
		expectedRepoResult    *entities.Messages
		expectedServiceResult *viewmodels.UDPMessages
	}{
		{
			testName:              "case 1",
			expectedRepoResult:    repoMessages,
			expectedServiceResult: messages,
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {

			// Define what mock repository should return
			messageRepository.On("GetMessages").Return(test.expectedRepoResult)

			// Get actual result
			result := udpService.GetAllMessages()

			// Compare
			assert.Equal(t, test.expectedServiceResult, result)

		})
	}

}
