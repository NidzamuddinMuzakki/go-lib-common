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
		logCtx  string
	)

	if ok2 {
		logCtx = v2.logCtx
		logger.Error(ctx, `error`, e, logger.Tag{Key: "logCtx", Value: logCtx})
	} else {
		logger.Error(ctx, `error`, e)
	}

	if !ok1 || (ok1 && v1.StatusCode >= 500) || v2.isNotify {
		rgs.GetSentry().CaptureException(e)
		// send notif
		slackMessage := rgs.GetNotif().GetFormattedMessage(logCtx, ctx, e)
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

type httpResp struct {
	GinCtx *gin.Context
}

func (h *httpResp) Return(statusCode int, response interface{}) {
	if h != nil {
		h.GinCtx.JSON(statusCode, response)
	}
}

// HttpResp is helper to logger the error, send response and send notification (if statusCode >= 500)
func HttpResp(ctx context.Context, e error, p ParamHttpErrResp) *httpResp {
	var (
		c   = p.GinCtx
		rgs = p.Registry
		hr  = &httpResp{GinCtx: c}

		v1, ok1 = MapErrorResponse[GetErrKey(e)]
		v2, ok2 = e.(*err)
		logCtx  string
	)

	if e == nil && !ok2 {
		return hr
	}

	if ok2 {
		logCtx = v2.logCtx
		logger.Error(ctx, `error`, e, logger.Tag{Key: "logCtx", Value: logCtx})
	} else {
		logger.Error(ctx, `error`, e)
	}

	if !ok1 || (ok1 && v1.StatusCode >= 500) || v2.isNotify {
		rgs.GetSentry().CaptureException(e)
		// send notif
		slackMessage := rgs.GetNotif().GetFormattedMessage(logCtx, ctx, e)
		errSlack := rgs.GetNotif().Send(ctx, slackMessage)
		if errSlack != nil {
			logger.Error(ctx, "Error sending notif to slack", errSlack)
		}
	}

	if ok2 && v2.isSuccessResp {
		return hr
	}

	if !ok1 {
		c.JSON(http.StatusInternalServerError, response.Response{
			Message: http.StatusText(http.StatusInternalServerError),
			Status:  response.StatusFail,
		})
		return nil
	}

	c.JSON(v1.StatusCode, v1.Response)
	return nil
}
