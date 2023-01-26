// Code generated by mockery v2.16.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// Option is an autogenerated mock type for the Option type
type Option struct {
	mock.Mock
}

// Execute provides a mock function with given fields: r
func (_m *Option) Execute(r *registry.registry) {
	_m.Called(r)
}

type mockConstructorTestingTNewOption interface {
	mock.TestingT
	Cleanup(func())
}

// NewOption creates a new instance of Option. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewOption(t mockConstructorTestingTNewOption) *Option {
	mock := &Option{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
