package panic_recovery_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	slackMock "bitbucket.org/moladinTech/go-lib-common/client/notification/slack/mocks"
	"bitbucket.org/moladinTech/go-lib-common/constant"
	"bitbucket.org/moladinTech/go-lib-common/middleware/gin/panic_recovery"
	sentryMock "bitbucket.org/moladinTech/go-lib-common/sentry/mocks"
	"bitbucket.org/moladinTech/go-lib-common/validator"
	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewPanicRecovery_ShouldSucceedWithValidation(t *testing.T) {
	t.Run("Should Succeed New Panic Recovery", func(t *testing.T) {
		dummy := "dummy"
		sentryClient := sentryMock.NewISentry(t)
		slack := slackMock.NewISlack(t)
		panicRecovery := panic_recovery.NewPanicRecovery(validator.New(),
			panic_recovery.WithConfigEnv(dummy),
			panic_recovery.WithSentry(sentryClient),
			panic_recovery.WithSlack(slack),
		)
		require.NotNil(t, panicRecovery)
	})
}

func TestNewPanicRecovery_ShouldSucceedWithoutSlack(t *testing.T) {
	t.Run("Should Succeed New Panic Recovery Without Slack", func(t *testing.T) {
		dummy := "dummy"
		sentryClient := sentryMock.NewISentry(t)
		panicRecovery := panic_recovery.NewPanicRecovery(validator.New(),
			panic_recovery.WithConfigEnv(dummy),
			panic_recovery.WithSentry(sentryClient),
		)
		require.NotNil(t, panicRecovery)
	})
}

func TestPanicRecoveryMiddleware_ShouldSucceedWithEnvironmentProduction(t *testing.T) {
	t.Run("Should succeeed with environment production", func(t *testing.T) {
		panicText := "panic on middleware"
		slackMessage := "dummy slack message"
		dummyEventID := sentry.NewEvent().EventID
		span := sentry.Span{}
		sentryClient := sentryMock.NewISentry(t)
		slack := slackMock.NewISlack(t)
		sentryClient.On("HandlingPanic", panicText).
			Once()
		sentryClient.On("StartSpan", mock.Anything, mock.Anything).
			Return(&span).
			Twice()
		sentryClient.On("SpanContext", mock.Anything).
			Return(context.Background()).
			Twice()
		sentryClient.On("Finish", mock.Anything).
			Twice()
		sentryClient.On("CaptureException", errors.New(fmt.Sprintf("panic: %v", panicText))).
			Return(&dummyEventID).
			Once()
		slack.On("GetFormattedMessage", mock.Anything, mock.Anything, mock.Anything).
			Return(slackMessage).
			Once()
		slack.On("Send", mock.Anything, mock.Anything).
			Return(nil).
			Once()

		panicRecovery := panic_recovery.NewPanicRecovery(validator.New(),
			panic_recovery.WithConfigEnv(constant.EnvProduction),
			panic_recovery.WithSentry(sentryClient),
			panic_recovery.WithSlack(slack),
		)
		require.NotNil(t, panicRecovery)

		mockResponse := `{"status":"fail","message":"Server error. Contact admin for more information."}`

		gin.SetMode(gin.TestMode)
		r := gin.Default()
		r.GET("/", panicRecovery.PanicRecoveryMiddleware(), func(ctx *gin.Context) {
			panic(panicText)
		})
		req, _ := http.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		responseData, _ := io.ReadAll(w.Body)
		require.Equal(t, mockResponse, string(responseData))
		require.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestPanicRecoveryMiddleware_ShouldSucceedWithEnvironmentProductionWithoutSlack(t *testing.T) {
	t.Run("Should succeeed with environment production without slack", func(t *testing.T) {
		panicText := "panic on middleware"
		dummyEventID := sentry.NewEvent().EventID
		span := sentry.Span{}
		sentryClient := sentryMock.NewISentry(t)
		sentryClient.On("HandlingPanic", panicText).
			Once()
		sentryClient.On("StartSpan", mock.Anything, mock.Anything).
			Return(&span).
			Once()
		sentryClient.On("SpanContext", mock.Anything).
			Return(context.Background()).
			Once()
		sentryClient.On("Finish", mock.Anything).
			Once()
		sentryClient.On("CaptureException", errors.New(fmt.Sprintf("panic: %v", panicText))).
			Return(&dummyEventID).
			Once()

		panicRecovery := panic_recovery.NewPanicRecovery(validator.New(),
			panic_recovery.WithConfigEnv(constant.EnvProduction),
			panic_recovery.WithSentry(sentryClient),
		)
		require.NotNil(t, panicRecovery)

		mockResponse := `{"status":"fail","message":"Server error. Contact admin for more information."}`

		gin.SetMode(gin.TestMode)
		r := gin.Default()
		r.GET("/", panicRecovery.PanicRecoveryMiddleware(), func(ctx *gin.Context) {
			panic(panicText)
		})
		req, _ := http.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		responseData, _ := io.ReadAll(w.Body)
		require.Equal(t, mockResponse, string(responseData))
		require.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestPanicRecoveryMiddleware_ShouldSucceedWithEnvironmentStaging(t *testing.T) {
	t.Run("Should succeeed with environment staging", func(t *testing.T) {
		panicText := "panic on middleware"
		slackMessage := "dummy slack message"
		dummyEventID := sentry.NewEvent().EventID
		span := sentry.Span{}
		sentryClient := sentryMock.NewISentry(t)
		slack := slackMock.NewISlack(t)
		sentryClient.On("HandlingPanic", panicText).
			Once()
		sentryClient.On("StartSpan", mock.Anything, mock.Anything).
			Return(&span).
			Twice()
		sentryClient.On("SpanContext", mock.Anything).
			Return(context.Background()).
			Twice()
		sentryClient.On("Finish", mock.Anything).
			Twice()
		sentryClient.On("CaptureException", errors.New(fmt.Sprintf("panic: %v", panicText))).
			Return(&dummyEventID).
			Once()
		slack.On("GetFormattedMessage", mock.Anything, mock.Anything, mock.Anything).
			Return(slackMessage).
			Once()
		slack.On("Send", mock.Anything, mock.Anything).
			Return(nil).
			Once()

		panicRecovery := panic_recovery.NewPanicRecovery(validator.New(),
			panic_recovery.WithConfigEnv(constant.EnvStaging),
			panic_recovery.WithSentry(sentryClient),
			panic_recovery.WithSlack(slack),
		)
		require.NotNil(t, panicRecovery)

		mockResponse := `{"status":"fail","message":"panic: ` + panicText + `"}`

		gin.SetMode(gin.TestMode)
		r := gin.Default()
		r.GET("/", panicRecovery.PanicRecoveryMiddleware(), func(ctx *gin.Context) {
			panic(panicText)
		})
		req, _ := http.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		responseData, _ := io.ReadAll(w.Body)
		require.Equal(t, mockResponse, string(responseData))
		require.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestPanicRecoveryMiddleware_ShouldSucceedWithEnvironmentStagingWithoutSlack(t *testing.T) {
	t.Run("Should succeeed with environment staging wihtout slack", func(t *testing.T) {
		panicText := "panic on middleware"
		dummyEventID := sentry.NewEvent().EventID
		span := sentry.Span{}
		sentryClient := sentryMock.NewISentry(t)
		sentryClient.On("HandlingPanic", panicText).
			Once()
		sentryClient.On("StartSpan", mock.Anything, mock.Anything).
			Return(&span).
			Once()
		sentryClient.On("SpanContext", mock.Anything).
			Return(context.Background()).
			Once()
		sentryClient.On("Finish", mock.Anything).
			Once()
		sentryClient.On("CaptureException", errors.New(fmt.Sprintf("panic: %v", panicText))).
			Return(&dummyEventID).
			Once()

		panicRecovery := panic_recovery.NewPanicRecovery(validator.New(),
			panic_recovery.WithConfigEnv(constant.EnvStaging),
			panic_recovery.WithSentry(sentryClient),
		)
		require.NotNil(t, panicRecovery)

		mockResponse := `{"status":"fail","message":"panic: ` + panicText + `"}`

		gin.SetMode(gin.TestMode)
		r := gin.Default()
		r.GET("/", panicRecovery.PanicRecoveryMiddleware(), func(ctx *gin.Context) {
			panic(panicText)
		})
		req, _ := http.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		responseData, _ := io.ReadAll(w.Body)
		require.Equal(t, mockResponse, string(responseData))
		require.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestPanicRecoveryMiddleware_ShouldSucceedWithEnvironmentDevelopment(t *testing.T) {
	t.Run("Should succeeed with environment development", func(t *testing.T) {
		panicText := "panic on middleware"
		slackMessage := "dummy slack message"
		dummyEventID := sentry.NewEvent().EventID
		span := sentry.Span{}
		sentryClient := sentryMock.NewISentry(t)
		slack := slackMock.NewISlack(t)
		sentryClient.On("HandlingPanic", panicText).
			Once()
		sentryClient.On("StartSpan", mock.Anything, mock.Anything).
			Return(&span).
			Twice()
		sentryClient.On("SpanContext", mock.Anything).
			Return(context.Background()).
			Twice()
		sentryClient.On("Finish", mock.Anything).
			Twice()
		sentryClient.On("CaptureException", errors.New(fmt.Sprintf("panic: %v", panicText))).
			Return(&dummyEventID).
			Once()
		slack.On("GetFormattedMessage", mock.Anything, mock.Anything, mock.Anything).
			Return(slackMessage).
			Once()
		slack.On("Send", mock.Anything, mock.Anything).
			Return(nil).
			Once()

		panicRecovery := panic_recovery.NewPanicRecovery(validator.New(),
			panic_recovery.WithConfigEnv(constant.EnvDevelopment),
			panic_recovery.WithSentry(sentryClient),
			panic_recovery.WithSlack(slack),
		)
		require.NotNil(t, panicRecovery)

		mockResponse := `{"status":"fail","message":"panic: ` + panicText + `"}`

		gin.SetMode(gin.TestMode)
		r := gin.Default()
		r.GET("/", panicRecovery.PanicRecoveryMiddleware(), func(ctx *gin.Context) {
			panic(panicText)
		})
		req, _ := http.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		responseData, _ := io.ReadAll(w.Body)
		require.Equal(t, mockResponse, string(responseData))
		require.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestPanicRecoveryMiddleware_ShouldSucceedWithEnvironmentDevelopmentWithoutSlack(t *testing.T) {
	t.Run("Should succeeed with environment development without slack", func(t *testing.T) {
		panicText := "panic on middleware"
		dummyEventID := sentry.NewEvent().EventID
		span := sentry.Span{}
		sentryClient := sentryMock.NewISentry(t)
		sentryClient.On("HandlingPanic", panicText).
			Once()
		sentryClient.On("StartSpan", mock.Anything, mock.Anything).
			Return(&span).
			Once()
		sentryClient.On("SpanContext", mock.Anything).
			Return(context.Background()).
			Once()
		sentryClient.On("Finish", mock.Anything).
			Once()
		sentryClient.On("CaptureException", errors.New(fmt.Sprintf("panic: %v", panicText))).
			Return(&dummyEventID).
			Once()

		panicRecovery := panic_recovery.NewPanicRecovery(validator.New(),
			panic_recovery.WithConfigEnv(constant.EnvDevelopment),
			panic_recovery.WithSentry(sentryClient),
		)
		require.NotNil(t, panicRecovery)

		mockResponse := `{"status":"fail","message":"panic: ` + panicText + `"}`

		gin.SetMode(gin.TestMode)
		r := gin.Default()
		r.GET("/", panicRecovery.PanicRecoveryMiddleware(), func(ctx *gin.Context) {
			panic(panicText)
		})
		req, _ := http.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		responseData, _ := io.ReadAll(w.Body)
		require.Equal(t, mockResponse, string(responseData))
		require.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestPanicRecoveryMiddleware_ShouldErrorSendSlack(t *testing.T) {
	t.Run("Should error send slack", func(t *testing.T) {
		panicText := "panic on middleware"
		slackMessage := "dummy slack message"
		dummyEventID := sentry.NewEvent().EventID
		span := sentry.Span{}
		sentryClient := sentryMock.NewISentry(t)
		slack := slackMock.NewISlack(t)
		sentryClient.On("HandlingPanic", panicText).
			Once()
		sentryClient.On("StartSpan", mock.Anything, mock.Anything).
			Return(&span).
			Twice()
		sentryClient.On("SpanContext", mock.Anything).
			Return(context.Background()).
			Twice()
		sentryClient.On("Finish", mock.Anything).
			Twice()
		sentryClient.On("CaptureException", errors.New(fmt.Sprintf("panic: %v", panicText))).
			Return(&dummyEventID).
			Once()
		sentryClient.On("CaptureException", errors.New("error slack")).
			Return(&dummyEventID).
			Once()
		slack.On("GetFormattedMessage", mock.Anything, mock.Anything, mock.Anything).
			Return(slackMessage).
			Once()
		slack.On("Send", mock.Anything, mock.Anything).
			Return(errors.New("error slack")).
			Once()

		panicRecovery := panic_recovery.NewPanicRecovery(validator.New(),
			panic_recovery.WithConfigEnv(constant.EnvDevelopment),
			panic_recovery.WithSentry(sentryClient),
			panic_recovery.WithSlack(slack),
		)
		require.NotNil(t, panicRecovery)

		mockResponse := `{"status":"fail","message":"panic: ` + panicText + `"}`

		gin.SetMode(gin.TestMode)
		r := gin.Default()
		r.GET("/", panicRecovery.PanicRecoveryMiddleware(), func(ctx *gin.Context) {
			panic(panicText)
		})
		req, _ := http.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		responseData, _ := io.ReadAll(w.Body)
		require.Equal(t, mockResponse, string(responseData))
		require.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestNewPanicRecovery_ErrorOnValidation(t *testing.T) {
	t.Run("Error On Validation New Panic Recovery", func(t *testing.T) {
		dummy := "dummy"
		require.Panics(t, func() {
			panic_recovery.NewPanicRecovery(validator.New(),
				panic_recovery.WithConfigEnv(dummy),
			)
		})
	})
}
