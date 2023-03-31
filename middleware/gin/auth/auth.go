package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strconv"
	"strings"
	"time"

	"bitbucket.org/moladinTech/go-lib-activity-log/model"
	moladinEvoClient "bitbucket.org/moladinTech/go-lib-common/client/moladin_evo"
	"bitbucket.org/moladinTech/go-lib-common/constant"
	"bitbucket.org/moladinTech/go-lib-common/middleware/gin/auth/rbac"
	responseModel "bitbucket.org/moladinTech/go-lib-common/response/model"
	"bitbucket.org/moladinTech/go-lib-common/sentry"
	"bitbucket.org/moladinTech/go-lib-common/signature"
	commonValidator "bitbucket.org/moladinTech/go-lib-common/validator"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/exp/slices"
)

type IMiddlewareAuth interface {
	AuthToken() gin.HandlerFunc
	AuthXApiKey() gin.HandlerFunc
	Auth() gin.HandlerFunc
	AuthSignature() gin.HandlerFunc
	AuthRoleRBAC(allowedRoles map[string]bool, applicationCode string) gin.HandlerFunc
	AuthPermissionRBAC(allowedPermissions map[string]bool, applicationCode string) gin.HandlerFunc
}

type MiddlewareAuthPackage struct {
	Sentry                  sentry.ISentry `validate:"required"`
	MoladinEvoClient        moladinEvoClient.IMoladinEvo
	Signature               signature.GenerateAndVerify
	RBACClient              rbac.IClientRBAC
	SignatureExpirationTime *uint
	ConfigApiKey            string
	PermittedRoles          []string
	ServiceName             string
	SecretKey               string
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

func WithSignatureExpirationTime(signatureExpirationTime *uint) Option {
	return func(s *MiddlewareAuthPackage) {
		s.SignatureExpirationTime = signatureExpirationTime
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

func WithRBACClient(rbacClient rbac.IClientRBAC) Option {
	return func(s *MiddlewareAuthPackage) {
		s.RBACClient = rbacClient
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
				c.AbortWithStatusJSON(http.StatusUnauthorized, responseModel.Response{
					Status: responseModel.StatusFail, Message: err.Error(),
				})
				return
			}

			if !slices.Contains(a.PermittedRoles, user.Role.Name) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, responseModel.Response{
					Status: responseModel.StatusFail, Message: http.StatusText(http.StatusUnauthorized),
				})
				return
			}

			c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), constant.XUserDetail, user))
			c.Next()
			return
		}

		c.AbortWithStatusJSON(http.StatusUnauthorized, responseModel.Response{
			Status: responseModel.StatusFail, Message: http.StatusText(http.StatusUnauthorized)})

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
			c.AbortWithStatusJSON(http.StatusUnauthorized, responseModel.Response{
				Status: responseModel.StatusFail, Message: http.StatusText(http.StatusUnauthorized)})
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
			c.JSON(http.StatusUnauthorized, responseModel.Response{
				Message: http.StatusText(http.StatusUnauthorized),
				Status:  responseModel.StatusFail,
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
				c.JSON(http.StatusUnauthorized, responseModel.Response{
					Message: http.StatusText(http.StatusUnauthorized),
					Status:  responseModel.StatusFail,
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
				c.JSON(http.StatusUnauthorized, responseModel.Response{
					Message: http.StatusText(http.StatusUnauthorized),
					Status:  responseModel.StatusFail,
				})
				c.Abort()
				return
			}
			c.Next()
			return
		}

		c.JSON(http.StatusUnauthorized, responseModel.Response{
			Message: http.StatusText(http.StatusUnauthorized),
			Status:  responseModel.StatusFail,
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
			c.AbortWithStatusJSON(http.StatusUnauthorized, responseModel.Response{
				Status: responseModel.StatusFail, Message: http.StatusText(http.StatusUnauthorized)})
			return
		}

		if a.SignatureExpirationTime != nil {
			requestAtUnix, _ := strconv.ParseInt(requestAt, 10, 64)
			requestAtUnix += int64((time.Hour * time.Duration(*a.SignatureExpirationTime)).Seconds())
			tnUnix := time.Now().Unix()
			if tnUnix > requestAtUnix {
				c.AbortWithStatusJSON(http.StatusUnauthorized, responseModel.Response{
					Status: responseModel.StatusFail, Message: http.StatusText(http.StatusUnauthorized)})
				return
			}
		}

		key := serviceNameSender + ":" + a.ServiceName + ":" + requestID + ":" + requestAt + ":" + a.SecretKey
		match := a.Signature.Verify(c.Request.Context(), key, requestSignature)
		if !match {
			c.AbortWithStatusJSON(http.StatusUnauthorized, responseModel.Response{
				Status: responseModel.StatusFail, Message: http.StatusText(http.StatusUnauthorized)})
			return
		}

		c.Next()
		return
	}
}

func (a *MiddlewareAuthPackage) AuthRoleRBAC(allowedRoles map[string]bool, applicationCode string) gin.HandlerFunc {
	return func(gc *gin.Context) {
		var token = gc.GetHeader(constant.AuthorizationHeader)
		if token == "" {
			gc.AbortWithStatusJSON(http.StatusUnauthorized, responseModel.Response{
				Status: responseModel.StatusFail, Message: http.StatusText(http.StatusUnauthorized)})
			return
		}

		tokenArr := strings.Split(token, " ")
		lastIndex := len(tokenArr) - 1
		token = tokenArr[lastIndex]
		isAllowed, userDetail, err := a.RBACClient.IsRoleAllowed(allowedRoles, token, applicationCode)
		if err != nil {
			gc.AbortWithStatusJSON(http.StatusUnauthorized, responseModel.Response{
				Status: responseModel.StatusFail, Message: err.Error(),
			})
			return
		}

		if !isAllowed {
			gc.AbortWithStatusJSON(http.StatusUnauthorized, responseModel.Response{
				Status: responseModel.StatusFail, Message: http.StatusText(http.StatusUnauthorized),
			})
			return
		}

		gc.Request = gc.Request.WithContext(context.WithValue(gc.Request.Context(), constant.XUserDetail, userDetail))
		gc.Next()
	}
}

func (a *MiddlewareAuthPackage) AuthPermissionRBAC(allowedPermissions map[string]bool, applicationCode string) gin.HandlerFunc {
	return func(gc *gin.Context) {
		var token = gc.GetHeader(constant.AuthorizationHeader)
		if token == "" {
			gc.AbortWithStatusJSON(http.StatusUnauthorized, responseModel.Response{
				Status: responseModel.StatusFail, Message: http.StatusText(http.StatusUnauthorized)})
			return
		}

		tokenArr := strings.Split(token, " ")
		lastIndex := len(tokenArr) - 1
		token = tokenArr[lastIndex]
		isAllowed, userDetail, err := a.RBACClient.IsPermissionAllowed(allowedPermissions, token, applicationCode)
		if err != nil {
			gc.AbortWithStatusJSON(http.StatusUnauthorized, responseModel.Response{
				Status: responseModel.StatusFail, Message: err.Error(),
			})
			return
		}

		if !isAllowed {
			gc.AbortWithStatusJSON(http.StatusUnauthorized, responseModel.Response{
				Status: responseModel.StatusFail, Message: http.StatusText(http.StatusUnauthorized),
			})
			return
		}

		gc.Request = gc.Request.WithContext(context.WithValue(gc.Request.Context(), constant.XUserDetail, userDetail))
		gc.Next()
	}
}
