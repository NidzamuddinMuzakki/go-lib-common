//go:generate mockery --name=IMiddlewareTracer
package tracer

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"bitbucket.org/moladinTech/go-lib-common/logger"
	"bitbucket.org/moladinTech/go-lib-common/sentry"
	commonValidator "bitbucket.org/moladinTech/go-lib-common/validator"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type IMiddlewareTracer interface {
	Tracer() gin.HandlerFunc
}

type MiddlewareTracerPackage struct {
	Sentry sentry.ISentry `validate:"required"`
}

func WithSentry(sentry sentry.ISentry) Option {
	return func(s *MiddlewareTracerPackage) {
		s.Sentry = sentry
	}
}

type Option func(*MiddlewareTracerPackage)

func NewTracer(
	validator *validator.Validate,
	options ...Option,
) IMiddlewareTracer {
	middlewareTracerPackage := &MiddlewareTracerPackage{}
	for _, option := range options {
		option(middlewareTracerPackage)
	}

	err := validator.Struct(middlewareTracerPackage)
	if err != nil {
		panic(commonValidator.ToErrResponse(err))
	}
	return middlewareTracerPackage
}

func (s *MiddlewareTracerPackage) Tracer() gin.HandlerFunc {

	const (
		tagLatency = "latency"
		tagPath    = "path"
		tagQuery   = "query"
		tagBody    = "body"
		tagHeader  = "header"
		tagType    = "type"
		tagMethod  = "method"
	)
	return func(c *gin.Context) {
		const logCtx = "common.middleware.gin.tracer.Tracer"
		reqCtx := c.Request.Context()
		span := s.Sentry.StartSpan(reqCtx, logCtx)
		defer span.Finish()
		// trace when request is getting started
		start := time.Now()
		ctxReq := logger.AddLoggingTag(reqCtx, logger.Tag{Key: tagPath, Value: c.Request.URL.Path})
		ctxReq = logger.AddLoggingTag(ctxReq, logger.Tag{Key: tagMethod, Value: c.Request.Method})
		ctxReq = logger.AddLoggingTag(ctxReq, logger.Tag{Key: tagHeader, Value: c.Request.Header})
		if c.Request.Method != http.MethodGet {
			ctxReq = logger.AddLoggingTag(ctxReq, logger.Tag{Key: tagBody, Value: func() string {
				// Read the content
				var bodyBytes []byte
				var err error
				if c.Request.Body != nil {
					bodyBytes, err = ioutil.ReadAll(c.Request.Body)
					if err != nil {
						logger.Error(ctxReq, `error`, err)
					}
				}
				// Restore the io.ReadCloser to its original state
				c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

				// encode to base64
				encodedData := make([]byte, base64.StdEncoding.EncodedLen(len(bodyBytes)))
				base64.StdEncoding.Encode(encodedData, bodyBytes)

				return string(bodyBytes)
			}()})
		}

		logger.Info(ctxReq, `request`)

		s.Sentry.SetStartTransaction(
			c.Request.Context(),
			"middleware.Tracer",
			fmt.Sprintf("%s %s", c.Request.Method, c.FullPath()),
			func(ctx context.Context) (string, uint8) {
				c.Request = c.Request.WithContext(ctx)
				c.Next()
				var statusSpan uint8
				status := fmt.Sprint(c.Writer.Status())

				switch c.Writer.Status() / 100 {
				case 2:
					statusSpan = uint8(sentry.STATUS_OK)
				case 4:
					statusSpan = uint8(sentry.STATUS_INVALID_ARGUMENT)
				case 5:
					statusSpan = uint8(sentry.STATUS_INTERNAL_SERVER_ERROR)
				}
				return status, statusSpan
			},
		)

		// trace after request ended and log response with logger

		latency := time.Since(start)
		ctxResp := logger.AddLoggingTag(c.Request.Context(), logger.Tag{Key: tagLatency, Value: latency})
		ctxResp = logger.AddLoggingTag(ctxResp, logger.Tag{Key: tagMethod, Value: c.Request.Method})
		ctxResp = logger.AddLoggingTag(ctxResp, logger.Tag{Key: tagPath, Value: c.Request.URL.Path})
		ctxResp = logger.AddLoggingTag(ctxResp, logger.Tag{Key: tagQuery, Value: c.Request.URL.RawQuery})
		ctxResp = logger.AddLoggingTag(ctxResp, logger.Tag{Key: tagHeader, Value: c.Writer.Header()})
		logger.Info(ctxResp, `response`)

	}
}