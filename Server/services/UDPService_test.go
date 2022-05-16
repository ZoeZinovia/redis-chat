package services

import (
	"github.com/stretchr/testify/assert"
	"redisChat/Server/repositories/entities"
	"redisChat/Server/repositories/mocks"

	"context"
	"testing"
)

var ctx context.Context = context.TODO()

func TestGetLastMessagesForUDP(t *testing.T) {

	// Define mocks
	messageRepository := new(mocks.MessageRepository)

	// Define service
	udpService := NewUDPService(messageRepository)

	// Define desired response from repository layer
	repoMessages := &entities.Messages{
		Messages: []entities.Message{
			{
				ID:      3,
				Message: "Message3",
				User:    "John",
			},
			{
				ID:      6,
				Message: "Message6",
				User:    "Maria",
			},
			{
				ID:      1,
				Message: "Message1",
				User:    "John",
			},
		},
	}

	tests := []struct {
		testName              string
		expectedRepoResult    *entities.Messages
		expectedServiceResult string
		expectedServiceError  error
	}{
		{
			testName:              "case 1",
			expectedRepoResult:    repoMessages,
			expectedServiceResult: `{"messages":[{"id":6,"msg":"Message6","user":"Maria"},{"id":3,"msg":"Message3","user":"John"},{"id":1,"msg":"Message1","user":"John"}]}`,
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {

			// Define what mock repository should return
			messageRepository.On("RetrieveAll", ctx).Return(test.expectedRepoResult, test.expectedServiceError)

			// Get actual result
			result, err := udpService.GetLastMessagesForUDP(ctx)

			// Compare
			assert.Equal(t, test.expectedServiceResult, result)
			assert.Equal(t, test.expectedServiceError, err)
		})
	}

}
