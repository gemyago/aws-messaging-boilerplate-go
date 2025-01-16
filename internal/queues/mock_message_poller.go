// Code generated by mockery. DO NOT EDIT.

//go:build !release

package queues

import (
	context "context"

	awsapi "github.com/gemyago/aws-sqs-boilerplate-go/internal/services/awsapi"

	mock "github.com/stretchr/testify/mock"
)

// mockMessagePoller is an autogenerated mock type for the messagePoller type
type mockMessagePoller struct {
	mock.Mock
}

type mockMessagePoller_Expecter struct {
	mock *mock.Mock
}

func (_m *mockMessagePoller) EXPECT() *mockMessagePoller_Expecter {
	return &mockMessagePoller_Expecter{mock: &_m.Mock}
}

// RegisterQueue provides a mock function with given fields: queue
func (_m *mockMessagePoller) RegisterQueue(queue awsapi.MessagesPollerQueue) {
	_m.Called(queue)
}

// mockMessagePoller_RegisterQueue_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RegisterQueue'
type mockMessagePoller_RegisterQueue_Call struct {
	*mock.Call
}

// RegisterQueue is a helper method to define mock.On call
//   - queue awsapi.MessagesPollerQueue
func (_e *mockMessagePoller_Expecter) RegisterQueue(queue interface{}) *mockMessagePoller_RegisterQueue_Call {
	return &mockMessagePoller_RegisterQueue_Call{Call: _e.mock.On("RegisterQueue", queue)}
}

func (_c *mockMessagePoller_RegisterQueue_Call) Run(run func(queue awsapi.MessagesPollerQueue)) *mockMessagePoller_RegisterQueue_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(awsapi.MessagesPollerQueue))
	})
	return _c
}

func (_c *mockMessagePoller_RegisterQueue_Call) Return() *mockMessagePoller_RegisterQueue_Call {
	_c.Call.Return()
	return _c
}

func (_c *mockMessagePoller_RegisterQueue_Call) RunAndReturn(run func(awsapi.MessagesPollerQueue)) *mockMessagePoller_RegisterQueue_Call {
	_c.Run(run)
	return _c
}

// Start provides a mock function with given fields: ctx
func (_m *mockMessagePoller) Start(ctx context.Context) error {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for Start")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockMessagePoller_Start_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Start'
type mockMessagePoller_Start_Call struct {
	*mock.Call
}

// Start is a helper method to define mock.On call
//   - ctx context.Context
func (_e *mockMessagePoller_Expecter) Start(ctx interface{}) *mockMessagePoller_Start_Call {
	return &mockMessagePoller_Start_Call{Call: _e.mock.On("Start", ctx)}
}

func (_c *mockMessagePoller_Start_Call) Run(run func(ctx context.Context)) *mockMessagePoller_Start_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *mockMessagePoller_Start_Call) Return(_a0 error) *mockMessagePoller_Start_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockMessagePoller_Start_Call) RunAndReturn(run func(context.Context) error) *mockMessagePoller_Start_Call {
	_c.Call.Return(run)
	return _c
}

// newMockMessagePoller creates a new instance of mockMessagePoller. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func newMockMessagePoller(t interface {
	mock.TestingT
	Cleanup(func())
}) *mockMessagePoller {
	mock := &mockMessagePoller{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
