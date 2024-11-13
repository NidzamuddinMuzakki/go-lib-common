// Code generated by mockery v2.20.2. DO NOT EDIT.

package mocks

import (
	context "context"

	exporter "bitbucket.org/moladinTech/go-lib-common/exporter"
	mock "github.com/stretchr/testify/mock"
)

// Exporter is an autogenerated mock type for the Exporter type
type Exporter struct {
	mock.Mock
}

// Export provides a mock function with given fields: _a0, _a1
func (_m *Exporter) Export(_a0 context.Context, _a1 interface{}) (exporter.ResultExport, error) {
	ret := _m.Called(_a0, _a1)

	var r0 exporter.ResultExport
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, interface{}) (exporter.ResultExport, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, interface{}) exporter.ResultExport); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(exporter.ResultExport)
	}

	if rf, ok := ret.Get(1).(func(context.Context, interface{}) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewExporter interface {
	mock.TestingT
	Cleanup(func())
}

// NewExporter creates a new instance of Exporter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewExporter(t mockConstructorTestingTNewExporter) *Exporter {
	mock := &Exporter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
