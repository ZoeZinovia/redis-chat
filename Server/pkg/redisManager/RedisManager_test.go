package redisManager

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redismock"
	"github.com/stretchr/testify/assert"
	"redisChat/Server/repositories/entities"
	"testing"
)

var ctx = context.TODO()

func TestSave(t *testing.T) {
	db, mock := redismock.NewClientMock()

	redis := RedisClient{client: db}

	tests := []struct {
		testName   string
		inputKey   int
		inputValue *entities.Message
	}{
		{
			testName: "case 1",
			inputKey: 1,
			inputValue: &entities.Message{
				ID:      1,
				Message: "Message1",
				User:    "John",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			// Define what mock redis should expect
			val, err := json.Marshal(test.inputValue)
			if err != nil {
				return
			}

			mock.ExpectSet(fmt.Sprintf("%d", test.inputKey), val, 0).SetVal("")
			mock.ExpectScan(0, "*", 0).SetVal([]string{}, 0)
			mock.ExpectDel(fmt.Sprintf("%d", test.inputKey)).SetVal(0)

			// Get actual result
			err = redis.Save(ctx, test.inputKey, test.inputValue)

			// Compare
			assert.NoError(t, err)
		})
	}
}

func TestRetrieveByKey(t *testing.T) {
	db, mock := redismock.NewClientMock()

	redis := RedisClient{client: db}

	tests := []struct {
		testName           string
		inputKey           int
		expectedValue      *entities.Message
		expectedRepoResult string
	}{
		{
			testName: "case 1",
			inputKey: 1,
			expectedValue: &entities.Message{
				ID:      1,
				Message: "Message1",
				User:    "John",
			},
			expectedRepoResult: `{"id":1,"msg":"Message1","user":"John"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			// Define what mock redis should expect
			mock.ExpectGet(fmt.Sprintf("%d", test.inputKey)).SetVal(test.expectedRepoResult)

			// Get actual result
			result := &entities.Message{}
			err := redis.RetrieveByKey(ctx, test.inputKey, result)

			// Compare
			assert.NoError(t, err)
			assert.Equal(t, test.expectedValue, result)
		})
	}
}
