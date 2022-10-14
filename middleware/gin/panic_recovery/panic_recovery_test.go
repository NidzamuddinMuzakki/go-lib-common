package panic_recovery_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

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
		sentry := sentryMock.NewISentry(t)
		panicRecovery := panic_recovery.NewPanicRecovery(validator.New(),
			panic_recovery.WithConfigEnv(dummy),
			panic_recovery.WithSentry(sentry),
		)
		require.NotNil(t, panicRecovery)
	})
}

func TestPanicRecoveryMiddleware_ShouldSucceedWithEnvironmentProduction(t *testing.T) {
	t.Run("Should succeeed with environment production", func(t *testing.T) {
		panicText := "panic on middleware"
		span := sentry.Span{}
		sentry := sentryMock.NewISentry(t)
		sentry.On("HandlingPanic", panicText).
			Once()
		sentry.On("StartSpan", mock.Anything, mock.Anything).
			Return(&span).
			Once()

		panicRecovery := panic_recovery.NewPanicRecovery(validator.New(),
			panic_recovery.WithConfigEnv(constant.EnvProduction),
			panic_recovery.WithSentry(sentry),
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

		responseData, _ := ioutil.ReadAll(w.Body)
		require.Equal(t, mockResponse, string(responseData))
		require.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestPanicRecoveryMiddleware_ShouldSucceedWithEnvironmentStaging(t *testing.T) {
	t.Run("Should succeeed with environment staging", func(t *testing.T) {
		panicText := "panic on middleware"
		span := sentry.Span{}
		sentry := sentryMock.NewISentry(t)
		sentry.On("HandlingPanic", panicText).
			Once()
		sentry.On("StartSpan", mock.Anything, mock.Anything).
			Return(&span).
			Once()

		panicRecovery := panic_recovery.NewPanicRecovery(validator.New(),
			panic_recovery.WithConfigEnv(constant.EnvStaging),
			panic_recovery.WithSentry(sentry),
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

		responseData, _ := ioutil.ReadAll(w.Body)
		require.Equal(t, mockResponse, string(responseData))
		require.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestPanicRecoveryMiddleware_ShouldSucceedWithEnvironmentDevelopment(t *testing.T) {
	t.Run("Should succeeed with environment development", func(t *testing.T) {
		panicText := "panic on middleware"
		span := sentry.Span{}
		sentry := sentryMock.NewISentry(t)
		sentry.On("HandlingPanic", panicText).
			Once()
		sentry.On("StartSpan", mock.Anything, mock.Anything).
			Return(&span).
			Once()

		panicRecovery := panic_recovery.NewPanicRecovery(validator.New(),
			panic_recovery.WithConfigEnv(constant.EnvDevelopment),
			panic_recovery.WithSentry(sentry),
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

		responseData, _ := ioutil.ReadAll(w.Body)
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
