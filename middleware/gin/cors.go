package gin

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func CORS(additionalHeader ...string) gin.HandlerFunc {
	allowHeaders := []string{
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
	}

	// append additional header from config
	allowHeaders = append(allowHeaders, additionalHeader...)

	allowMethods := []string{
		"POST",
		"GET",
		"PUT",
		"PATCH",
		"DELETE",
		"OPTIONS",
	}

	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", strings.Join(allowHeaders, ","))
		c.Writer.Header().Set("Access-Control-Allow-Methods", strings.Join(allowMethods, ","))

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
