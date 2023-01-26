package limiter

import (
	"context"
	"time"

	"bitbucket.org/moladinTech/go-lib-common/cache"
	commonValidator "bitbucket.org/moladinTech/go-lib-common/validator"
	"github.com/go-playground/validator/v10"
)

type SlackLimiter struct {
	cache   cache.Cacher `validate:"required"`
	seconds *uint        `validate:"required"`
}

type ISlackLimiter interface {
	LimitChecker(ctx context.Context, data cache.Data) (isSuccessSet bool, err error)
}

type Option func(*SlackLimiter)

func WithCache(Cache cache.Cacher) Option {
	return func(s *SlackLimiter) {
		s.cache = Cache
	}
}

func WithTTL(seconds *uint) Option {
	return func(s *SlackLimiter) {
		var defaultSeconds uint = 600
		// if empty, set default value to 10 minutes
		if seconds == nil {
			seconds = &defaultSeconds
			return
		} else {
			if *seconds <= 0 {
				seconds = &defaultSeconds
			}
			s.seconds = seconds
		}
	}
}

func NewSlackLimiter(validator *validator.Validate, options ...Option) *SlackLimiter {
	slackLimiterPackages := &SlackLimiter{}

	for _, option := range options {
		option(slackLimiterPackages)
	}
	err := validator.Struct(slackLimiterPackages)
	if err != nil {
		panic(commonValidator.ToErrResponse(err))
	}

	return slackLimiterPackages
}

// LimitChecker
// if isSuccessSet = true, it means setNX successfully set new key to redis.
// that means the data.Key didnt exist before.
func (c *SlackLimiter) LimitChecker(ctx context.Context, data cache.Data) (isSuccessSet bool, err error) {
	ttl := time.Duration(*c.seconds) * time.Second
	isSuccessSet, err = c.cache.SetNx(ctx, data, ttl)
	if err != nil {
		return false, err
	}

	return isSuccessSet, nil
}
