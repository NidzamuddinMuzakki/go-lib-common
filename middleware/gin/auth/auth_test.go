package auth_test

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"bitbucket.org/moladinTech/go-lib-activity-log/model"
	"bitbucket.org/moladinTech/go-lib-common/cast"
	moladinEvoMock "bitbucket.org/moladinTech/go-lib-common/client/moladin_evo/mocks"
	"bitbucket.org/moladinTech/go-lib-common/constant"
	"bitbucket.org/moladinTech/go-lib-common/middleware/gin/auth"
	sentryMock "bitbucket.org/moladinTech/go-lib-common/sentry/mocks"
	"bitbucket.org/moladinTech/go-lib-common/signature"
	signatureMock "bitbucket.org/moladinTech/go-lib-common/signature/mocks"
	"bitbucket.org/moladinTech/go-lib-common/validator"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewAuth_ShouldSucceedWithValidation(t *testing.T) {
	t.Run("Should Succeed New Auth", func(t *testing.T) {
		xApiKey := "dummy"
		sentry := sentryMock.NewISentry(t)
		moladinEvo := moladinEvoMock.NewIMoladinEvo(t)
		auth := auth.NewAuth(validator.New(),
			auth.WithSentry(sentry),
			auth.WithConfigApiKey(xApiKey),
			auth.WithPermittedRoles([]string{"finance"}),
			auth.WithMoladinEvoClient(moladinEvo),
		)
		require.NotNil(t, auth)
	})
}

func TestNewAuth_ErrorOnValidation(t *testing.T) {
	t.Run("Error On Validation New Auth", func(t *testing.T) {
		xApiKey := "dummy"
		moladinEvo := moladinEvoMock.NewIMoladinEvo(t)
		require.Panics(t, func() {
			auth.NewAuth(validator.New(),
				auth.WithConfigApiKey(xApiKey),
				auth.WithPermittedRoles([]string{"finance"}),
				auth.WithMoladinEvoClient(moladinEvo),
			)
		})
	})
}

// ////////// Func AuthToken
func TestAuthToken_ErrorOnEmptyHeader(t *testing.T) {
	t.Run("Error on empty header", func(t *testing.T) {
		span := sentry.Span{}
		sentry := sentryMock.NewISentry(t)
		sentry.On("StartSpan", mock.Anything, mock.Anything).
			Return(&span).
			Once()

		xApiKey := "dummy"
		moladinEvo := moladinEvoMock.NewIMoladinEvo(t)
		auth := auth.NewAuth(validator.New(),
			auth.WithSentry(sentry),
			auth.WithConfigApiKey(xApiKey),
			auth.WithPermittedRoles([]string{"finance"}),
			auth.WithMoladinEvoClient(moladinEvo),
		)
		require.NotNil(t, auth)

		mockResponse := `{"status":"fail","message":"Unauthorized"}`

		gin.SetMode(gin.TestMode)
		r := gin.Default()
		r.GET("/", auth.AuthToken())

		req, _ := http.NewRequest("GET", "/", nil)
		req.Header = http.Header{}
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		responseData, _ := ioutil.ReadAll(w.Body)
		require.Equal(t, mockResponse, string(responseData))
		require.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestAuthToken_ShouldSucceedOnRoleAuthorized(t *testing.T) {
	t.Run("Should succeed On Role Authorized", func(t *testing.T) {
		span := sentry.Span{}
		token := "Bearer token"
		sentry := sentryMock.NewISentry(t)
		sentry.On("StartSpan", mock.Anything, mock.Anything).
			Return(&span).
			Once()

		xApiKey := "dummy"
		moladinEvo := moladinEvoMock.NewIMoladinEvo(t)
		moladinEvo.On("UserDetail", mock.Anything, token).
			Return(model.UserDetail{
				UserId: 1,
				Name:   "John",
				Email:  "john@example.com",
				Role: model.UserRole{
					Id:   1,
					Name: "finance",
				},
			}, nil).
			Once()
		auth := auth.NewAuth(validator.New(),
			auth.WithSentry(sentry),
			auth.WithConfigApiKey(xApiKey),
			auth.WithPermittedRoles([]string{"finance"}),
			auth.WithMoladinEvoClient(moladinEvo),
		)
		require.NotNil(t, auth)

		mockResponse := ``

		gin.SetMode(gin.TestMode)
		r := gin.Default()
		r.GET("/", auth.AuthToken())

		req, _ := http.NewRequest("GET", "/", nil)
		req.Header = http.Header{}
		req.Header.Set(constant.AuthorizationHeader, token)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		responseData, _ := ioutil.ReadAll(w.Body)
		require.Equal(t, mockResponse, string(responseData))
		require.Equal(t, http.StatusOK, w.Code)
	})
}

func TestAuthToken_ErrorOnUnauthorizedRole(t *testing.T) {
	t.Run("Error on Unauthorized Role", func(t *testing.T) {
		span := sentry.Span{}
		token := "Bearer token"
		sentry := sentryMock.NewISentry(t)
		sentry.On("StartSpan", mock.Anything, mock.Anything).
			Return(&span).
			Once()

		xApiKey := "dummy"
		moladinEvo := moladinEvoMock.NewIMoladinEvo(t)
		moladinEvo.On("UserDetail", mock.Anything, token).
			Return(model.UserDetail{
				UserId: 1,
				Name:   "John",
				Email:  "john@example.com",
				Role: model.UserRole{
					Id:   1,
					Name: "admin",
				},
			}, nil).
			Once()
		auth := auth.NewAuth(validator.New(),
			auth.WithSentry(sentry),
			auth.WithConfigApiKey(xApiKey),
			auth.WithPermittedRoles([]string{"finance"}),
			auth.WithMoladinEvoClient(moladinEvo),
		)
		require.NotNil(t, auth)

		mockResponse := `{"status":"fail","message":"Unauthorized"}`

		gin.SetMode(gin.TestMode)
		r := gin.Default()
		r.GET("/", auth.AuthToken())

		req, _ := http.NewRequest("GET", "/", nil)
		req.Header = http.Header{}
		req.Header.Set(constant.AuthorizationHeader, token)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		responseData, _ := ioutil.ReadAll(w.Body)
		require.Equal(t, mockResponse, string(responseData))
		require.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestAuthToken_ErrorOnUserDetail(t *testing.T) {
	t.Run("Error on Unauthorized Role", func(t *testing.T) {
		span := sentry.Span{}
		token := "Bearer token"
		sentry := sentryMock.NewISentry(t)
		sentry.On("StartSpan", mock.Anything, mock.Anything).
			Return(&span).
			Once()

		xApiKey := "dummy"
		moladinEvo := moladinEvoMock.NewIMoladinEvo(t)
		moladinEvo.On("UserDetail", mock.Anything, token).
			Return(model.UserDetail{}, errors.New("Unauthorized")).
			Once()
		auth := auth.NewAuth(validator.New(),
			auth.WithSentry(sentry),
			auth.WithConfigApiKey(xApiKey),
			auth.WithPermittedRoles([]string{"finance"}),
			auth.WithMoladinEvoClient(moladinEvo),
		)
		require.NotNil(t, auth)

		mockResponse := `{"status":"fail","message":"Unauthorized"}`

		gin.SetMode(gin.TestMode)
		r := gin.Default()
		r.GET("/", auth.AuthToken())

		req, _ := http.NewRequest("GET", "/", nil)
		req.Header = http.Header{}
		req.Header.Set(constant.AuthorizationHeader, token)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		responseData, _ := ioutil.ReadAll(w.Body)
		require.Equal(t, mockResponse, string(responseData))
		require.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

// ////////// Func AuthXApiKey
func TestAuthXApiKey_ErrorOnEmptyHeader(t *testing.T) {
	t.Run("Error on empty header", func(t *testing.T) {
		span := sentry.Span{}
		sentry := sentryMock.NewISentry(t)
		sentry.On("StartSpan", mock.Anything, mock.Anything).
			Return(&span).
			Once()

		xApiKey := "dummy"
		moladinEvo := moladinEvoMock.NewIMoladinEvo(t)
		auth := auth.NewAuth(validator.New(),
			auth.WithSentry(sentry),
			auth.WithConfigApiKey(xApiKey),
			auth.WithPermittedRoles([]string{"finance"}),
			auth.WithMoladinEvoClient(moladinEvo),
		)
		require.NotNil(t, auth)

		mockResponse := `{"status":"fail","message":"Unauthorized"}`

		gin.SetMode(gin.TestMode)
		r := gin.Default()
		r.GET("/", auth.AuthXApiKey())

		req, _ := http.NewRequest("GET", "/", nil)
		req.Header = http.Header{}
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		responseData, _ := ioutil.ReadAll(w.Body)
		require.Equal(t, mockResponse, string(responseData))
		require.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestAuthXApiKey_ShouldSucceedOnMatchingApiKey(t *testing.T) {
	t.Run("Should succeed on matching api key", func(t *testing.T) {
		span := sentry.Span{}
		sentry := sentryMock.NewISentry(t)
		sentry.On("StartSpan", mock.Anything, mock.Anything).
			Return(&span).
			Once()

		xApiKey := "dummy"
		moladinEvo := moladinEvoMock.NewIMoladinEvo(t)
		auth := auth.NewAuth(validator.New(),
			auth.WithSentry(sentry),
			auth.WithConfigApiKey(xApiKey),
			auth.WithPermittedRoles([]string{"finance"}),
			auth.WithMoladinEvoClient(moladinEvo),
		)
		require.NotNil(t, auth)

		mockResponse := ``

		gin.SetMode(gin.TestMode)
		r := gin.Default()
		r.GET("/", auth.AuthXApiKey())

		req, _ := http.NewRequest("GET", "/", nil)
		req.Header = http.Header{}
		req.Header.Set(constant.XApiKeyHeader, xApiKey)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		responseData, _ := ioutil.ReadAll(w.Body)
		require.Equal(t, mockResponse, string(responseData))
		require.Equal(t, http.StatusOK, w.Code)
	})
}

// Funtion AuthSignature
func TestAuthSignature_ShouldSucceedOnMatchingSignature(t *testing.T) {
	t.Run("Should succeed on matching signature", func(t *testing.T) {
		span := sentry.Span{}
		sentry := sentryMock.NewISentry(t)
		sentry.On("StartSpan", mock.Anything, mock.Anything).
			Return(&span).
			Once()

		key := "service-sender:service-received:b3590f68-3d79-4ae5-8dfc-fd08b10c4e14:1670984752:secretbanget"
		s, err := signature.NewSignature(
			signature.WithAlgorithm(signature.BCrypt),
		)
		if err != nil {
			panic(err)
		}

		ctx := context.TODO()
		hashed, err := s.Generate(ctx,
			key)
		if err != nil {
			panic(err)
		}

		auth := auth.NewAuth(validator.New(),
			auth.WithSentry(sentry),
			auth.WithSignature(s),
			auth.WithServiceName("service-received"),
			auth.WithSecretKey("secretbanget"),
		)
		require.NotNil(t, auth)

		mockResponse := ``

		gin.SetMode(gin.TestMode)
		r := gin.Default()
		r.GET("/", auth.AuthSignature())

		req, _ := http.NewRequest("GET", "/", nil)
		req.Header = http.Header{}
		req.Header.Set(constant.XRequestSignatureHeader, hashed)
		req.Header.Set(constant.XServiceNameHeader, "service-sender")
		req.Header.Set(constant.XRequestIdHeader, "b3590f68-3d79-4ae5-8dfc-fd08b10c4e14")
		req.Header.Set(constant.XRequestAtHeader, "1670984752")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		responseData, _ := ioutil.ReadAll(w.Body)
		require.Equal(t, mockResponse, string(responseData))
		require.Equal(t, http.StatusOK, w.Code)
	})
}

func TestAuthSignature_ErrorOnExpirationValidation(t *testing.T) {
	t.Run("Error on expiration validation", func(t *testing.T) {
		span := sentry.Span{}
		sentry := sentryMock.NewISentry(t)
		sentry.On("StartSpan", mock.Anything, mock.Anything).
			Return(&span).
			Once()

		key := "service-sender:service-received:b3590f68-3d79-4ae5-8dfc-fd08b10c4e14:1676285674:secretbanget"
		s, err := signature.NewSignature(
			signature.WithAlgorithm(signature.BCrypt),
		)
		if err != nil {
			panic(err)
		}

		ctx := context.TODO()
		hashed, err := s.Generate(ctx,
			key)
		if err != nil {
			panic(err)
		}

		auth := auth.NewAuth(validator.New(),
			auth.WithSentry(sentry),
			auth.WithSignature(s),
			auth.WithSignatureExpirationTime(cast.NewPointer[uint](1)),
			auth.WithServiceName("service-received"),
			auth.WithSecretKey("secretbanget"),
		)
		require.NotNil(t, auth)

		mockResponse := `{"status":"fail","message":"Unauthorized"}`

		gin.SetMode(gin.TestMode)
		r := gin.Default()
		r.GET("/", auth.AuthSignature())

		req, _ := http.NewRequest("GET", "/", nil)
		req.Header = http.Header{}
		req.Header.Set(constant.XRequestSignatureHeader, hashed)
		req.Header.Set(constant.XServiceNameHeader, "service-sender")
		req.Header.Set(constant.XRequestIdHeader, "b3590f68-3d79-4ae5-8dfc-fd08b10c4e14")
		req.Header.Set(constant.XRequestAtHeader, "1676285674")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		responseData, _ := ioutil.ReadAll(w.Body)
		require.Equal(t, mockResponse, string(responseData))
		require.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestAuthSignature_ErrorOnEmptyHeader(t *testing.T) {
	t.Run("Error on empty header", func(t *testing.T) {
		span := sentry.Span{}
		sentry := sentryMock.NewISentry(t)
		sentry.On("StartSpan", mock.Anything, mock.Anything).
			Return(&span).
			Once()

		s, err := signature.NewSignature(
			signature.WithAlgorithm(signature.BCrypt),
		)
		if err != nil {
			panic(err)
		}

		auth := auth.NewAuth(validator.New(),
			auth.WithSentry(sentry),
			auth.WithSignature(s),
			auth.WithServiceName("service-received"),
			auth.WithSecretKey("secretbanget"),
		)
		require.NotNil(t, auth)

		mockResponse := `{"status":"fail","message":"Unauthorized"}`

		gin.SetMode(gin.TestMode)
		r := gin.Default()
		r.GET("/", auth.AuthSignature())

		req, _ := http.NewRequest("GET", "/", nil)
		req.Header = http.Header{}
		req.Header.Set(constant.XServiceNameHeader, "service-sender")
		req.Header.Set(constant.XRequestIdHeader, "b3590f68-3d79-4ae5-8dfc-fd08b10c4e14")
		req.Header.Set(constant.XRequestAtHeader, "1670984752")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		responseData, _ := ioutil.ReadAll(w.Body)
		require.Equal(t, mockResponse, string(responseData))
		require.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestAuthSignature_ErrorOnKeyAndSignatureNotMatch(t *testing.T) {
	t.Run("Error on key and signature not match", func(t *testing.T) {
		span := sentry.Span{}
		sentry := sentryMock.NewISentry(t)
		sentry.On("StartSpan", mock.Anything, mock.Anything).
			Return(&span).
			Once()

		sign := signatureMock.NewGenerateAndVerify(t)
		sign.On("Verify", mock.Anything, mock.Anything, mock.Anything).
			Return(false).
			Once()

		auth := auth.NewAuth(validator.New(),
			auth.WithSentry(sentry),
			auth.WithSignature(sign),
			auth.WithServiceName("service-received"),
			auth.WithSecretKey("secretbanget"),
		)
		require.NotNil(t, auth)

		mockResponse := `{"status":"fail","message":"Unauthorized"}`

		gin.SetMode(gin.TestMode)
		r := gin.Default()
		r.GET("/", auth.AuthSignature())

		req, _ := http.NewRequest("GET", "/", nil)
		req.Header = http.Header{}
		req.Header.Set(constant.XRequestSignatureHeader, "dummy")
		req.Header.Set(constant.XServiceNameHeader, "service-sender")
		req.Header.Set(constant.XRequestIdHeader, "b3590f68-3d79-4ae5-8dfc-fd08b10c4e14")
		req.Header.Set(constant.XRequestAtHeader, "1670984752")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		responseData, _ := ioutil.ReadAll(w.Body)
		require.Equal(t, mockResponse, string(responseData))
		require.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
