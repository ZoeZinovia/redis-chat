package handlers

import (
	"net/http"
	"net/http/httptest"
	repoMock "redisChat/Client/repositories/mocks"
	servMock "redisChat/Client/services/mocks"
	"redisChat/Client/services/viewmodels"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestGetMessages tests the handler function that gets all messages
func TestGetMessages(t *testing.T) {
	udpService := new(servMock.UDPService)
	messageRepository := new(repoMock.MessageRepository)
	messageHandler := MessageHandler{
		udpService:        udpService,
		messageRepository: messageRepository,
	}

	expectedMessage := &viewmodels.UDPMessages{
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
		testName           string
		expectedResult     *viewmodels.UDPMessages
		expectedError      error
		expectedStatusCode int
	}{
		{
			testName:           "case 1",
			expectedResult:     expectedMessage,
			expectedStatusCode: http.StatusOK,
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {

			// Create fake request
			req, err := http.NewRequest("GET", "/messages", nil)
			assert.NoError(t, err)

			// Define what mock service should return
			udpService.On("GetAllMessages").Return(test.expectedResult)

			// Create response recorder and set handler
			res := httptest.NewRecorder()
			handler := http.HandlerFunc(messageHandler.GetMessages)

			// Call handler and pass req and res
			handler.ServeHTTP(res, req)

			// Check that status code is as expected
			assert.Equal(t, test.expectedStatusCode, res.Code)

		})
	}

}
