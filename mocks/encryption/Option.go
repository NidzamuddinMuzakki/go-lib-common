// Code generated by mockery v2.16.0. DO NOT EDIT.

package mocks

import (
	encryption "bitbucket.org/moladinTech/go-lib-common/encryption"
	mock "github.com/stretchr/testify/mock"
)

// Option is an autogenerated mock type for the Option type
type Option struct {
	mock.Mock
}

// Execute provides a mock function with given fields: _a0
func (_m *Option) Execute(_a0 *encryption.EncryptionPackage) {
	_m.Called(_a0)
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