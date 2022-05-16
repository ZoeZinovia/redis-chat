package repositories

import (
	"github.com/stretchr/testify/assert"
	"redisChat/Server/repositories/entities"
	"redisChat/Server/services/viewmodels"
	"testing"
)

func TestToEntity(t *testing.T) {
	// Define input viewmodel
	message := &viewmodels.UDPMessage{
		ID:      1,
		Message: "Message1",
		User:    "John",
	}

	// Define desired entity
	entityMessage := &entities.Message{
		ID:      1,
		Message: "Message1",
		User:    "John",
	}

	tests := []struct {
		testName       string
		inputViewmodel *viewmodels.UDPMessage
		expectedEntity *entities.Message
	}{
		{
			testName:       "case 1",
			inputViewmodel: message,
			expectedEntity: entityMessage,
		},
		{
			testName:       "case 2",
			inputViewmodel: &viewmodels.UDPMessage{},
			expectedEntity: &entities.Message{},
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {

			// Get actual result
			result := entities.ToEntity(test.inputViewmodel)

			// Compare
			assert.Equal(t, test.expectedEntity, result)
		})
	}
}

func TestToViewmodel(t *testing.T) {
	// Define input entity
	entityMessage := &entities.Message{
		ID:      1,
		Message: "Message1",
		User:    "John",
	}

	// Define desired entity
	message := &viewmodels.UDPMessage{
		ID:      1,
		Message: "Message1",
		User:    "John",
	}

	tests := []struct {
		testName          string
		inputEntity       *entities.Message
		expectedViewmodel *viewmodels.UDPMessage
	}{
		{
			testName:          "case 1",
			inputEntity:       entityMessage,
			expectedViewmodel: message,
		},
		{
			testName:          "case 2",
			inputEntity:       &entities.Message{},
			expectedViewmodel: &viewmodels.UDPMessage{},
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {

			// Get actual result
			result := entities.ToViewmodel(test.inputEntity)

			// Compare
			assert.Equal(t, test.expectedViewmodel, result)
		})
	}
}
