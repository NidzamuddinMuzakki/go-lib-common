// Code generated by mockery v2.12.2. DO NOT EDIT.

package mocks

import (
	testing "testing"

	mock "github.com/stretchr/testify/mock"

	time "time"
)

// TimeItf is an autogenerated mock type for the TimeItf type
type TimeItf struct {
	mock.Mock
}

// Now provides a mock function with given fields:
func (_m *TimeItf) Now() time.Time {
	ret := _m.Called()

	var r0 time.Time
	if rf, ok := ret.Get(0).(func() time.Time); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(time.Time)
	}

	return r0
}

// ToDateTime provides a mock function with given fields:
func (_m *TimeItf) ToDateTime() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// NewTimeItf creates a new instance of TimeItf. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewTimeItf(t testing.TB) *TimeItf {
	mock := &TimeItf{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
