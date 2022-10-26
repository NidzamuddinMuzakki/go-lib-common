package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net/http"

	"bitbucket.org/moladinTech/go-lib-activity-log/model"
	moladinEvoClient "bitbucket.org/moladinTech/go-lib-common/client/moladin_evo"
	"bitbucket.org/moladinTech/go-lib-common/constant"
	"bitbucket.org/moladinTech/go-lib-common/response"
	"bitbucket.org/moladinTech/go-lib-common/sentry"
	commonValidator "bitbucket.org/moladinTech/go-lib-common/validator"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/exp/slices"
)

type IMiddlewareAuth interface {
	AuthToken() gin.HandlerFunc
	AuthXApiKey() gin.HandlerFunc
	Auth() gin.HandlerFunc
}

type MiddlewareAuthPackage struct {
	Sentry           sentry.ISentry `validate:"required"`
	MoladinEvoClient moladinEvoClient.IMoladinEvo
	ConfigApiKey     string
	PermittedRoles   []string
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