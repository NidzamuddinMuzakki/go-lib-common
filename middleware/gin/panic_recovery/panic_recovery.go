//go:generate mockery --name=IMiddlewarePanicRecovery
package panic_recovery

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"

	"bitbucket.org/moladinTech/go-lib-common/client/notification/slack"
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
	Slack     slack.ISlack   `validate:"required"`
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

func WithSlack(slack slack.ISlack) Option {
	return func(s *MiddlewarePanicRecoveryPackage) {
		s.Slack = slack
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

func (p *MiddlewarePanicRecoveryPackage) sendSlack(
	ctx context.Context,
	err error,
) {
	const logCtx = "common.middleware.gin.panic_recovery.sendSlack"
	span := p.Sentry.StartSpan(ctx, logCtx)
	ctx = p.Sentry.SpanContext(*span)
	defer p.Sentry.Finish(span)

	slackMessage := p.Slack.GetFormattedMessage(logCtx, ctx, err)
	errSlack := p.Slack.Send(ctx, slackMessage)
	if errSlack != nil {
		p.Sentry.CaptureException(errSlack)
		logger.Error(ctx, "Error sending notif to slack", err,
			logger.Tag{
				Key:   "logCtx",
				Value: logCtx,
			},
			logger.Tag{
				Key:   "slackError",
				Value: err,
			},
		)
	}
}

func (p *MiddlewarePanicRecoveryPackage) PanicRecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			reqCtx := c.Request.Context()
			const logCtx = "common.middleware.gin.panic_recovery.PanicRecoveryMiddleware"
			span := p.Sentry.StartSpan(reqCtx, logCtx)
			ctx := p.Sentry.SpanContext(*span)
			defer p.Sentry.Finish(span)

			if pnc := recover(); pnc != nil {
				errStr := fmt.Sprintf("panic: %v", pnc)
				p.Sentry.CaptureException(errors.New(errStr))
				logger.Error(ctx, logCtx, errors.New(errStr), logger.Tag{Key: "debug", Value: string(debug.Stack())})
				responseMsg := "Server error. Contact admin for more information."
				if p.ConfigEnv != constant.EnvProduction {
					responseMsg = errStr
				}
				p.Sentry.HandlingPanic(pnc)
				p.sendSlack(ctx, errors.New(errStr))
				c.AbortWithStatusJSON(
					http.StatusInternalServerError,
					response.Response{
						Status:  response.StatusFail,
						Message: responseMsg,
					},
				)
				return
			}
		}()
		c.Next()
	}
}
