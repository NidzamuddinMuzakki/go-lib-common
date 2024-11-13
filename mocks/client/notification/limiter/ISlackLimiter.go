// Code generated by mockery v2.20.2. DO NOT EDIT.

package mocks

import (
	context "context"

	cache "bitbucket.org/moladinTech/go-lib-common/cache"

	mock "github.com/stretchr/testify/mock"
)

// ISlackLimiter is an autogenerated mock type for the ISlackLimiter type
type ISlackLimiter struct {
	mock.Mock
}

// LimitChecker provides a mock function with given fields: ctx, data
func (_m *ISlackLimiter) LimitChecker(ctx context.Context, data cache.Data) (bool, error) {
	ret := _m.Called(ctx, data)

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, cache.Data) (bool, error)); ok {
		return rf(ctx, data)
	}
	if rf, ok := ret.Get(0).(func(context.Context, cache.Data) bool); ok {
		r0 = rf(ctx, data)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(context.Context, cache.Data) error); ok {
		r1 = rf(ctx, data)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewISlackLimiter interface {
	mock.TestingT
	Cleanup(func())
}

// NewISlackLimiter creates a new instance of ISlackLimiter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewISlackLimiter(t mockConstructorTestingTNewISlackLimiter) *ISlackLimiter {
	mock := &ISlackLimiter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
