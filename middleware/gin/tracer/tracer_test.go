package tracer_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"bitbucket.org/moladinTech/go-lib-common/middleware/gin/tracer"
	sentryMock "bitbucket.org/moladinTech/go-lib-common/sentry/mocks"
	"bitbucket.org/moladinTech/go-lib-common/validator"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewTracer_ShouldSucceedWithValidation(t *testing.T) {
	t.Run("Should Succeed New Tracer", func(t *testing.T) {
		sentry := sentryMock.NewISentry(t)
		tracer := tracer.NewTracer(validator.New(),
			tracer.WithSentry(sentry),
		)
		require.NotNil(t, tracer)
	})
}

func TestTracer_ShouldSucceed(t *testing.T) {
	t.Run("Should Succeed Tracer", func(t *testing.T) {
		sentry := sentryMock.NewISentry(t)
		sentry.On("SetStartTransaction",
			mock.Anything,
			"middleware.Tracer",
			"POST /",
			mock.Anything,
		).Once()

		tracer := tracer.NewTracer(validator.New(),
			tracer.WithSentry(sentry),
		)
		require.NotNil(t, tracer)

		mockResponse := ``

		gin.SetMode(gin.TestMode)
		r := gin.Default()
		r.POST("/", tracer.Tracer())

		payload := map[string]string{
			"ID": "data ID",
		}
		jsonValue, _ := json.Marshal(payload)

		req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(jsonValue))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		responseData, _ := ioutil.ReadAll(w.Body)
		require.Equal(t, mockResponse, string(responseData))
		require.Equal(t, http.StatusOK, w.Code)
	})
}

func TestNewTracer_ErrorOnValidation(t *testing.T) {
	t.Run("Error On Validation New Tracer", func(t *testing.T) {
		require.Panics(t, func() {
			tracer.NewTracer(validator.New())
		})
	})
}
