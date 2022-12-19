package limiter

import (
	"bitbucket.org/moladinTech/go-lib-common/cache"
	"bitbucket.org/moladinTech/go-lib-common/response"
	commonValidator "bitbucket.org/moladinTech/go-lib-common/validator"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"time"
)

type IMiddlewareLimiter interface {
	Limit(key string, ttl time.Duration, limit uint) gin.HandlerFunc
}

type MiddlewareLimitPackage struct {
	Cache cache.Cacher `validate:"required"`
}

type Option func(*MiddlewareLimitPackage)

func WithCacher(cache cache.Cacher) Option {
	return func(s *MiddlewareLimitPackage) {
		s.Cache = cache
	}
}

func NewLimiter(validator *validator.Validate, options ...Option) IMiddlewareLimiter {
	middlewareLimitPackage := &MiddlewareLimitPackage{}

	for _, option := range options {
		option(middlewareLimitPackage)
	}
	err := validator.Struct(middlewareLimitPackage)
	if err != nil {
		panic(commonValidator.ToErrResponse(err))
	}

	return middlewareLimitPackage
}

func (a *MiddlewareLimitPackage) Limit(key string, ttl time.Duration, limit uint) gin.HandlerFunc {
	return func(c *gin.Context) {

		incr, err := a.Cache.Incr(c.Request.Context(), key)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, response.Response{
				Status:  response.StatusError,
				Message: http.StatusText(http.StatusInternalServerError),
			})
			return
		}

		if incr.Val() == 1 {
			_, err := a.Cache.Expire(c.Request.Context(), key, ttl)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, response.Response{
					Status:  response.StatusError,
					Message: http.StatusText(http.StatusInternalServerError),
				})
				return
			}
		}

		if incr.Val() > int64(limit) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, response.Response{
				Status:  response.StatusFail,
				Message: http.StatusText(http.StatusTooManyRequests),
			})
			return
		}

		c.Next()
	}
}
