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
	"sec-ch-ua",
	"sec-ch-ua-mobile",
	"sec-ch-ua-platform",
	"Content-Type",
	"content-type",
	"Content-Length",
	"content-length",
	"Accept",
	"accept",
	"Origin",
	"origin",
	"Referer",
	"referer",
	"User-Agent",
	"user-agent",
	"Accept-Encoding",
	"accept-encoding",
	"X-CSRF-Token",
	"x-csrf-token",
	"Authorization",
	"authorization",
	"Cache-Control",
	"cache-control",
	"X-Requested-With",
	"x-requested-with",
	"X-Request-Id",
	"x-request-id",
	"X-Origin-Path",
	"x-origin-path",
	"x-Service-Name",
	"x-service-name",
	"x-Api-Key",
	"x-api-key",
	"X-Menu-Slug",
	"x-menu-slug",
	"x-menu-test-additional",
}

func TestCors_ShouldSucceed(t *testing.T) {
	t.Run("Should Succeed CORS", func(t *testing.T) {

		mockResponse := ``

		gin.SetMode(gin.TestMode)
		r := gin.Default()
		r.Routes()
		r.GET("/", middlewareGin.CORS("x-menu-test-additional"))

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
