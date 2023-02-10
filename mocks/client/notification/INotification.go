// Code generated by mockery v2.18.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// INotification is an autogenerated mock type for the INotification type
type INotification struct {
	mock.Mock
}

// GetFormattedMessage provides a mock function with given fields: logCtx, ctx, message
func (_m *INotification) GetFormattedMessage(logCtx string, ctx context.Context, message interface{}) string {
	ret := _m.Called(logCtx, ctx, message)

	var r0 string
	if rf, ok := ret.Get(0).(func(string, context.Context, interface{}) string); ok {
		r0 = rf(logCtx, ctx, message)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Health provides a mock function with given fields: ctx
func (_m *INotification) Health(ctx context.Context) error {
	ret := _m.Called(ctx)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Send provides a mock function with given fields: ctx, message
func (_m *INotification) Send(ctx context.Context, message string) error {
	ret := _m.Called(ctx, message)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, message)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewINotification interface {
	mock.TestingT
	Cleanup(func())
}

// NewINotification creates a new instance of INotification. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewINotification(t mockConstructorTestingTNewINotification) *INotification {
	mock := &INotification{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
