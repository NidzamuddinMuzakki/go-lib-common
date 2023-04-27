package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"bitbucket.org/moladinTech/go-lib-activity-log/model"
	moladinEvoClient "bitbucket.org/moladinTech/go-lib-common/client/moladin_evo"
	"bitbucket.org/moladinTech/go-lib-common/constant"
	"bitbucket.org/moladinTech/go-lib-common/logger"
	"bitbucket.org/moladinTech/go-lib-common/middleware/gin/auth/rbac"
	responseModel "bitbucket.org/moladinTech/go-lib-common/response/model"
	"bitbucket.org/moladinTech/go-lib-common/sentry"
	"bitbucket.org/moladinTech/go-lib-common/signature"
	commonValidator "bitbucket.org/moladinTech/go-lib-common/validator"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"golang.org/x/exp/slices"
)

type XUserType string

const (
	XSTRBAC      XUserType = "RBAC"
	XSTEvo       XUserType = "EVO"
	XSTSignature XUserType = "Signature"
	XSTAPIKey    XUserType = "APIKey"
)

type IMiddlewareAuth interface {
	Auth(rbacPermissions []string) gin.HandlerFunc
	AuthToken() gin.HandlerFunc
	AuthXApiKey() gin.HandlerFunc
	AuthSignature() gin.HandlerFunc
	AuthRoleRBAC(allowedRoles []string) gin.HandlerFunc
	AuthPermissionRBAC(allowedPermissions []string) gin.HandlerFunc
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

func (a *MiddlewareAuthPackage) Auth(rbacPermissions []string) gin.HandlerFunc {
	return func(gc *gin.Context) {
		const logCtx = "common.middleware.gin.auth.Auth"

		reqCtx := gc.Request.Context()
		span := a.Sentry.StartSpan(reqCtx, logCtx)
		defer span.Finish()

		authorizationToken := gc.GetHeader(constant.AuthorizationHeader)
		xApiKey := gc.GetHeader(constant.XApiKeyHeader)
		xServiceName := gc.GetHeader(constant.XServiceNameHeader)
		xSignature := gc.GetHeader(constant.XRequestSignatureHeader)

		if xApiKey == "" && xServiceName == "" && authorizationToken == "" && xSignature == "" {
			gc.AbortWithStatusJSON(http.StatusUnauthorized, responseModel.Response{
				Message: http.StatusText(http.StatusUnauthorized),
				Status:  responseModel.StatusFail,
			})
			return
		}

		if authorizationToken != "" {
			if a.RBACClient != nil {
				token := strings.Split(authorizationToken, " ")
				if len(token) < 2 {
					gc.AbortWithStatusJSON(http.StatusUnauthorized, responseModel.Response{
						Message: http.StatusText(http.StatusUnauthorized),
						Status:  responseModel.StatusFail,
					})
					return
				}

				_, err := uuid.Parse(token[1])
				if err == nil {
					isPermitted, user, errCheck := a.RBACClient.IsPermissionAllowed(rbacPermissions, token[1])
					if errCheck != nil {
						logger.Error(gc.Request.Context(), err.Error(), err, logger.Tag{
							Key:   "logCtx",
							Value: logCtx,
						})
					}

					if errCheck == nil && isPermitted {
						gc.Request = gc.Request.WithContext(context.WithValue(gc.Request.Context(), constant.AuthorizationHeader, authorizationToken))
						gc.Request = gc.Request.WithContext(context.WithValue(gc.Request.Context(), constant.XUserType, XSTRBAC))
						gc.Request = gc.Request.WithContext(context.WithValue(gc.Request.Context(), constant.XUserDetail, user))
						gc.Next()
						return
					}
				}
			}

			if a.MoladinEvoClient != nil {
				user, err := a.MoladinEvoClient.UserDetail(gc.Request.Context(), authorizationToken)
				if err != nil {
					logger.Error(gc.Request.Context(), err.Error(), err, logger.Tag{
						Key:   "logCtx",
						Value: logCtx,
					})
					gc.AbortWithStatusJSON(http.StatusUnauthorized, responseModel.Response{
						Status:  responseModel.StatusFail,
						Message: http.StatusText(http.StatusUnauthorized),
					})
					return
				}

				if slices.Contains(a.PermittedRoles, user.Role.Name) {
					gc.Request = gc.Request.WithContext(context.WithValue(gc.Request.Context(), constant.AuthorizationHeader, authorizationToken))
					gc.Request = gc.Request.WithContext(context.WithValue(gc.Request.Context(), constant.XUserType, XSTEvo))
					gc.Request = gc.Request.WithContext(context.WithValue(gc.Request.Context(), constant.XUserDetail, user))
					gc.Next()
					return
				}
			}
		}

		if xSignature != "" && xServiceName != "" {
			isPermitted := a.Signature.Verify(gc.Request.Context(), a.SecretKey, xSignature)
			if isPermitted {
				gc.Request = gc.Request.WithContext(context.WithValue(gc.Request.Context(), constant.AuthorizationHeader, xSignature))
				gc.Request = gc.Request.WithContext(context.WithValue(gc.Request.Context(), constant.XUserType, XSTSignature))
				gc.Request = gc.Request.WithContext(context.WithValue(gc.Request.Context(), constant.XUserDetail, model.UserDetail{
					UserId: 0,
					Name:   xServiceName,
					Email:  fmt.Sprintf("%s@moladin.com", xServiceName),
				}))
				gc.Next()
				return
			}
		}

		if xApiKey != "" && xServiceName != "" {
			token := []byte(xServiceName + xServiceName)
			validateKey := sha256.Sum256(token)
			if xApiKey == hex.EncodeToString(validateKey[:]) {
				gc.Request = gc.Request.WithContext(context.WithValue(gc.Request.Context(), constant.AuthorizationHeader, xApiKey))
				gc.Request = gc.Request.WithContext(context.WithValue(gc.Request.Context(), constant.XUserType, XSTAPIKey))
				gc.Request = gc.Request.WithContext(context.WithValue(gc.Request.Context(), constant.XUserDetail, model.UserDetail{
					UserId: 0,
					Name:   xServiceName,
					Email:  fmt.Sprintf("%s@moladin.com", xServiceName),
				}))
				gc.Next()
				return
			}
		}

		gc.AbortWithStatusJSON(http.StatusUnauthorized, responseModel.Response{
			Message: http.StatusText(http.StatusUnauthorized),
			Status:  responseModel.StatusFail,
		})
		return
	}
}

func (a *MiddlewareAuthPackage) AuthToken() gin.HandlerFunc {
	return func(gc *gin.Context) {
		const logCtx = "common.middleware.gin.auth.AuthToken"

		reqCtx := gc.Request.Context()
		span := a.Sentry.StartSpan(reqCtx, logCtx)
		defer span.Finish()

		var token = gc.GetHeader(constant.AuthorizationHeader)
		if token != "" {
			user, err := a.MoladinEvoClient.UserDetail(gc.Request.Context(), token)
			if err != nil {
				logger.Error(gc.Request.Context(), err.Error(), err, logger.Tag{
					Key:   "logCtx",
					Value: logCtx,
				})
				gc.AbortWithStatusJSON(http.StatusUnauthorized, responseModel.Response{
					Status:  responseModel.StatusFail,
					Message: http.StatusText(http.StatusUnauthorized),
				})
				return
			}

			if !slices.Contains(a.PermittedRoles, user.Role.Name) {
				gc.AbortWithStatusJSON(http.StatusUnauthorized, responseModel.Response{
					Status:  responseModel.StatusFail,
					Message: http.StatusText(http.StatusUnauthorized),
				})
				return
			}

			gc.Request = gc.Request.WithContext(context.WithValue(gc.Request.Context(), constant.XUserType, XSTEvo))
			gc.Request = gc.Request.WithContext(context.WithValue(gc.Request.Context(), constant.XUserDetail, user))
			gc.Next()
			return
		}

		gc.AbortWithStatusJSON(http.StatusUnauthorized, responseModel.Response{
			Status:  responseModel.StatusFail,
			Message: http.StatusText(http.StatusUnauthorized)},
		)
	}
}

func (a *MiddlewareAuthPackage) AuthXApiKey() gin.HandlerFunc {
	return func(gc *gin.Context) {
		const logCtx = "common.middleware.gin.auth.AuthXApiKey"

		reqCtx := gc.Request.Context()
		span := a.Sentry.StartSpan(reqCtx, logCtx)
		defer span.Finish()

		var apiKey = gc.GetHeader(constant.XApiKeyHeader)
		if apiKey != a.ConfigApiKey {
			gc.AbortWithStatusJSON(http.StatusUnauthorized, responseModel.Response{
				Status:  responseModel.StatusFail,
				Message: http.StatusText(http.StatusUnauthorized)},
			)
			return
		}

		serviceNameSender := gc.GetHeader(constant.XServiceNameHeader)
		gc.Request = gc.Request.WithContext(context.WithValue(gc.Request.Context(), constant.XUserType, XSTAPIKey))
		gc.Request = gc.Request.WithContext(context.WithValue(gc.Request.Context(), constant.XUserDetail, model.UserDetail{
			UserId: 0,
			Name:   serviceNameSender,
			Email:  fmt.Sprintf("%s@moladin.com", serviceNameSender),
		}))
		gc.Next()
	}
}

func (a *MiddlewareAuthPackage) AuthSignature() gin.HandlerFunc {
	return func(gc *gin.Context) {
		const logCtx = "common.middleware.gin.auth.AuthSignature"

		reqCtx := gc.Request.Context()
		span := a.Sentry.StartSpan(reqCtx, logCtx)
		defer span.Finish()

		var (
			serviceNameSender = gc.GetHeader(constant.XServiceNameHeader)
			requestSignature  = gc.GetHeader(constant.XRequestSignatureHeader)
			requestID         = gc.GetHeader(constant.XRequestIdHeader)
			requestAt         = gc.GetHeader(constant.XRequestAtHeader)
		)

		if serviceNameSender == "" || requestSignature == "" || requestID == "" || requestAt == "" {
			gc.AbortWithStatusJSON(http.StatusUnauthorized, responseModel.Response{
				Status:  responseModel.StatusFail,
				Message: http.StatusText(http.StatusUnauthorized)},
			)
			return
		}

		if a.SignatureExpirationTime != nil {
			requestAtUnix, _ := strconv.ParseInt(requestAt, 10, 64)
			requestAtUnix += int64((time.Hour * time.Duration(*a.SignatureExpirationTime)).Seconds())
			tnUnix := time.Now().Unix()
			if tnUnix > requestAtUnix {
				gc.AbortWithStatusJSON(http.StatusUnauthorized, responseModel.Response{
					Status:  responseModel.StatusFail,
					Message: http.StatusText(http.StatusUnauthorized)},
				)
				return
			}
		}

		key := serviceNameSender + ":" + a.ServiceName + ":" + requestID + ":" + requestAt + ":" + a.SecretKey
		match := a.Signature.Verify(gc.Request.Context(), key, requestSignature)
		if !match {
			gc.AbortWithStatusJSON(http.StatusUnauthorized, responseModel.Response{
				Status:  responseModel.StatusFail,
				Message: http.StatusText(http.StatusUnauthorized)},
			)
			return
		}

		gc.Request = gc.Request.WithContext(context.WithValue(gc.Request.Context(), constant.XUserType, XSTSignature))
		gc.Request = gc.Request.WithContext(context.WithValue(gc.Request.Context(), constant.XUserDetail, model.UserDetail{
			UserId: 0,
			Name:   serviceNameSender,
			Email:  fmt.Sprintf("%s@moladin.com", serviceNameSender),
		}))
		gc.Next()
	}
}

func (a *MiddlewareAuthPackage) AuthRoleRBAC(allowedRoles []string) gin.HandlerFunc {
	return func(gc *gin.Context) {
		const logCtx = "common.middleware.gin.auth.AuthRoleRBAC"

		reqCtx := gc.Request.Context()
		span := a.Sentry.StartSpan(reqCtx, logCtx)
		defer span.Finish()

		var token = gc.GetHeader(constant.AuthorizationHeader)
		if token == "" {
			gc.AbortWithStatusJSON(http.StatusUnauthorized, responseModel.Response{
				Status:  responseModel.StatusFail,
				Message: http.StatusText(http.StatusUnauthorized)},
			)
			return
		}

		tokenArr := strings.Split(token, " ")
		lastIndex := len(tokenArr) - 1
		token = tokenArr[lastIndex]
		isAllowed, userDetail, err := a.RBACClient.IsRoleAllowed(allowedRoles, token)
		if err != nil {
			logger.Error(gc.Request.Context(), err.Error(), err, logger.Tag{
				Key:   "logCtx",
				Value: logCtx,
			})
			gc.AbortWithStatusJSON(http.StatusUnauthorized, responseModel.Response{
				Status:  responseModel.StatusFail,
				Message: http.StatusText(http.StatusUnauthorized),
			})
			return
		}

		if !isAllowed {
			gc.AbortWithStatusJSON(http.StatusUnauthorized, responseModel.Response{
				Status:  responseModel.StatusFail,
				Message: http.StatusText(http.StatusUnauthorized)},
			)
			return
		}

		gc.Request = gc.Request.WithContext(context.WithValue(gc.Request.Context(), constant.XUserType, XSTRBAC))
		gc.Request = gc.Request.WithContext(context.WithValue(gc.Request.Context(), constant.XUserDetail, userDetail))
		gc.Next()
	}
}

func (a *MiddlewareAuthPackage) AuthPermissionRBAC(allowedPermissions []string) gin.HandlerFunc {
	return func(gc *gin.Context) {
		const logCtx = "common.middleware.gin.auth.AuthPermissionRBAC"

		reqCtx := gc.Request.Context()
		span := a.Sentry.StartSpan(reqCtx, logCtx)
		defer span.Finish()

		var token = gc.GetHeader(constant.AuthorizationHeader)
		if token == "" {
			gc.AbortWithStatusJSON(http.StatusUnauthorized, responseModel.Response{
				Status:  responseModel.StatusFail,
				Message: http.StatusText(http.StatusUnauthorized)},
			)
			return
		}

		tokenArr := strings.Split(token, " ")
		lastIndex := len(tokenArr) - 1
		token = tokenArr[lastIndex]
		isAllowed, userDetail, err := a.RBACClient.IsPermissionAllowed(allowedPermissions, token)
		if err != nil || !isAllowed {
			if err != nil {
				logger.Error(gc.Request.Context(), err.Error(), err, logger.Tag{
					Key:   "logCtx",
					Value: logCtx,
				})
			}
			gc.AbortWithStatusJSON(http.StatusUnauthorized, responseModel.Response{
				Status:  responseModel.StatusFail,
				Message: http.StatusText(http.StatusUnauthorized),
			})
			return
		}

		gc.Request = gc.Request.WithContext(context.WithValue(gc.Request.Context(), constant.XUserType, XSTRBAC))
		gc.Request = gc.Request.WithContext(context.WithValue(gc.Request.Context(), constant.XUserDetail, userDetail))
		gc.Next()
	}
}
