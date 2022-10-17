package errors

import (
	"context"
	"net/http"

	"bitbucket.org/moladinTech/go-lib-common/logger"
	"bitbucket.org/moladinTech/go-lib-common/response"
	"github.com/gin-gonic/gin"
)

// helper to logger the error and send response
func HttpErrResp(ctx context.Context, err error, c *gin.Context, args ...interface{}) {
	logger.Error(ctx, `error`, err)
	if val, ok := MapErrorResponse[GetErrKey(err)]; ok {
		c.JSON(val.StatusCode, val.Response)
	} else { // default return message when error happen
		c.JSON(http.StatusInternalServerError, response.Response{
			Message: http.StatusText(http.StatusInternalServerError),
			Status:  response.StatusFail,
		})
	}
}
