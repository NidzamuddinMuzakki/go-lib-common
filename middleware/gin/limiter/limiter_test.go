package limiter

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"bitbucket.org/moladinTech/go-lib-common/cast"
	mocks "bitbucket.org/moladinTech/go-lib-common/mocks/cache"
	"bitbucket.org/moladinTech/go-lib-common/validator"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestShouldBeAbleToGetRequestWithDefaultValue(t *testing.T) {
	red := &redis.IntCmd{}
	red.SetVal(1)
	ttl := &redis.DurationCmd{}
	ttl.SetVal(-1)

	mockCacher := mocks.NewCacher(t)
	mockCacher.On("Incr", context.TODO(), mock.Anything).Return(red, nil)
	mockCacher.On("Ttl", context.TODO(), mock.Anything).Return(ttl, nil)
	mockCacher.On("Expire", context.TODO(), mock.Anything, mock.Anything).Return(nil, nil)

	_, limiter := NewLimiter(
		validator.New(),
		WithCacher(mockCacher),
		WithServiceName("test-service"),
	)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/", limiter, func(c *gin.Context) {
		c.JSON(http.StatusOK, map[string]interface{}{
			"status":  "success",
			"message": "success",
		})
	})

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header = http.Header{}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	mockResponse := `{"message":"success","status":"success"}`

	responseData, _ := io.ReadAll(w.Body)
	require.Equal(t, mockResponse, string(responseData))
	require.Equal(t, http.StatusOK, w.Code)

}

func TestShouldBeAbleToGetRequestWithCustomValue(t *testing.T) {
	red := &redis.IntCmd{}
	red.SetVal(1)
	ttl := &redis.DurationCmd{}
	ttl.SetVal(-1)

	mockCacher := mocks.NewCacher(t)
	mockCacher.On("Incr", context.TODO(), mock.Anything).Return(red, nil)
	mockCacher.On("Ttl", context.TODO(), mock.Anything).Return(ttl, nil)
	mockCacher.On("Expire", context.TODO(), mock.Anything, mock.Anything).Return(nil, nil)

	mlp, _ := NewLimiter(
		validator.New(),
		WithCacher(mockCacher),
		WithServiceName("test-service"),
	)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/", mlp.WithCustomLimit(RateLimit{Limit: cast.NewPointer[uint](10), TTL: cast.NewPointer[uint](1)}), func(c *gin.Context) {
		c.JSON(http.StatusOK, map[string]interface{}{
			"status":  "success",
			"message": "success",
		})
	})

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header = http.Header{}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	mockResponse := `{"message":"success","status":"success"}`

	responseData, _ := io.ReadAll(w.Body)
	require.Equal(t, mockResponse, string(responseData))
	require.Equal(t, http.StatusOK, w.Code)

}

func TestShouldBeUnableToGetRequest(t *testing.T) {
	red := &redis.IntCmd{}
	red.SetVal(1)
	ttl := &redis.DurationCmd{}
	ttl.SetVal(-1)

	mockCacher := mocks.NewCacher(t)
	mockCacher.On("Incr", context.TODO(), mock.Anything).Return(red, nil)
	mockCacher.On("Ttl", context.TODO(), mock.Anything).Return(ttl, nil)
	mockCacher.On("Expire", context.TODO(), mock.Anything, mock.Anything).Return(nil, nil)

	mlp, _ := NewLimiter(
		validator.New(),
		WithCacher(mockCacher),
		WithServiceName("test-service"),
	)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/", mlp.WithCustomLimit(RateLimit{Limit: cast.NewPointer[uint](0), TTL: cast.NewPointer[uint](30)}))

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header = http.Header{}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	mockResponse := `{"status":"fail","message":"Too Many Requests"}`

	responseData, _ := io.ReadAll(w.Body)
	require.Equal(t, mockResponse, string(responseData))
	require.Equal(t, http.StatusTooManyRequests, w.Code)
}
