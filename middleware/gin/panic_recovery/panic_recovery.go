package panic_recovery

import (
	"errors"
	"fmt"
	"net/http"

	"bitbucket.org/moladinTech/go-lib-common/constant"
	"bitbucket.org/moladinTech/go-lib-common/logger"
	"bitbucket.org/moladinTech/go-lib-common/response"
	"bitbucket.org/moladinTech/go-lib-common/sentry"
	"github.com/gin-gonic/gin"
)

type IMiddlewarePanicRecovery interface {
	PanicRecoveryMiddleware() gin.HandlerFunc
}

type MiddlewarePanicRecoveryPackage struct {
	configEnv string
	sentry    sentry.ISentry
}

func WithConfigEnv(configEnv string) Option {
	return func(s *MiddlewarePanicRecoveryPackage) {
		s.configEnv = configEnv
	}
}
func WithSentry(sentry sentry.ISentry) Option {
	return func(s *MiddlewarePanicRecoveryPackage) {
		s.sentry = sentry
	}
}

type Option func(*MiddlewarePanicRecoveryPackage)

func NewPanicRecovery(
	options ...Option,
) IMiddlewarePanicRecovery {
	middlewarePanicRecoveryPackage := &MiddlewarePanicRecoveryPackage{}

	for _, option := range options {
		option(middlewarePanicRecoveryPackage)
	}

	return middlewarePanicRecoveryPackage
}

func (p *MiddlewarePanicRecoveryPackage) PanicRecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if pnc := recover(); pnc != nil {
				const logCtx = "middleware.gin.PanicRecoveryMiddleware"
				errStr := fmt.Sprintf("panic: %v", pnc)

				reqCtx := c.Request.Context()

				logger.Error(reqCtx, logCtx, errors.New(errStr))
				responseMsg := "Server error. Contact admin for more information."
				if p.configEnv != constant.EnvProduction {
					responseMsg = errStr
				}
				p.sentry.HandlingPanic(pnc)
				c.AbortWithStatusJSON(http.StatusInternalServerError, response.Response{Status: response.StatusFail, Message: responseMsg})
				return
			}
		}()
		c.Next()
	}
}
