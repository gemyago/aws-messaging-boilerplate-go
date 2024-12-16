// Code generated by mockery. DO NOT EDIT.

//go:build !release

package v1controllers

import (
	context "context"

	handlers "github.com/gemyago/aws-sqs-boilerplate-go/internal/api/http/v1routes/handlers"
	mock "github.com/stretchr/testify/mock"

	models "github.com/gemyago/aws-sqs-boilerplate-go/internal/api/http/v1routes/models"
)

// mockMessagesCommands is an autogenerated mock type for the messagesCommands type
type mockMessagesCommands struct {
	mock.Mock
}

type mockMessagesCommands_Expecter struct {
	mock *mock.Mock
}

func (_m *mockMessagesCommands) EXPECT() *mockMessagesCommands_Expecter {
	return &mockMessagesCommands_Expecter{mock: &_m.Mock}
}

// ProcessMessage provides a mock function with given fields: _a0, _a1
func (_m *mockMessagesCommands) ProcessMessage(_a0 context.Context, _a1 *models.Message) error {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for ProcessMessage")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.Message) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockMessagesCommands_ProcessMessage_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ProcessMessage'
type mockMessagesCommands_ProcessMessage_Call struct {
	*mock.Call
}

// ProcessMessage is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 *models.Message
func (_e *mockMessagesCommands_Expecter) ProcessMessage(_a0 interface{}, _a1 interface{}) *mockMessagesCommands_ProcessMessage_Call {
	return &mockMessagesCommands_ProcessMessage_Call{Call: _e.mock.On("ProcessMessage", _a0, _a1)}
}

func (_c *mockMessagesCommands_ProcessMessage_Call) Run(run func(_a0 context.Context, _a1 *models.Message)) *mockMessagesCommands_ProcessMessage_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*models.Message))
	})
	return _c
}

func (_c *mockMessagesCommands_ProcessMessage_Call) Return(_a0 error) *mockMessagesCommands_ProcessMessage_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockMessagesCommands_ProcessMessage_Call) RunAndReturn(run func(context.Context, *models.Message) error) *mockMessagesCommands_ProcessMessage_Call {
	_c.Call.Return(run)
	return _c
}

// PublishMessage provides a mock function with given fields: _a0, _a1
func (_m *mockMessagesCommands) PublishMessage(_a0 context.Context, _a1 *handlers.MessagesPublishMessageRequest) error {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for PublishMessage")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *handlers.MessagesPublishMessageRequest) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockMessagesCommands_PublishMessage_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'PublishMessage'
type mockMessagesCommands_PublishMessage_Call struct {
	*mock.Call
}

// PublishMessage is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 *handlers.MessagesPublishMessageRequest
func (_e *mockMessagesCommands_Expecter) PublishMessage(_a0 interface{}, _a1 interface{}) *mockMessagesCommands_PublishMessage_Call {
	return &mockMessagesCommands_PublishMessage_Call{Call: _e.mock.On("PublishMessage", _a0, _a1)}
}

func (_c *mockMessagesCommands_PublishMessage_Call) Run(run func(_a0 context.Context, _a1 *handlers.MessagesPublishMessageRequest)) *mockMessagesCommands_PublishMessage_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*handlers.MessagesPublishMessageRequest))
	})
	return _c
}

func (_c *mockMessagesCommands_PublishMessage_Call) Return(_a0 error) *mockMessagesCommands_PublishMessage_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockMessagesCommands_PublishMessage_Call) RunAndReturn(run func(context.Context, *handlers.MessagesPublishMessageRequest) error) *mockMessagesCommands_PublishMessage_Call {
	_c.Call.Return(run)
	return _c
}

// newMockMessagesCommands creates a new instance of mockMessagesCommands. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func newMockMessagesCommands(t interface {
	mock.TestingT
	Cleanup(func())
}) *mockMessagesCommands {
	mock := &mockMessagesCommands{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
