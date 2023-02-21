package logger_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"bitbucket.org/moladinTech/go-lib-common/constant"
	commonLogger "bitbucket.org/moladinTech/go-lib-common/logger"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func Test_NewContextFromParent(t *testing.T) {
	t.Parallel()
	httpRecorder := httptest.NewRecorder()
	ginCtx, _ := gin.CreateTestContext(httpRecorder)
	ginCtx.Request = &http.Request{
		Method: "get",
		URL: &url.URL{
			Host: "localhost",
		},
	}
	newCtx := context.WithValue(ginCtx.Request.Context(), constant.XRequestIdHeader, "1234")
	ginCtx.Request = ginCtx.Request.WithContext(newCtx)

	assert.NotNil(t, ginCtx.Request.Context())

	http.NewRequest("get", "localhost", nil)
	ctx := commonLogger.NewContextFromParent(ginCtx)

	expectedCtx := commonLogger.AddRequestID(
		context.WithValue(
			context.Background(),
			constant.XRequestIdHeader,
			"1234",
		),
		"1234")

	assert.Equal(t, expectedCtx, ctx)
}
