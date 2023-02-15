package response_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	slackMock "bitbucket.org/moladinTech/go-lib-common/client/notification/slack/mocks"
	"bitbucket.org/moladinTech/go-lib-common/registry"
	commonResponse "bitbucket.org/moladinTech/go-lib-common/response"
	sentryMock "bitbucket.org/moladinTech/go-lib-common/sentry/mocks"
	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_HttpRespNoError(t *testing.T) {
	httpResp := commonResponse.HttpResp(context.Background(), nil, commonResponse.ParamHttpErrResp{
		Registry: registry.NewRegistry(),
		GinCtx:   &gin.Context{},
	})

	assert.Equal(t, &gin.Context{}, httpResp.GinCtx)
}

func Test_HttpRespWithError(t *testing.T) {
	mockedSentry := sentryMock.NewISentry(t)
	mockedSlack := slackMock.NewISlack(t)

	dummyError := errors.New("error")
	registryClient := registry.NewRegistry(
		registry.WithSlack(mockedSlack),
		registry.WithSentry(mockedSentry),
		registry.WithNotif(mockedSlack),
	)

	mockedSentry.On("CaptureException", dummyError).
		Return(&sentry.NewEvent().EventID)
	mockedSlack.On("GetFormattedMessage",
		context.Background(), mock.Anything, dummyError,
	).Return("Hello")
	mockedSlack.On("Send",
		context.Background(), "Hello",
	).Return(nil)

	httpRecorder := httptest.NewRecorder()
	ginCtx, _ := gin.CreateTestContext(httpRecorder)
	paramHttpResp := commonResponse.ParamHttpErrResp{
		Err:      dummyError,
		GinCtx:   ginCtx,
		Registry: registryClient,
	}

	err := commonResponse.HttpResp(
		context.Background(),
		dummyError,
		paramHttpResp,
	)

	assert.Nil(t, err)
}

func Test_HttpErrRespWithError(t *testing.T) {
	mockedSentry := sentryMock.NewISentry(t)
	mockedSlack := slackMock.NewISlack(t)

	dummyError := errors.New("error")
	registryClient := registry.NewRegistry(
		registry.WithSlack(mockedSlack),
		registry.WithSentry(mockedSentry),
		registry.WithNotif(mockedSlack),
	)

	mockedSentry.On("CaptureException", dummyError).
		Return(&sentry.NewEvent().EventID)
	mockedSlack.On("GetFormattedMessage",
		context.Background(), mock.Anything, dummyError,
	).Return("Hello")
	mockedSlack.On("Send",
		context.Background(), "Hello",
	).Return(nil)

	httpRecorder := httptest.NewRecorder()
	ginCtx, _ := gin.CreateTestContext(httpRecorder)
	paramHttpResp := commonResponse.ParamHttpErrResp{
		Err:      dummyError,
		GinCtx:   ginCtx,
		Registry: registryClient,
	}

	commonResponse.HttpErrResp(
		context.Background(),
		paramHttpResp,
	)

	assert.Equal(t, ginCtx.Writer.Status(), http.StatusInternalServerError)
}
