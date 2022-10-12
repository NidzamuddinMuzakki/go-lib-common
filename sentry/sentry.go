package sentry

import (
	"context"
	"time"

	"bitbucket.org/moladinTech/go-lib-common/constant"
	commonContext "bitbucket.org/moladinTech/go-lib-common/context"

	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type function func(ctx context.Context) (string, uint8)
type Platform int
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
	Dsn        string
	Debug      bool
	Env        string
	SampleRate float64
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

type ISentry interface {
	SetStartTransaction(
		ctx context.Context,
		spanName string,
		transactionName string,
		fn function,
	)
	Trace(ctx context.Context, spanName string, fn function)
	StartSpan(ctx context.Context, spanName string) *sentry.Span
	Finish(span *sentry.Span)
	SetTag(sentrySpan *sentry.Span, name string, value string)
	CaptureException(exception error) *sentry.EventID
	GetGinMiddleware() gin.HandlerFunc
	Flush(timeout time.Duration) bool
	SetUserInfo(u UserInfoSentry)
	HandlingPanic(err interface{})
}

type Option func(*SentryPackage)

func NewSentry(
	options ...Option,
) (ISentry, error) {
	sentryPkg := &SentryPackage{}

	for _, option := range options {
		option(sentryPkg)
	}

	err := sentry.Init(sentry.ClientOptions{
		Dsn:              sentryPkg.Dsn,
		Debug:            sentryPkg.Debug,
		Environment:      sentryPkg.Env,
		TracesSampleRate: sentryPkg.SampleRate,
	})
	if err != nil {
		return nil, err
	}
	return sentryPkg, nil
}

// SetStartTransaction is used to describe each transaction
func (s *SentryPackage) SetStartTransaction(
	ctx context.Context,
	spanName string,
	transactionName string,
	fn function,
) {
	span := sentry.StartSpan(ctx, spanName, sentry.TransactionName(transactionName))
	defer span.Finish()
	xRequestID := commonContext.GetValueAsString(ctx, constant.XRequestIdHeader)
	span.TraceID = sentry.TraceID(uuid.MustParse(xRequestID))
	span.SetTag(constant.XRequestIdHeader, xRequestID)

	status, spanStatus := fn(span.Context())
	span.SetTag("real_response_status", status)
	span.Status = sentry.SpanStatus(spanStatus)
}

// Trace is used to describe each operation
func (s *SentryPackage) Trace(ctx context.Context, spanName string, fn function) {
	span := sentry.StartSpan(ctx, spanName)
	defer span.Finish()
	fn(span.Context())
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
	return sentry.Flush(2 * time.Second)
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
