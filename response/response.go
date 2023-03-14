package response

import (
	"context"
	"net/http"

	commonError "bitbucket.org/moladinTech/go-lib-common/errors"
	"bitbucket.org/moladinTech/go-lib-common/logger"
	"bitbucket.org/moladinTech/go-lib-common/registry"
	responseModel "bitbucket.org/moladinTech/go-lib-common/response/model"
	"github.com/gin-gonic/gin"
)

// RouteNotFound handle when user is hitting non-exist endpoint.
// It will imediately return error 404 not found.
func RouteNotFound(e *gin.Engine) {
	e.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, responseModel.Response{
			Message: http.StatusText(http.StatusNotFound),
			Status:  responseModel.StatusFail,
		})
	})
}

type ParamHttpErrResp struct {
	Err      error
	GinCtx   *gin.Context
	Registry registry.IRegistry
	Data     interface{}
}

// HttpErrResp is helper to logger the error, send response and send notification (if statusCode >= 500)
func HttpErrResp(ctx context.Context, p ParamHttpErrResp) {
	var (
		c   = p.GinCtx
		e   = p.Err
		rgs = p.Registry

		responseMap, isMapMatch    = commonError.MapErrorResponse[commonError.GetErrKey(e)]
		matchedError, isErrorMatch = commonError.ErrorMatcher(e)
		logCtx                     string
	)

	if isErrorMatch {
		logCtx = matchedError.GetLogCtx()
		logger.Error(ctx, `error`, e, logger.Tag{Key: "logCtx", Value: logCtx})
	} else {
		logger.Error(ctx, `error`, e)
	}

	SendNotifParam := commonError.ParamIsSendNotif{
		IsMapMatch:   isMapMatch,
		ResponseMap:  responseMap,
		IsErrorMatch: isErrorMatch,
		MatchedError: matchedError,
	}

	if commonError.IsCaptureErrorAndSendNotif(SendNotifParam) {
		rgs.GetSentry().CaptureException(e)
		// send notif
		slackMessage := rgs.GetNotif().GetFormattedMessage(ctx, logCtx, e)
		errSlack := rgs.GetNotif().Send(ctx, slackMessage)
		if errSlack != nil {
			logger.Error(ctx, "Error sending notif to slack", errSlack)
		}
	}

	if !isMapMatch {
		c.JSON(http.StatusInternalServerError, responseModel.Response{
			Message: http.StatusText(http.StatusInternalServerError),
			Status:  responseModel.StatusFail,
		})
		return
	}
	// append data
	if p.Data != nil {
		val, ok := responseMap.Response.(responseModel.Response)
		if ok {
			val.Data = p.Data
		}
		responseMap.Response = val
	}

	c.JSON(responseMap.StatusCode, responseMap.Response)

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

		responseMap, isMapMatch    = commonError.MapErrorResponse[commonError.GetErrKey(e)]
		matchedError, isErrorMatch = commonError.ErrorMatcher(e)
		logCtx                     string
	)

	if e == nil && !isErrorMatch {
		return hr
	}

	if isErrorMatch {
		logCtx = matchedError.GetLogCtx()
		logger.Error(ctx, `error`, e, logger.Tag{Key: "logCtx", Value: logCtx})
	} else {
		logger.Error(ctx, `error`, e)
	}

	SendNotifParam := commonError.ParamIsSendNotif{
		IsMapMatch:   isMapMatch,
		ResponseMap:  responseMap,
		IsErrorMatch: isErrorMatch,
		MatchedError: matchedError,
	}

	if commonError.IsCaptureErrorAndSendNotif(SendNotifParam) {
		rgs.GetSentry().CaptureException(e)
		// send notif
		slackMessage := rgs.GetNotif().GetFormattedMessage(ctx, logCtx, e)
		errSlack := rgs.GetNotif().Send(ctx, slackMessage)
		if errSlack != nil {
			logger.Error(ctx, "Error sending notif to slack", errSlack)
		}
	}

	if isErrorMatch && matchedError.GetIsSuccessResp() {
		return hr
	}

	if !isMapMatch {
		c.JSON(http.StatusInternalServerError, responseModel.Response{
			Message: http.StatusText(http.StatusInternalServerError),
			Status:  responseModel.StatusFail,
		})
		return nil
	}

	c.JSON(responseMap.StatusCode, responseMap.Response)
	return nil
}
