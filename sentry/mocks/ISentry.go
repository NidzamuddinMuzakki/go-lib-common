// Code generated by mockery v2.20.2. DO NOT EDIT.

package mocks

import (
	context "context"

	go_lib_commonsentry "bitbucket.org/moladinTech/go-lib-common/sentry"
	gin "github.com/gin-gonic/gin"

	http "net/http"

	mock "github.com/stretchr/testify/mock"

	sentry "github.com/getsentry/sentry-go"

	time "time"
)

// ISentry is an autogenerated mock type for the ISentry type
type ISentry struct {
	mock.Mock
}

// CaptureException provides a mock function with given fields: exception
func (_m *ISentry) CaptureException(exception error) *sentry.EventID {
	ret := _m.Called(exception)

	var r0 *sentry.EventID
	if rf, ok := ret.Get(0).(func(error) *sentry.EventID); ok {
		r0 = rf(exception)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*sentry.EventID)
		}
	}

	return r0
}

// Finish provides a mock function with given fields: span
func (_m *ISentry) Finish(span *sentry.Span) {
	_m.Called(span)
}

// Flush provides a mock function with given fields: timeout
func (_m *ISentry) Flush(timeout time.Duration) bool {
	ret := _m.Called(timeout)

	var r0 bool
	if rf, ok := ret.Get(0).(func(time.Duration) bool); ok {
		r0 = rf(timeout)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// GetGinMiddleware provides a mock function with given fields:
func (_m *ISentry) GetGinMiddleware() gin.HandlerFunc {
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

// HandlingPanic provides a mock function with given fields: err
func (_m *ISentry) HandlingPanic(err interface{}) {
	_m.Called(err)
}

// SetEventCapture provides a mock function with given fields: eventName, data
func (_m *ISentry) SetEventCapture(eventName string, data interface{}) {
	_m.Called(eventName, data)
}

// SetIntegrationCapture provides a mock function with given fields: eventName, request, response
func (_m *ISentry) SetIntegrationCapture(eventName string, request interface{}, response interface{}) {
	_m.Called(eventName, request, response)
}

// SetRequest provides a mock function with given fields: r
func (_m *ISentry) SetRequest(r *http.Request) {
	_m.Called(r)
}

// SetStartTransaction provides a mock function with given fields: ctx, spanName, transactionName, fn
func (_m *ISentry) SetStartTransaction(ctx context.Context, spanName string, transactionName string, fn func(context.Context, *sentry.Span) (string, uint8)) {
	_m.Called(ctx, spanName, transactionName, fn)
}

// SetTag provides a mock function with given fields: sentrySpan, name, value
func (_m *ISentry) SetTag(sentrySpan *sentry.Span, name string, value string) {
	_m.Called(sentrySpan, name, value)
}

// SetUserInfo provides a mock function with given fields: u
func (_m *ISentry) SetUserInfo(u go_lib_commonsentry.UserInfoSentry) {
	_m.Called(u)
}

// SpanContext provides a mock function with given fields: span
func (_m *ISentry) SpanContext(span sentry.Span) context.Context {
	ret := _m.Called(span)

	var r0 context.Context
	if rf, ok := ret.Get(0).(func(sentry.Span) context.Context); ok {
		r0 = rf(span)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(context.Context)
		}
	}

	return r0
}

// StartSpan provides a mock function with given fields: ctx, spanName
func (_m *ISentry) StartSpan(ctx context.Context, spanName string) *sentry.Span {
	ret := _m.Called(ctx, spanName)

	var r0 *sentry.Span
	if rf, ok := ret.Get(0).(func(context.Context, string) *sentry.Span); ok {
		r0 = rf(ctx, spanName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*sentry.Span)
		}
	}

	return r0
}

// Trace provides a mock function with given fields: ctx, spanName, fn
func (_m *ISentry) Trace(ctx context.Context, spanName string, fn func(context.Context, *sentry.Span)) {
	_m.Called(ctx, spanName, fn)
}

type mockConstructorTestingTNewISentry interface {
	mock.TestingT
	Cleanup(func())
}

// NewISentry creates a new instance of ISentry. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewISentry(t mockConstructorTestingTNewISentry) *ISentry {
	mock := &ISentry{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
