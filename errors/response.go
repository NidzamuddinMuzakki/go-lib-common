package errors

import (
	"context"
	"net/http"

	"bitbucket.org/moladinTech/go-lib-common/logger"
	"bitbucket.org/moladinTech/go-lib-common/registry"
	"bitbucket.org/moladinTech/go-lib-common/response"
	"github.com/gin-gonic/gin"
)

type ParamHttpErrResp struct {
	Err      error
	GinCtx   *gin.Context
	Registry registry.IRegistry
}

// HttpErrResp is helper to logger the error, send response and send notification (if statusCode >= 500)
func HttpErrResp(ctx context.Context, p ParamHttpErrResp) {
	var (
		c   = p.GinCtx
		e   = p.Err
		rgs = p.Registry

		v1, ok1 = MapErrorResponse[GetErrKey(e)]
		v2, ok2 = e.(*err)
	)

	if ok2 {
		logger.Error(ctx, `error`, e, logger.Tag{Key: "logCtx", Value: v2.logCtx})
	} else {
		logger.Error(ctx, `error`, e)
	}

	if ok2 && (!ok1 || (ok1 && v1.StatusCode >= 500) || v2.isNotify) {
		rgs.GetSentry().CaptureException(e)
		// send notif
		slackMessage := rgs.GetNotif().GetFormattedMessage(v2.logCtx, ctx, e)
		errSlack := rgs.GetNotif().Send(ctx, slackMessage)
		if errSlack != nil {
			logger.Error(ctx, "Error sending notif to slack", errSlack)
		}
	}

	if !ok1 {
		c.JSON(http.StatusInternalServerError, response.Response{
			Message: http.StatusText(http.StatusInternalServerError),
			Status:  response.StatusFail,
		})
		return
	}

	c.JSON(v1.StatusCode, v1.Response)
	return
}
