package errors

import (
	"context"
	"net/http"

	"bitbucket.org/moladinTech/go-lib-common/logger"
	"bitbucket.org/moladinTech/go-lib-common/registry"
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

// HttpErrRespAndSendNotif is helper to logger the error and send response and send notification (if statusCode >= 500)
func HttpErrRespAndSendNotif(ctx context.Context, err error, c *gin.Context, commonRegistry registry.IRegistry) {
	logger.Error(ctx, `error`, err)
	if val, ok := MapErrorResponse[GetErrKey(err)]; ok {
		if val.StatusCode >= 500 {
			commonRegistry.GetSentry().CaptureException(err)
			// send notif
			logCtx := getLogCtx(err)
			slackMessage := commonRegistry.GetSlack().GetFormattedMessage(logCtx, ctx, err)
			errSlack := commonRegistry.GetSlack().Send(ctx, slackMessage)
			if errSlack != nil {
				logger.Error(ctx, "Error sending notif to slack", err)
			}
		}

		// send response
		c.JSON(val.StatusCode, val.Response)
	} else { // default return message when error happen
		commonRegistry.GetSentry().CaptureException(err)
		// send notif
		logCtx := getLogCtx(err)
		slackMessage := commonRegistry.GetSlack().GetFormattedMessage(logCtx, ctx, err)
		errSlack := commonRegistry.GetSlack().Send(ctx, slackMessage)
		if errSlack != nil {
			logger.Error(ctx, "Error sending notif to slack", err)
		}

		// send response
		c.JSON(http.StatusInternalServerError, response.Response{
			Message: http.StatusText(http.StatusInternalServerError),
			Status:  response.StatusFail,
		})
	}
}
