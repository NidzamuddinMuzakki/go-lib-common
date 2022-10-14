package gin_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	middlewareGin "bitbucket.org/moladinTech/go-lib-common/middleware/gin"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

var allowHeaderRules = []string{
	"Content-Type",
	"Content-Length",
	"Accept-Encoding",
	"X-CSRF-Token",
	"Authorization",
	"accept",
	"origin",
	"Cache-Control",
	"X-Requested-With",
	"X-Request-Id",
	"X-Origin-Path",
	"x-service-name",
	"x-api-key",
}

func TestCors_ShouldSucceed(t *testing.T) {
	t.Run("Should Succeed CORS", func(t *testing.T) {

		mockResponse := ``

		gin.SetMode(gin.TestMode)
		r := gin.Default()
		r.Routes()
		r.GET("/", middlewareGin.CORS())

		req, _ := http.NewRequest("GET", "/", nil)
		req.Header = http.Header{}
		for key, allowHeaderRule := range allowHeaderRules {
			req.Header.Set(allowHeaderRule, strconv.Itoa(key))
		}
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		resultHeaders := strings.Split(w.Result().Header.Get("Access-Control-Allow-Headers"), ",")

		responseData, _ := ioutil.ReadAll(w.Body)
		require.Equal(t, mockResponse, string(responseData))
		require.Equal(t, http.StatusOK, w.Code)

		var mapAllowHeaderRule = map[string]struct{}{}
		var missingResultHeader = []string{}
		for _, allowHeaderRule := range allowHeaderRules {
			key := strings.ToLower(allowHeaderRule)
			mapAllowHeaderRule[key] = struct{}{}
		}
		for _, resAllowHeader := range resultHeaders {
			if _, ok := mapAllowHeaderRule[strings.ToLower(resAllowHeader)]; !ok {
				missingResultHeader = append(missingResultHeader, resAllowHeader)
			}
		}
		require.Empty(t, missingResultHeader)

		var mapResultHeader = map[string]struct{}{}
		var missingAllowHeaderRules = []string{}
		for _, resAllowHeader := range resultHeaders {
			key := strings.ToLower(resAllowHeader)
			mapResultHeader[key] = struct{}{}
		}
		for _, allowHeaderRule := range allowHeaderRules {
			if _, ok := mapResultHeader[strings.ToLower(allowHeaderRule)]; !ok {
				missingAllowHeaderRules = append(missingAllowHeaderRules, allowHeaderRule)
			}
		}
		require.Empty(t, missingAllowHeaderRules)

	})
}

func TestCors_ErrorOnMethoOptions(t *testing.T) {
	t.Run("Error on Method Options", func(t *testing.T) {

		mockResponse := ``

		gin.SetMode(gin.TestMode)
		r := gin.Default()
		r.Routes()
		r.OPTIONS("/", middlewareGin.CORS())

		req, _ := http.NewRequest("OPTIONS", "/", nil)
		req.Header = http.Header{}
		for key, allowHeaderRule := range allowHeaderRules {
			req.Header.Set(allowHeaderRule, strconv.Itoa(key))
		}
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		responseData, _ := ioutil.ReadAll(w.Body)
		require.Equal(t, mockResponse, string(responseData))
		require.Equal(t, http.StatusNoContent, w.Code)

	})
}
