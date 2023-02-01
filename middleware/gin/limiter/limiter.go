package limiter

import (
	"fmt"
	"net/http"
	"time"

	"bitbucket.org/moladinTech/go-lib-common/cache"
	"bitbucket.org/moladinTech/go-lib-common/cast"
	"bitbucket.org/moladinTech/go-lib-common/logger"
	"bitbucket.org/moladinTech/go-lib-common/response"
	commonValidator "bitbucket.org/moladinTech/go-lib-common/validator"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type IMiddlewareLimiter interface {
	WithCustomLimit(rt RateLimit) gin.HandlerFunc
}

type MiddlewareLimitPackage struct {
	Cache        cache.Cacher `validate:"required"`
	ServiceName  string       `validate:"required"`
	DefaultTTL   uint         `validate:"required"`
	DefaultLimit uint         `validate:"required"`
}

type Option func(*MiddlewareLimitPackage)

func WithCacher(cache cache.Cacher) Option {
	return func(s *MiddlewareLimitPackage) {
		s.Cache = cache
	}
}

func WithServiceName(serviceName string) Option {
	return func(s *MiddlewareLimitPackage) {
		s.ServiceName = serviceName
	}
}

func NewLimiter(validator *validator.Validate, options ...Option) (*MiddlewareLimitPackage, gin.HandlerFunc) {
	mlp := &MiddlewareLimitPackage{
		DefaultTTL:   1,
		DefaultLimit: 100,
	}

	for _, option := range options {
		option(mlp)
	}
	err := validator.Struct(mlp)
	if err != nil {
		panic(commonValidator.ToErrResponse(err))
	}

	hf := func(c *gin.Context) {
		var (
			logCtx = "middleware.gin.limiter.limiter.Limit"
			ctx    = c.Request.Context()

			key = fmt.Sprintf("%s:%s", mlp.ServiceName, c.FullPath())
		)

		incr, err := mlp.Cache.Incr(ctx, key)
		if err != nil {
			logger.Error(ctx, "error", err, logger.Tag{Key: "logCtx", Value: logCtx})
		}

		if incr.Val() == 1 {
			_, err := mlp.Cache.Expire(ctx, key, time.Duration(mlp.DefaultTTL)*time.Second)
			if err != nil {
				logger.Error(ctx, "error", err, logger.Tag{Key: "logCtx", Value: logCtx})
			}
		}

		if incr.Val() > int64(mlp.DefaultLimit) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, response.Response{
				Status:  response.StatusFail,
				Message: http.StatusText(http.StatusTooManyRequests),
			})
			return
		}

		c.Next()
	}

	return mlp, hf
}

type RateLimit struct {
	TTL   *uint
	Limit *uint
}

// WithCustomLimit custom rate limit
// With default rate limit is 100 requests per seconds
func (a *MiddlewareLimitPackage) WithCustomLimit(rt RateLimit) gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			logCtx = "middleware.gin.limiter.limiter.WithCustomLimit"
			ctx    = c.Request.Context()

			key = fmt.Sprintf("%s:%s", a.ServiceName, c.FullPath())
		)

		if rt.TTL == nil {
			rt.TTL = cast.NewPointer[uint](a.DefaultTTL)
		}

		if rt.Limit == nil {
			rt.Limit = cast.NewPointer[uint](a.DefaultLimit)
		}

		incr, err := a.Cache.Incr(ctx, key)
		if err != nil {
			logger.Error(ctx, "error", err, logger.Tag{Key: "logCtx", Value: logCtx})
		}

		if incr.Val() == 1 {
			_, err := a.Cache.Expire(ctx, key, time.Duration(*rt.TTL)*time.Second)
			if err != nil {
				logger.Error(ctx, "error", err, logger.Tag{Key: "logCtx", Value: logCtx})
			}
		}

		if incr.Val() > int64(*rt.Limit) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, response.Response{
				Status:  response.StatusFail,
				Message: http.StatusText(http.StatusTooManyRequests),
			})
			return
		}

		c.Next()
	}
}
