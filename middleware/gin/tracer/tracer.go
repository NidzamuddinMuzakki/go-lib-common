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
	commonSentry "bitbucket.org/moladinTech/go-lib-common/sentry"

	"github.com/gin-gonic/gin"
)

type IMiddlewareTracer interface {
	Tracer() gin.HandlerFunc
}

type MiddlewareTracerPackage struct {
	sentry commonSentry.ISentry
}

func WithSentry(sentry commonSentry.ISentry) Option {
	return func(s *MiddlewareTracerPackage) {
		s.sentry = sentry
	}
}

type Option func(*MiddlewareTracerPackage)

func NewTracer(
	options ...Option,
) IMiddlewareTracer {
	middlewareTracerPackage := &MiddlewareTracerPackage{}
	for _, option := range options {
		option(middlewareTracerPackage)
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
		// trace when request is getting started
		start := time.Now()
		ctxReq := logger.AddLoggingTag(c.Request.Context(), logger.Tag{Key: tagPath, Value: c.Request.URL.Path})
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

		s.sentry.SetStartTransaction(
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
					statusSpan = uint8(commonSentry.STATUS_OK)
				case 4:
					statusSpan = uint8(commonSentry.STATUS_INVALID_ARGUMENT)
				case 5:
					statusSpan = uint8(commonSentry.STATUS_INTERNAL_SERVER_ERROR)
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
