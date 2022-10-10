package gin

import (
	"context"
	"net/http"

	"bitbucket.org/moladinTech/go-lib-activity-log/model"
	moladinEvoClient "bitbucket.org/moladinTech/go-lib-common/client/moladin_evo"
	"bitbucket.org/moladinTech/go-lib-common/constant"
	"bitbucket.org/moladinTech/go-lib-common/response"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slices"
)

type MiddlewareAuth interface {
	AuthToken() gin.HandlerFunc
	AuthXApiKey() gin.HandlerFunc
}

type middlewareAuth struct {
	moladinEvoClient moladinEvoClient.MoladinEvo
	configApiKey     string
	permittedRoles   []string
}

func NewAuth(moladinEvoClient moladinEvoClient.MoladinEvo, configApiKey string, permittedRoles []string) MiddlewareAuth {
	return &middlewareAuth{
		moladinEvoClient: moladinEvoClient,
		configApiKey:     configApiKey,
		permittedRoles:   permittedRoles,
	}
}

func (a *middlewareAuth) AuthToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		var token = c.GetHeader("authorization")
		if token != "" {
			user, err := a.moladinEvoClient.UserDetail(c.Request.Context(), token)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, response.Response{Status: response.StatusFail, Message: err.Error()})
				return
			}

			if !slices.Contains(a.permittedRoles, user.Role.Name) {
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

func (a *middlewareAuth) AuthXApiKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		var apiKey = c.GetHeader("x-api-key")
		if apiKey != a.configApiKey {
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
