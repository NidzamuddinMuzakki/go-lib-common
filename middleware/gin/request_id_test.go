package gin_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"bitbucket.org/moladinTech/go-lib-common/constant"
	middlewareGin "bitbucket.org/moladinTech/go-lib-common/middleware/gin"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestRequestId_ShouldSucceedWithHeader(t *testing.T) {
	t.Run("Should Succeed Request ID", func(t *testing.T) {

		mockResponse := ``

		gin.SetMode(gin.TestMode)
		r := gin.Default()
		r.Routes()
		r.GET("/", middlewareGin.RequestID())

		req, _ := http.NewRequest("GET", "/", nil)
		req.Header = http.Header{}
		req.Header.Set(constant.XRequestIdHeader, uuid.NewString())
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		responseData, _ := ioutil.ReadAll(w.Body)

		require.Equal(t, mockResponse, string(responseData))
		require.Equal(t, http.StatusOK, w.Code)
	})
}

func TestRequestId_ShouldSucceedWithoutHeader(t *testing.T) {
	t.Run("Should Succeed Request ID With Header", func(t *testing.T) {

		mockResponse := ``

		gin.SetMode(gin.TestMode)
		r := gin.Default()
		r.Routes()
		r.GET("/", middlewareGin.RequestID())

		req, _ := http.NewRequest("GET", "/", nil)
		req.Header = http.Header{}
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		responseData, _ := ioutil.ReadAll(w.Body)

		require.Equal(t, mockResponse, string(responseData))
		require.Equal(t, http.StatusOK, w.Code)
	})
}
