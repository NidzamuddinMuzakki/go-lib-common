// Code generated by mockery v2.20.2. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// IEncryption is an autogenerated mock type for the IEncryption type
type IEncryption struct {
	mock.Mock
}

// Decrypt provides a mock function with given fields: data, salt
func (_m *IEncryption) Decrypt(data string, salt []byte) ([]byte, error) {
	ret := _m.Called(data, salt)

	var r0 []byte
	var r1 error
	if rf, ok := ret.Get(0).(func(string, []byte) ([]byte, error)); ok {
		return rf(data, salt)
	}
	if rf, ok := ret.Get(0).(func(string, []byte) []byte); ok {
		r0 = rf(data, salt)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	if rf, ok := ret.Get(1).(func(string, []byte) error); ok {
		r1 = rf(data, salt)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Encrypt provides a mock function with given fields: data, salt
func (_m *IEncryption) Encrypt(data string, salt []byte) ([]byte, error) {
	ret := _m.Called(data, salt)

	var r0 []byte
	var r1 error
	if rf, ok := ret.Get(0).(func(string, []byte) ([]byte, error)); ok {
		return rf(data, salt)
	}
	if rf, ok := ret.Get(0).(func(string, []byte) []byte); ok {
		r0 = rf(data, salt)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	if rf, ok := ret.Get(1).(func(string, []byte) error); ok {
		r1 = rf(data, salt)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GenerateSalt provides a mock function with given fields: key
func (_m *IEncryption) GenerateSalt(key string) []byte {
	ret := _m.Called(key)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(string) []byte); ok {
		r0 = rf(key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	return r0
}

type mockConstructorTestingTNewIEncryption interface {
	mock.TestingT
	Cleanup(func())
}

// NewIEncryption creates a new instance of IEncryption. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewIEncryption(t mockConstructorTestingTNewIEncryption) *IEncryption {
	mock := &IEncryption{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
