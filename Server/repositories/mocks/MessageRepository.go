// Code generated by mockery v2.12.2. DO NOT EDIT.

package mocks

import (
	context "context"
	entities "redisChat/Server/repositories/entities"

	mock "github.com/stretchr/testify/mock"

	testing "testing"
)

// MessageRepository is an autogenerated mock type for the MessageRepository type
type MessageRepository struct {
	mock.Mock
}

// Delete provides a mock function with given fields: ctx, messageID
func (_m *MessageRepository) Delete(ctx context.Context, messageID int) error {
	ret := _m.Called(ctx, messageID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int) error); ok {
		r0 = rf(ctx, messageID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Flush provides a mock function with given fields: ctx
func (_m *MessageRepository) Flush(ctx context.Context) {
	_m.Called(ctx)
}

// RetrieveAll provides a mock function with given fields: ctx
func (_m *MessageRepository) RetrieveAll(ctx context.Context) (*entities.Messages, error) {
	ret := _m.Called(ctx)

	var r0 *entities.Messages
	if rf, ok := ret.Get(0).(func(context.Context) *entities.Messages); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entities.Messages)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RetrieveByMessageID provides a mock function with given fields: ctx, messageID
func (_m *MessageRepository) RetrieveByMessageID(ctx context.Context, messageID int) (*entities.Message, error) {
	ret := _m.Called(ctx, messageID)

	var r0 *entities.Message
	if rf, ok := ret.Get(0).(func(context.Context, int) *entities.Message); ok {
		r0 = rf(ctx, messageID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entities.Message)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, messageID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Save provides a mock function with given fields: ctx, message
func (_m *MessageRepository) Save(ctx context.Context, message *entities.Message) error {
	ret := _m.Called(ctx, message)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *entities.Message) error); ok {
		r0 = rf(ctx, message)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewMessageRepository creates a new instance of MessageRepository. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewMessageRepository(t testing.TB) *MessageRepository {
	mock := &MessageRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
