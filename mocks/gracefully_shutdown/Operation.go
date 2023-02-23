// Code generated by mockery v2.20.2. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// Operation is an autogenerated mock type for the Operation type
type Operation struct {
	mock.Mock
}

// Execute provides a mock function with given fields: ctx
func (_m *Operation) Execute(ctx context.Context) error {
	ret := _m.Called(ctx)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewOperation interface {
	mock.TestingT
	Cleanup(func())
}

// NewOperation creates a new instance of Operation. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewOperation(t mockConstructorTestingTNewOperation) *Operation {
	mock := &Operation{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
