// Code generated by mockery v2.16.0. DO NOT EDIT.

package mocks

import (
	gin "github.com/gin-gonic/gin"
	mock "github.com/stretchr/testify/mock"
)

// IMiddlewareAuth is an autogenerated mock type for the IMiddlewareAuth type
type IMiddlewareAuth struct {
	mock.Mock
}

// Auth provides a mock function with given fields:
func (_m *IMiddlewareAuth) Auth() gin.HandlerFunc {
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

// AuthSignature provides a mock function with given fields:
func (_m *IMiddlewareAuth) AuthSignature() gin.HandlerFunc {
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

// AuthToken provides a mock function with given fields:
func (_m *IMiddlewareAuth) AuthToken() gin.HandlerFunc {
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

// AuthXApiKey provides a mock function with given fields:
func (_m *IMiddlewareAuth) AuthXApiKey() gin.HandlerFunc {
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

type mockConstructorTestingTNewIMiddlewareAuth interface {
	mock.TestingT
	Cleanup(func())
}

// NewIMiddlewareAuth creates a new instance of IMiddlewareAuth. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewIMiddlewareAuth(t mockConstructorTestingTNewIMiddlewareAuth) *IMiddlewareAuth {
	mock := &IMiddlewareAuth{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}