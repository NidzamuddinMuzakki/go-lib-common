// Code generated by mockery v2.16.0. DO NOT EDIT.

package mocks

import (
	context "context"

	kafka "bitbucket.org/moladinTech/go-lib-common/kafka"
	mock "github.com/stretchr/testify/mock"
)

// IPublisher is an autogenerated mock type for the IPublisher type
type IPublisher struct {
	mock.Mock
}

// Publish provides a mock function with given fields: ctx, topic, message
func (_m *IPublisher) Publish(ctx context.Context, topic kafka.Topic, message kafka.IMessage) (int32, int64, error) {
	ret := _m.Called(ctx, topic, message)

	var r0 int32
	if rf, ok := ret.Get(0).(func(context.Context, kafka.Topic, kafka.IMessage) int32); ok {
		r0 = rf(ctx, topic, message)
	} else {
		r0 = ret.Get(0).(int32)
	}

	var r1 int64
	if rf, ok := ret.Get(1).(func(context.Context, kafka.Topic, kafka.IMessage) int64); ok {
		r1 = rf(ctx, topic, message)
	} else {
		r1 = ret.Get(1).(int64)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, kafka.Topic, kafka.IMessage) error); ok {
		r2 = rf(ctx, topic, message)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

type mockConstructorTestingTNewIPublisher interface {
	mock.TestingT
	Cleanup(func())
}

// NewIPublisher creates a new instance of IPublisher. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewIPublisher(t mockConstructorTestingTNewIPublisher) *IPublisher {
	mock := &IPublisher{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
