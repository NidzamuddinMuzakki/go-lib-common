// Code generated by mockery v2.16.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// GenerateAndVerify is an autogenerated mock type for the GenerateAndVerify type
type GenerateAndVerify struct {
	mock.Mock
}

// Generate provides a mock function with given fields: ctx, key
func (_m *GenerateAndVerify) Generate(ctx context.Context, key string) (string, error) {
	ret := _m.Called(ctx, key)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, string) string); ok {
		r0 = rf(ctx, key)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Verify provides a mock function with given fields: ctx, key, sign
func (_m *GenerateAndVerify) Verify(ctx context.Context, key string, sign string) bool {
	ret := _m.Called(ctx, key, sign)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, string, string) bool); ok {
		r0 = rf(ctx, key, sign)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

type mockConstructorTestingTNewGenerateAndVerify interface {
	mock.TestingT
	Cleanup(func())
}

// NewGenerateAndVerify creates a new instance of GenerateAndVerify. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewGenerateAndVerify(t mockConstructorTestingTNewGenerateAndVerify) *GenerateAndVerify {
	mock := &GenerateAndVerify{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
