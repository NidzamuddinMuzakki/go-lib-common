package sentry_test

import (
	"context"
	"errors"
	"log"
	"os"
	"testing"

	"bitbucket.org/moladinTech/go-lib-common/constant"
	commonSentry "bitbucket.org/moladinTech/go-lib-common/sentry"
	"bitbucket.org/moladinTech/go-lib-common/validator"
	"github.com/getsentry/sentry-go"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

func TestNewSentry_ShouldSucceedWithValidation(t *testing.T) {
	t.Parallel()

	t.Run("Should Succeed New Sentry", func(t *testing.T) {
		err := godotenv.Load("../.env.test")
		require.NoError(t, err)

		dsn := os.Getenv("SENTRY_DSN")
		sentry := commonSentry.NewSentry(
			validator.New(),
			commonSentry.WithDebug(true),
			commonSentry.WithDsn(dsn),
			commonSentry.WithEnv(constant.EnvProduction),
			commonSentry.WithSampleRate(1.0),
		)
		require.NotNil(t, sentry)
	})
}

func TestSetStartTransaction_ShouldSucceedWithTrace(t *testing.T) {
	t.Parallel()

	t.Run("Should Succeed SetStartTransaction", func(t *testing.T) {
		err := godotenv.Load("../.env.test")
		require.NoError(t, err)

		dsn := os.Getenv("SENTRY_DSN")
		newSentry := commonSentry.NewSentry(
			validator.New(),
			commonSentry.WithDebug(true),
			commonSentry.WithDsn(dsn),
			commonSentry.WithEnv(constant.EnvProduction),
			commonSentry.WithSampleRate(1.0),
		)
		require.NotNil(t, newSentry)

		key := constant.XRequestIdHeader
		val := uuid.NewString()
		ctx := context.WithValue(context.TODO(), key, val)
		newSentry.SetStartTransaction(
			ctx,
			"span name test",
			"transaction name",
			func(ctx context.Context, span *sentry.Span) (string, uint8) {
				newSentry.Trace(ctx, "span name test 2", func(ctx context.Context, span *sentry.Span) {
					log.Println("It works")
				})
				return "200", uint8(commonSentry.STATUS_OK)
			},
		)

	})
}

func TestSetStartTransaction_ShouldSucceedWithTraceManual(t *testing.T) {
	t.Parallel()

	t.Run("Should Succeed SetStartTransaction", func(t *testing.T) {
		err := godotenv.Load("../.env.test")
		require.NoError(t, err)

		dsn := os.Getenv("SENTRY_DSN")
		s := commonSentry.NewSentry(
			validator.New(),
			commonSentry.WithDebug(true),
			commonSentry.WithDsn(dsn),
			commonSentry.WithEnv(constant.EnvProduction),
			commonSentry.WithSampleRate(1.0),
		)
		require.NotNil(t, s)

		key := constant.XRequestIdHeader
		val := uuid.NewString()
		ctx := context.WithValue(context.TODO(), key, val)
		s.SetStartTransaction(
			ctx,
			"span name test",
			"transaction name",
			func(ctx context.Context, span *sentry.Span) (string, uint8) {
				sp := s.StartSpan(ctx, "span name test 2")
				defer s.Finish(sp)
				s.SetTag(span, "consumeGroupId", "123")
				return "200", uint8(commonSentry.STATUS_OK)
			},
		)

	})
}

func TestSetStartTransaction_ShouldSucceedWithTraceAndTag(t *testing.T) {
	t.Parallel()

	t.Run("Should Succeed SetStartTransaction", func(t *testing.T) {
		err := godotenv.Load("../.env.test")
		require.NoError(t, err)

		dsn := os.Getenv("SENTRY_DSN")
		newSentry := commonSentry.NewSentry(
			validator.New(),
			commonSentry.WithDebug(true),
			commonSentry.WithDsn(dsn),
			commonSentry.WithEnv(constant.EnvProduction),
			commonSentry.WithSampleRate(1.0),
		)
		require.NotNil(t, newSentry)

		key := constant.XRequestIdHeader
		val := uuid.NewString()
		ctx := context.WithValue(context.TODO(), key, val)
		newSentry.SetStartTransaction(
			ctx,
			"span name test",
			"transaction name",
			func(ctx context.Context, span *sentry.Span) (string, uint8) {
				newSentry.Trace(ctx, "span name test 2", func(ctx context.Context, span *sentry.Span) {
					newSentry.SetTag(span, "tag test", "tag value test")
				})
				return "200", uint8(commonSentry.STATUS_OK)
			},
		)

	})
}

func TestSetStartTransaction_ShouldSucceedWithTraceAndUserInfo(t *testing.T) {
	t.Parallel()

	t.Run("Should Succeed SetStartTransaction", func(t *testing.T) {
		err := godotenv.Load("../.env.test")
		require.NoError(t, err)

		dsn := os.Getenv("SENTRY_DSN")
		newSentry := commonSentry.NewSentry(
			validator.New(),
			commonSentry.WithDebug(true),
			commonSentry.WithDsn(dsn),
			commonSentry.WithEnv(constant.EnvProduction),
			commonSentry.WithSampleRate(1.0),
		)
		require.NotNil(t, newSentry)

		key := constant.XRequestIdHeader
		val := uuid.NewString()
		ctx := context.WithValue(context.TODO(), key, val)
		newSentry.SetStartTransaction(
			ctx,
			"span name test",
			"transaction name",
			func(ctx context.Context, span *sentry.Span) (string, uint8) {
				newSentry.Trace(ctx, "span name test 2", func(ctx context.Context, span *sentry.Span) {
					newSentry.SetUserInfo(commonSentry.UserInfoSentry{
						ID:       "1",
						Username: "John",
						Email:    "john@example.com",
					})
				})
				return "200", uint8(commonSentry.STATUS_OK)
			},
		)

	})
}

func TestSetStartTransaction_ShouldSucceedWithTraceAndCaptureException(t *testing.T) {
	t.Parallel()

	t.Run("Should Succeed SetStartTransaction", func(t *testing.T) {
		err := godotenv.Load("../.env.test")
		require.NoError(t, err)

		dsn := os.Getenv("SENTRY_DSN")
		newSentry := commonSentry.NewSentry(
			validator.New(),
			commonSentry.WithDebug(true),
			commonSentry.WithDsn(dsn),
			commonSentry.WithEnv(constant.EnvProduction),
			commonSentry.WithSampleRate(1.0),
		)
		require.NotNil(t, newSentry)

		key := constant.XRequestIdHeader
		val := uuid.NewString()
		ctx := context.WithValue(context.TODO(), key, val)
		newSentry.SetStartTransaction(
			ctx,
			"span name test",
			"transaction name",
			func(ctx context.Context, span *sentry.Span) (string, uint8) {
				newSentry.Trace(ctx, "span name test 2", func(ctx context.Context, span *sentry.Span) {
					newSentry.CaptureException(errors.New("error"))
				})
				return "200", uint8(commonSentry.STATUS_OK)
			},
		)

	})
}

func TestSetStartTransaction_ShouldSucceedWithTraceAndHandlingPanic(t *testing.T) {
	t.Parallel()

	t.Run("Should Succeed SetStartTransaction", func(t *testing.T) {
		err := godotenv.Load("../.env.test")
		require.NoError(t, err)

		dsn := os.Getenv("SENTRY_DSN")
		newSentry := commonSentry.NewSentry(
			validator.New(),
			commonSentry.WithDebug(true),
			commonSentry.WithDsn(dsn),
			commonSentry.WithEnv(constant.EnvProduction),
			commonSentry.WithSampleRate(1.0),
		)
		require.NotNil(t, newSentry)

		key := constant.XRequestIdHeader
		val := uuid.NewString()
		ctx := context.WithValue(context.TODO(), key, val)
		newSentry.SetStartTransaction(
			ctx,
			"span name test",
			"transaction name",
			func(ctx context.Context, span *sentry.Span) (string, uint8) {
				newSentry.Trace(ctx, "span name test 2", func(ctx context.Context, span *sentry.Span) {
					defer func() {
						if pnc := recover(); pnc != nil {
							newSentry.HandlingPanic(pnc)
						}
					}()
					panic("panic test")

				})
				return "200", uint8(commonSentry.STATUS_OK)
			},
		)

	})
}
