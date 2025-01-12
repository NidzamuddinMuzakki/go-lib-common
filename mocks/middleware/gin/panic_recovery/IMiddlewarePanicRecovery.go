// Code generated by mockery v2.20.2. DO NOT EDIT.

package mocks

import (
	gin "github.com/gin-gonic/gin"
	mock "github.com/stretchr/testify/mock"
)

// IMiddlewarePanicRecovery is an autogenerated mock type for the IMiddlewarePanicRecovery type
type IMiddlewarePanicRecovery struct {
	mock.Mock
}

// PanicRecoveryMiddleware provides a mock function with given fields:
func (_m *IMiddlewarePanicRecovery) PanicRecoveryMiddleware() gin.HandlerFunc {
	ret := _m.Called()

	var r0 gin.HandlerFunc
	if rf, ok := ret.Get(0).(func() gin.HandlerFunc); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(gin.HandlerFunc)
		}
	}

	return r0
}

type mockConstructorTestingTNewIMiddlewarePanicRecovery interface {
	mock.TestingT
	Cleanup(func())
}

// NewIMiddlewarePanicRecovery creates a new instance of IMiddlewarePanicRecovery. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewIMiddlewarePanicRecovery(t mockConstructorTestingTNewIMiddlewarePanicRecovery) *IMiddlewarePanicRecovery {
	mock := &IMiddlewarePanicRecovery{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
