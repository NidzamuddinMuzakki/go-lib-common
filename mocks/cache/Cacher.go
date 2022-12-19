// Code generated by mockery v2.13.1. DO NOT EDIT.

package mocks

import (
	context "context"

	cache "bitbucket.org/moladinTech/go-lib-common/cache"

	mock "github.com/stretchr/testify/mock"

	redis "github.com/go-redis/redis/v8"

	time "time"
)

// Cacher is an autogenerated mock type for the Cacher type
type Cacher struct {
	mock.Mock
}

// BatchGet provides a mock function with given fields: ctx, keys, dest
func (_m *Cacher) BatchGet(ctx context.Context, keys []cache.Key, dest interface{}) error {
	ret := _m.Called(ctx, keys, dest)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []cache.Key, interface{}) error); ok {
		r0 = rf(ctx, keys, dest)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// BatchSet provides a mock function with given fields: ctx, datas, duration
func (_m *Cacher) BatchSet(ctx context.Context, datas []cache.Data, duration time.Duration) error {
	ret := _m.Called(ctx, datas, duration)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []cache.Data, time.Duration) error); ok {
		r0 = rf(ctx, datas, duration)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete provides a mock function with given fields: ctx, key
func (_m *Cacher) Delete(ctx context.Context, key cache.Key) error {
	ret := _m.Called(ctx, key)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, cache.Key) error); ok {
		r0 = rf(ctx, key)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Expire provides a mock function with given fields: ctx, key, ttl
func (_m *Cacher) Expire(ctx context.Context, key string, ttl time.Duration) (*redis.BoolCmd, error) {
	ret := _m.Called(ctx, key, ttl)

	var r0 *redis.BoolCmd
	if rf, ok := ret.Get(0).(func(context.Context, string, time.Duration) *redis.BoolCmd); ok {
		r0 = rf(ctx, key, ttl)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*redis.BoolCmd)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, time.Duration) error); ok {
		r1 = rf(ctx, key, ttl)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Get provides a mock function with given fields: ctx, key, dest
func (_m *Cacher) Get(ctx context.Context, key cache.Key, dest interface{}) error {
	ret := _m.Called(ctx, key, dest)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, cache.Key, interface{}) error); ok {
		r0 = rf(ctx, key, dest)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Incr provides a mock function with given fields: ctx, key
func (_m *Cacher) Incr(ctx context.Context, key string) (*redis.IntCmd, error) {
	ret := _m.Called(ctx, key)

	var r0 *redis.IntCmd
	if rf, ok := ret.Get(0).(func(context.Context, string) *redis.IntCmd); ok {
		r0 = rf(ctx, key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*redis.IntCmd)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Set provides a mock function with given fields: ctx, data, duration
func (_m *Cacher) Set(ctx context.Context, data cache.Data, duration time.Duration) error {
	ret := _m.Called(ctx, data, duration)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, cache.Data, time.Duration) error); ok {
		r0 = rf(ctx, data, duration)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewCacher interface {
	mock.TestingT
	Cleanup(func())
}

// NewCacher creates a new instance of Cacher. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewCacher(t mockConstructorTestingTNewCacher) *Cacher {
	mock := &Cacher{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
