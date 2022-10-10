package gin

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

type MiddlewarePanicRecovery interface {
	PanicRecoveryMiddleware() gin.HandlerFunc
}

type middlewarePanicRecovery struct {
	configEnv string
}

func NewPanicRecovery(configEnv string) MiddlewarePanicRecovery {
	return &middlewarePanicRecovery{
		configEnv: configEnv,
	}
}

func (p *middlewarePanicRecovery) PanicRecoveryMiddleware() gin.HandlerFunc {
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
				sentry.HandlingPanic(pnc)
				c.AbortWithStatusJSON(http.StatusInternalServerError, response.Response{Status: response.StatusFail, Message: responseMsg})
				return
			}
		}()
		c.Next()
	}
}
