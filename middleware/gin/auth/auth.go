package auth

import (
	"bitbucket.org/moladinTech/go-lib-activity-log/model"
	moladinEvoClient "bitbucket.org/moladinTech/go-lib-common/client/moladin_evo"
	"bitbucket.org/moladinTech/go-lib-common/constant"
	"bitbucket.org/moladinTech/go-lib-common/response"
	"bitbucket.org/moladinTech/go-lib-common/sentry"
	"bitbucket.org/moladinTech/go-lib-common/signature"
	commonValidator "bitbucket.org/moladinTech/go-lib-common/validator"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/exp/slices"
	"net/http"
)

type IMiddlewareAuth interface {
	AuthToken() gin.HandlerFunc
	AuthXApiKey() gin.HandlerFunc
	Auth() gin.HandlerFunc
	AuthSignature() gin.HandlerFunc
}

type MiddlewareAuthPackage struct {
	Sentry           sentry.ISentry `validate:"required"`
	MoladinEvoClient moladinEvoClient.IMoladinEvo
	Signature        signature.GenerateAndVerify
	ConfigApiKey     string
	PermittedRoles   []string
	ServiceName      string
	SecretKey        string
}

func WithSentry(sentry sentry.ISentry) Option {
	return func(s *MiddlewareAuthPackage) {
		s.Sentry = sentry
	}
}

func WithMoladinEvoClient(moladinEvoClient moladinEvoClient.IMoladinEvo) Option {
	return func(s *MiddlewareAuthPackage) {
		s.MoladinEvoClient = moladinEvoClient
	}
}

func WithSignature(signature signature.GenerateAndVerify) Option {
	return func(s *MiddlewareAuthPackage) {
		s.Signature = signature
	}
}

func WithConfigApiKey(configApiKey string) Option {
	return func(s *MiddlewareAuthPackage) {
		s.ConfigApiKey = configApiKey
	}
}

func WithPermittedRoles(permittedRoles []string) Option {
	return func(s *MiddlewareAuthPackage) {
		s.PermittedRoles = permittedRoles
	}
}

func WithServiceName(serviceName string) Option {
	return func(s *MiddlewareAuthPackage) {
		s.ServiceName = serviceName
	}
}

func WithSecretKey(secretKey string) Option {
	return func(s *MiddlewareAuthPackage) {
		s.SecretKey = secretKey
	}
}

type Option func(*MiddlewareAuthPackage)

func NewAuth(
	validator *validator.Validate,
	options ...Option,
) IMiddlewareAuth {
	middlewareAuthPackage := &MiddlewareAuthPackage{}

	for _, option := range options {
		option(middlewareAuthPackage)
	}

	err := validator.Struct(middlewareAuthPackage)
	if err != nil {
		panic(commonValidator.ToErrResponse(err))
	}

	return middlewareAuthPackage
}

func (a *MiddlewareAuthPackage) AuthToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		const logCtx = "common.middleware.gin.auth.AuthToken"
		reqCtx := c.Request.Context()
		span := a.Sentry.StartSpan(reqCtx, logCtx)
		defer span.Finish()

		var token = c.GetHeader(constant.AuthorizationHeader)
		if token != "" {
			user, err := a.MoladinEvoClient.UserDetail(c.Request.Context(), token)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, response.Response{Status: response.StatusFail, Message: err.Error()})
				return
			}

			if !slices.Contains(a.PermittedRoles, user.Role.Name) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, response.Response{Status: response.StatusFail, Message: http.StatusText(http.StatusUnauthorized)})
				return
			}

			c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), constant.XUserDetail, user))
			c.Next()
			return
		}

		c.AbortWithStatusJSON(http.StatusUnauthorized, response.Response{Status: response.StatusFail, Message: http.StatusText(http.StatusUnauthorized)})

	}
}

func (a *MiddlewareAuthPackage) AuthXApiKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		const logCtx = "common.middleware.gin.auth.AuthXApiKey"
		reqCtx := c.Request.Context()
		span := a.Sentry.StartSpan(reqCtx, logCtx)
		defer span.Finish()

		var apiKey = c.GetHeader(constant.XApiKeyHeader)
		if apiKey != a.ConfigApiKey {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.Response{Status: response.StatusFail, Message: http.StatusText(http.StatusUnauthorized)})
			return
		}

		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), constant.XUserDetail, model.UserDetail{
			UserId: 0,
			Name:   "SYSTEM",
			Email:  "system@moladin.com",
		}))
		c.Next()
	}
}

func (a *MiddlewareAuthPackage) Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizationToken := c.GetHeader(constant.AuthorizationHeader)
		xApiKey := c.GetHeader(constant.XApiKeyHeader)
		xServiceName := c.GetHeader(constant.XServiceNameHeader)
		if xApiKey == "" && xServiceName == "" && authorizationToken == "" {
			c.JSON(http.StatusUnauthorized, response.Response{
				Message: http.StatusText(http.StatusUnauthorized),
				Status:  response.StatusFail,
			})
			c.Abort()
			return
		}

		if authorizationToken != "" {
			var (
				token     = c.GetHeader(constant.AuthorizationHeader)
				user, err = a.MoladinEvoClient.UserDetail(c.Request.Context(), token)
			)
			if err != nil {
				c.JSON(http.StatusUnauthorized, response.Response{
					Message: http.StatusText(http.StatusUnauthorized),
					Status:  response.StatusFail,
				})
				c.Abort()
				return
			}
			c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), constant.XUserDetail, user))
			c.Next()
			return
		}

		if xApiKey != "" && xServiceName != "" {
			token := []byte(xServiceName + xServiceName)
			validateKey := sha256.Sum256(token)
			if xApiKey != hex.EncodeToString(validateKey[:]) {
				c.JSON(http.StatusUnauthorized, response.Response{
					Message: http.StatusText(http.StatusUnauthorized),
					Status:  response.StatusFail,
				})
				c.Abort()
				return
			}
			c.Next()
			return
		}

		c.JSON(http.StatusUnauthorized, response.Response{
			Message: http.StatusText(http.StatusUnauthorized),
			Status:  response.StatusFail,
		})
		c.Abort()
	}
}

func (a *MiddlewareAuthPackage) AuthSignature() gin.HandlerFunc {
	return func(c *gin.Context) {
		const logCtx = "common.middleware.gin.auth.AuthSignature"
		reqCtx := c.Request.Context()
		span := a.Sentry.StartSpan(reqCtx, logCtx)
		defer span.Finish()

		var (
			serviceNameSender = c.GetHeader(constant.XServiceNameHeader)
			requestSignature  = c.GetHeader(constant.XRequestSignatureHeader)
			requestID         = c.GetHeader(constant.XRequestIdHeader)
			requestAt         = c.GetHeader(constant.XRequestAtHeader)
		)

		if serviceNameSender == "" || requestSignature == "" || requestID == "" || requestAt == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.Response{Status: response.StatusFail, Message: http.StatusText(http.StatusUnauthorized)})
			return
		}

		key := serviceNameSender + ":" + a.ServiceName + ":" + requestID + ":" + requestAt + ":" + a.SecretKey
		match := a.Signature.Verify(c.Request.Context(), key, requestSignature)
		if !match {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.Response{Status: response.StatusFail, Message: http.StatusText(http.StatusUnauthorized)})
			return
		}

		c.Next()
		return
	}
}
