// Code generated by mockery v2.20.2. DO NOT EDIT.

package mocks

import (
	aws "bitbucket.org/moladinTech/go-lib-common/client/aws"
	auth "bitbucket.org/moladinTech/go-lib-common/middleware/gin/auth"

	cache "bitbucket.org/moladinTech/go-lib-common/cache"

	encryption "bitbucket.org/moladinTech/go-lib-common/encryption"

	exporter "bitbucket.org/moladinTech/go-lib-common/exporter"

	gcp "bitbucket.org/moladinTech/go-lib-common/client/gcp"

	limiter "bitbucket.org/moladinTech/go-lib-common/middleware/gin/limiter"

	mock "github.com/stretchr/testify/mock"

	moladin_evo "bitbucket.org/moladinTech/go-lib-common/client/moladin_evo"

	notification "bitbucket.org/moladinTech/go-lib-common/client/notification"

	panic_recovery "bitbucket.org/moladinTech/go-lib-common/middleware/gin/panic_recovery"

	sentry "bitbucket.org/moladinTech/go-lib-common/sentry"

	slack "bitbucket.org/moladinTech/go-lib-common/client/notification/slack"

	time "bitbucket.org/moladinTech/go-lib-common/time"

	tracer "bitbucket.org/moladinTech/go-lib-common/middleware/gin/tracer"

	validator "github.com/go-playground/validator/v10"
)

// IRegistry is an autogenerated mock type for the IRegistry type
type IRegistry struct {
	mock.Mock
}

// GetAuthMiddleware provides a mock function with given fields:
func (_m *IRegistry) GetAuthMiddleware() auth.IMiddlewareAuth {
	ret := _m.Called()

	var r0 auth.IMiddlewareAuth
	if rf, ok := ret.Get(0).(func() auth.IMiddlewareAuth); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(auth.IMiddlewareAuth)
		}
	}

	return r0
}

// GetCache provides a mock function with given fields:
func (_m *IRegistry) GetCache() cache.Cacher {
	ret := _m.Called()

	var r0 cache.Cacher
	if rf, ok := ret.Get(0).(func() cache.Cacher); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(cache.Cacher)
		}
	}

	return r0
}

// GetEncryption provides a mock function with given fields:
func (_m *IRegistry) GetEncryption() encryption.IEncryption {
	ret := _m.Called()

	var r0 encryption.IEncryption
	if rf, ok := ret.Get(0).(func() encryption.IEncryption); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(encryption.IEncryption)
		}
	}

	return r0
}

// GetExporterCSV provides a mock function with given fields:
func (_m *IRegistry) GetExporterCSV() exporter.Exporter {
	ret := _m.Called()

	var r0 exporter.Exporter
	if rf, ok := ret.Get(0).(func() exporter.Exporter); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(exporter.Exporter)
		}
	}

	return r0
}

// GetExporterExcel provides a mock function with given fields:
func (_m *IRegistry) GetExporterExcel() exporter.Exporter {
	ret := _m.Called()

	var r0 exporter.Exporter
	if rf, ok := ret.Get(0).(func() exporter.Exporter); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(exporter.Exporter)
		}
	}

	return r0
}

// GetGCS provides a mock function with given fields:
func (_m *IRegistry) GetGCS() gcp.GCSClient {
	ret := _m.Called()

	var r0 gcp.GCSClient
	if rf, ok := ret.Get(0).(func() gcp.GCSClient); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(gcp.GCSClient)
		}
	}

	return r0
}

// GetLimiterMiddleware provides a mock function with given fields:
func (_m *IRegistry) GetLimiterMiddleware() limiter.IMiddlewareLimiter {
	ret := _m.Called()

	var r0 limiter.IMiddlewareLimiter
	if rf, ok := ret.Get(0).(func() limiter.IMiddlewareLimiter); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(limiter.IMiddlewareLimiter)
		}
	}

	return r0
}

// GetMoladinEvo provides a mock function with given fields:
func (_m *IRegistry) GetMoladinEvo() moladin_evo.IMoladinEvo {
	ret := _m.Called()

	var r0 moladin_evo.IMoladinEvo
	if rf, ok := ret.Get(0).(func() moladin_evo.IMoladinEvo); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(moladin_evo.IMoladinEvo)
		}
	}

	return r0
}

// GetNotif provides a mock function with given fields:
func (_m *IRegistry) GetNotif() notification.INotification {
	ret := _m.Called()

	var r0 notification.INotification
	if rf, ok := ret.Get(0).(func() notification.INotification); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(notification.INotification)
		}
	}

	return r0
}

// GetPanicRecoveryMiddleware provides a mock function with given fields:
func (_m *IRegistry) GetPanicRecoveryMiddleware() panic_recovery.IMiddlewarePanicRecovery {
	ret := _m.Called()

	var r0 panic_recovery.IMiddlewarePanicRecovery
	if rf, ok := ret.Get(0).(func() panic_recovery.IMiddlewarePanicRecovery); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(panic_recovery.IMiddlewarePanicRecovery)
		}
	}

	return r0
}

// GetS3 provides a mock function with given fields:
func (_m *IRegistry) GetS3() aws.S3 {
	ret := _m.Called()

	var r0 aws.S3
	if rf, ok := ret.Get(0).(func() aws.S3); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(aws.S3)
		}
	}

	return r0
}

// GetSentry provides a mock function with given fields:
func (_m *IRegistry) GetSentry() sentry.ISentry {
	ret := _m.Called()

	var r0 sentry.ISentry
	if rf, ok := ret.Get(0).(func() sentry.ISentry); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(sentry.ISentry)
		}
	}

	return r0
}

// GetSlack provides a mock function with given fields:
func (_m *IRegistry) GetSlack() slack.ISlack {
	ret := _m.Called()

	var r0 slack.ISlack
	if rf, ok := ret.Get(0).(func() slack.ISlack); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(slack.ISlack)
		}
	}

	return r0
}

// GetTime provides a mock function with given fields:
func (_m *IRegistry) GetTime() time.TimeItf {
	ret := _m.Called()

	var r0 time.TimeItf
	if rf, ok := ret.Get(0).(func() time.TimeItf); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(time.TimeItf)
		}
	}

	return r0
}

// GetTraceMiddleware provides a mock function with given fields:
func (_m *IRegistry) GetTraceMiddleware() tracer.IMiddlewareTracer {
	ret := _m.Called()

	var r0 tracer.IMiddlewareTracer
	if rf, ok := ret.Get(0).(func() tracer.IMiddlewareTracer); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(tracer.IMiddlewareTracer)
		}
	}

	return r0
}

// GetValidator provides a mock function with given fields:
func (_m *IRegistry) GetValidator() *validator.Validate {
	ret := _m.Called()

	var r0 *validator.Validate
	if rf, ok := ret.Get(0).(func() *validator.Validate); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*validator.Validate)
		}
	}

	return r0
}

type mockConstructorTestingTNewIRegistry interface {
	mock.TestingT
	Cleanup(func())
}

// NewIRegistry creates a new instance of IRegistry. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewIRegistry(t mockConstructorTestingTNewIRegistry) *IRegistry {
	mock := &IRegistry{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
