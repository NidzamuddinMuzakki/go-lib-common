//go:generate mockery --name=IMiddlewarePanicRecovery
package panic_recovery

import (
	"errors"
	"fmt"
	"net/http"

	"bitbucket.org/moladinTech/go-lib-common/constant"
	"bitbucket.org/moladinTech/go-lib-common/logger"
	"bitbucket.org/moladinTech/go-lib-common/response"
	"bitbucket.org/moladinTech/go-lib-common/sentry"
	commonValidator "bitbucket.org/moladinTech/go-lib-common/validator"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type IMiddlewarePanicRecovery interface {
	PanicRecoveryMiddleware() gin.HandlerFunc
}

type MiddlewarePanicRecoveryPackage struct {
	ConfigEnv string         `validate:"required"`
	Sentry    sentry.ISentry `validate:"required"`
}

func WithConfigEnv(configEnv string) Option {
	return func(s *MiddlewarePanicRecoveryPackage) {
		s.ConfigEnv = configEnv
	}
}
func WithSentry(sentry sentry.ISentry) Option {
	return func(s *MiddlewarePanicRecoveryPackage) {
		s.Sentry = sentry
	}
}

type Option func(*MiddlewarePanicRecoveryPackage)

func NewPanicRecovery(
	validator *validator.Validate,
	options ...Option,
) IMiddlewarePanicRecovery {
	middlewarePanicRecoveryPackage := &MiddlewarePanicRecoveryPackage{}

	for _, option := range options {
		option(middlewarePanicRecoveryPackage)
	}

	err := validator.Struct(middlewarePanicRecoveryPackage)
	if err != nil {
		panic(commonValidator.ToErrResponse(err))
	}

	return middlewarePanicRecoveryPackage
}

func (p *MiddlewarePanicRecoveryPackage) PanicRecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			reqCtx := c.Request.Context()
			const logCtx = "common.middleware.gin.panic_recovery.PanicRecoveryMiddleware"
			span := p.Sentry.StartSpan(reqCtx, logCtx)
			defer span.Finish()

			if pnc := recover(); pnc != nil {
				errStr := fmt.Sprintf("panic: %v", pnc)
				logger.Error(reqCtx, logCtx, errors.New(errStr))
				responseMsg := "Server error. Contact admin for more information."
				if p.ConfigEnv != constant.EnvProduction {
					responseMsg = errStr
				}
				p.Sentry.HandlingPanic(pnc)
				c.AbortWithStatusJSON(http.StatusInternalServerError, response.Response{Status: response.StatusFail, Message: responseMsg})
				return
			}
		}()
		c.Next()
	}
}
