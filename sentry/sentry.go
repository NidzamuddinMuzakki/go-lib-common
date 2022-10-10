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

type Config struct {
	Dsn        string
	Debug      bool
	Env        string
	SampleRate float64
}

func Init(ctx context.Context, config Config) error {
	err := sentry.Init(sentry.ClientOptions{
		Dsn:              config.Dsn,
		Debug:            config.Debug,
		Environment:      config.Env,
		TracesSampleRate: config.SampleRate,
	})
	if err != nil {
		return err
	}
	return nil
}

// SetStartTransaction is used to describe each transaction
func SetStartTransaction(
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
func Trace(ctx context.Context, spanName string, fn function) {
	span := sentry.StartSpan(ctx, spanName)
	defer span.Finish()
	fn(span.Context())
}

func StartSpan(ctx context.Context, spanName string) *sentry.Span {
	return sentry.StartSpan(ctx, spanName)
}

func Finish(span *sentry.Span) {
	span.Finish()
}

func SetTag(sentrySpan *sentry.Span, name string, value string) {
	sentrySpan.SetTag(name, value)
}

func CaptureException(exception error) *sentry.EventID {
	return sentry.CaptureException(exception)
}

func GetGinMiddleware() gin.HandlerFunc {
	return sentrygin.New(sentrygin.Options{
		Repanic: true,
	})
}

func Flush(timeout time.Duration) bool {
	return sentry.Flush(2 * time.Second)
}

// SetUserInfo is used to describe user information
func SetUserInfo(u UserInfoSentry) {
	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetUser(sentry.User{
			ID:       u.ID,
			Username: u.Username,
			Email:    u.Email,
		})
	})
}

func HandlingPanic(err interface{}) {
	sentry.CurrentHub().Recover(err)
	sentry.Flush(time.Second * 5)
}
