// Code generated by mockery v2.12.2. DO NOT EDIT.

package mocks

import (
	entities "redisChat/Client/repositories/entities"

	mock "github.com/stretchr/testify/mock"

	testing "testing"
)

// MessageRepository is an autogenerated mock type for the MessageRepository type
type MessageRepository struct {
	mock.Mock
}

// GetMessageUserByID provides a mock function with given fields: ID
func (_m *MessageRepository) GetMessageUserByID(ID int) string {
	ret := _m.Called(ID)

	var r0 string
	if rf, ok := ret.Get(0).(func(int) string); ok {
		r0 = rf(ID)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// GetMessages provides a mock function with given fields:
func (_m *MessageRepository) GetMessages() *entities.Messages {
	ret := _m.Called()

	var r0 *entities.Messages
	if rf, ok := ret.Get(0).(func() *entities.Messages); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entities.Messages)
		}
	}

	return r0
}

// UpdateMessages provides a mock function with given fields: messages
func (_m *MessageRepository) UpdateMessages(messages *entities.Messages) {
	_m.Called(messages)
}

// NewMessageRepository creates a new instance of MessageRepository. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewMessageRepository(t testing.TB) *MessageRepository {
	mock := &MessageRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
