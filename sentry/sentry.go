//go:generate mockery --name=ISentry
package sentry

import (
	"context"
	"net/http"
	"time"

	"golang.org/x/exp/slices"

	"bitbucket.org/moladinTech/go-lib-common/constant"
	commonContext "bitbucket.org/moladinTech/go-lib-common/context"
	commonValidator "bitbucket.org/moladinTech/go-lib-common/validator"

	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type UserInfoSentry struct {
	ID       string
	Username string
	Email    string
}

const (
	STATUS_OK                    = sentry.SpanStatusOK
	STATUS_INVALID_ARGUMENT      = sentry.SpanStatusInvalidArgument
	STATUS_INTERNAL_SERVER_ERROR = sentry.SpanStatusInternalError
)

type SentryPackage struct {
	Dsn                   string  `validate:"required"`
	Env                   string  `validate:"required"`
	SampleRate            float64 `validate:"required"`
	EnableTracing         bool
	BlacklistTransactions []string
	Debug                 bool
}

func WithDsn(dsn string) Option {
	return func(s *SentryPackage) {
		s.Dsn = dsn
	}
}
func WithDebug(debug bool) Option {
	return func(s *SentryPackage) {
		s.Debug = debug
	}
}
func WithEnv(env string) Option {
	return func(s *SentryPackage) {
		s.Env = env
	}
}
func WithSampleRate(sampleRate float64) Option {
	return func(s *SentryPackage) {
		s.SampleRate = sampleRate
	}
}
func WithEnableTracing(enableTracing bool) Option {
	return func(s *SentryPackage) {
		s.EnableTracing = enableTracing
	}
}
func WithBlacklistTransactions(TransactionNames []string) Option {
	return func(s *SentryPackage) {
		s.BlacklistTransactions = TransactionNames
	}
}

type ISentry interface {
	SetStartTransaction(
		ctx context.Context,
		spanName string,
		transactionName string,
		fn func(ctx context.Context, span *sentry.Span) (string, uint8),
	)
	Trace(ctx context.Context, spanName string, fn func(ctx context.Context, span *sentry.Span))
	StartSpan(ctx context.Context, spanName string) *sentry.Span
	Finish(span *sentry.Span)
	SetTag(sentrySpan *sentry.Span, name string, value string)
	CaptureException(exception error) *sentry.EventID
	GetGinMiddleware() gin.HandlerFunc
	Flush(timeout time.Duration) bool
	SetUserInfo(u UserInfoSentry)
	HandlingPanic(err interface{})
	SpanContext(span sentry.Span) context.Context
	SetRequest(r *http.Request)
	SetIntegrationCapture(
		eventName string,
		request interface{},
		response interface{},
	)
	SetEventCapture(
		eventName string,
		data interface{},
	)
}

type Option func(*SentryPackage)

func NewSentry(
	validator *validator.Validate,
	options ...Option,
) ISentry {
	sentryPkg := &SentryPackage{
		EnableTracing: true,
	}

	for _, option := range options {
		option(sentryPkg)
	}
	err := validator.Struct(sentryPkg)
	if err != nil {
		panic(commonValidator.ToErrResponse(err))
	}

	err = sentry.Init(sentry.ClientOptions{
		EnableTracing:    sentryPkg.EnableTracing,
		Dsn:              sentryPkg.Dsn,
		Debug:            sentryPkg.Debug,
		Environment:      sentryPkg.Env,
		TracesSampleRate: sentryPkg.SampleRate,
		BeforeSendTransaction: func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
			if len(sentryPkg.BlacklistTransactions) > 0 && slices.Contains(sentryPkg.BlacklistTransactions, event.Transaction) {
				return nil
			}

			return event
		},
		BeforeSend: func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
			if len(sentryPkg.BlacklistTransactions) > 0 && slices.Contains(sentryPkg.BlacklistTransactions, event.Transaction) {
				return nil
			}

			return event
		},
	})
	if err != nil {
		panic(err)
	}

	return sentryPkg
}

// SetStartTransaction is used to describe each transaction
func (s *SentryPackage) SetStartTransaction(
	ctx context.Context,
	spanName string,
	transactionName string,
	fn func(ctx context.Context, span *sentry.Span) (string, uint8),
) {
	span := sentry.StartSpan(ctx, spanName, sentry.TransactionName(transactionName))
	defer span.Finish()
	xRequestID := commonContext.GetValueAsString(ctx, constant.XRequestIdHeader)
	span.TraceID = sentry.TraceID(uuid.MustParse(xRequestID))
	span.SetTag(constant.XRequestIdHeader, xRequestID)

	status, spanStatus := fn(span.Context(), span)
	span.SetTag("real_response_status", status)
	span.Status = sentry.SpanStatus(spanStatus)
}

// Trace is used to describe each operation
func (s *SentryPackage) Trace(ctx context.Context, spanName string, fn func(ctx context.Context, span *sentry.Span)) {
	span := sentry.StartSpan(ctx, spanName)
	defer span.Finish()
	fn(span.Context(), span)
}

func (s *SentryPackage) StartSpan(ctx context.Context, spanName string) *sentry.Span {
	return sentry.StartSpan(ctx, spanName)
}

func (s *SentryPackage) Finish(span *sentry.Span) {
	span.Finish()
}

func (s *SentryPackage) SetTag(sentrySpan *sentry.Span, name string, value string) {
	sentrySpan.SetTag(name, value)
}

func (s *SentryPackage) CaptureException(exception error) *sentry.EventID {
	return sentry.CaptureException(exception)
}

func (s *SentryPackage) GetGinMiddleware() gin.HandlerFunc {
	return sentrygin.New(sentrygin.Options{
		Repanic: true,
	})
}

func (s *SentryPackage) Flush(timeout time.Duration) bool {
	return sentry.Flush(timeout * time.Second)
}

// SetUserInfo is used to describe user information
func (s *SentryPackage) SetUserInfo(u UserInfoSentry) {
	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetUser(sentry.User{
			ID:       u.ID,
			Username: u.Username,
			Email:    u.Email,
		})
	})
}

func (s *SentryPackage) HandlingPanic(err interface{}) {
	sentry.CurrentHub().Recover(err)
	sentry.Flush(time.Second * 5)
}

func Span() sentry.Span {
	return sentry.Span{}
}

func (s *SentryPackage) SpanContext(span sentry.Span) context.Context {
	return span.Context()
}

// SetRequest add request to the current scope.
func (s *SentryPackage) SetRequest(r *http.Request) {
	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetRequest(r)
	})
}

// SetIntegrationCapture event capturing for specific request and response data
func (s *SentryPackage) SetIntegrationCapture(
	eventName string,
	request interface{},
	response interface{},
) {
	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetContext(eventName, map[string]interface{}{
			"request":  request,
			"response": response,
		})
	})
}

// SetIntegrationCapture event capturing for specific data
func (s *SentryPackage) SetEventCapture(
	eventName string,
	data interface{},
) {
	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetContext(eventName, map[string]interface{}{
			"data": data,
		})
	})
}
