package repositories

import (
	"redisChat/Client/repositories/entities"
	"redisChat/Client/services/viewmodels"

	"github.com/stretchr/testify/assert"
	"testing"
)

func TestToEntity(t *testing.T) {
	// Define input viewmodel
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

	// Define desired entity
	entityMessages := &entities.Messages{
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

	tests := []struct {
		testName       string
		inputViewmodel *viewmodels.UDPMessages
		expectedEntity *entities.Messages
	}{
		{
			testName:       "case 1",
			inputViewmodel: messages,
			expectedEntity: entityMessages,
		},
		{
			testName:       "case 2",
			inputViewmodel: &viewmodels.UDPMessages{},
			expectedEntity: &entities.Messages{
				Messages: &[]entities.Message{},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {

			// Get actual result
			result := &entities.Messages{}
			result.ToEntity(test.inputViewmodel)

			// Compare
			assert.Equal(t, test.expectedEntity, result)
		})
	}
}

func TestToViewmodel(t *testing.T) {
	// Define input entity
	entityMessages := &entities.Messages{
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

	// Define desired entity
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
		testName          string
		inputEntity       *entities.Messages
		expectedViewmodel *viewmodels.UDPMessages
	}{
		{
			testName:          "case 1",
			inputEntity:       entityMessages,
			expectedViewmodel: messages,
		},
		{
			testName:    "case 2",
			inputEntity: &entities.Messages{},
			expectedViewmodel: &viewmodels.UDPMessages{
				Messages: []viewmodels.UDPMessage{},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {

			// Get actual result
			result := &viewmodels.UDPMessages{}
			test.inputEntity.ToEntity(result)

			// Compare
			assert.Equal(t, test.expectedViewmodel, result)
		})
	}
}
