package client

import (
	"context"
	"fmt"
	"net/http"

	actLogModel "bitbucket.org/moladinTech/go-lib-activity-log/model"
	"bitbucket.org/moladinTech/go-lib-common/constant"
	"bitbucket.org/moladinTech/go-lib-common/logger"

	"github.com/parnurzeal/gorequest"
	"github.com/pkg/errors"
)

var ErrHealthCheck = errors.New("failed health check moladin-evo")

type MoladinEvo interface {
	Health(ctx context.Context) error
	UserDetail(ctx context.Context, token string) (actLogModel.UserDetail, error)
}

type moladinEvo struct {
	client        *gorequest.SuperAgent
	baseURL       string
	xServicesName string
}

type Option struct {
	BaseURL       string
	XServicesName string
}

func NewMoladinEvo(opt Option) *moladinEvo {

	return &moladinEvo{
		client:        gorequest.New(),
		baseURL:       opt.BaseURL,
		xServicesName: opt.XServicesName,
	}
}

func (c *moladinEvo) Health(ctx context.Context) error {
	const logCtx = "common.client.moladin_evo.Health"

	resp, _, err := c.client.Clone().Get(fmt.Sprintf("%s", c.baseURL)).
		End()

	if len(err) > 0 {
		logger.Error(ctx, ErrHealthCheck.Error(), err[0], logger.Tag{Key: "logCtx", Value: logCtx})
		return err[0]
	}

	if resp.StatusCode != http.StatusOK {
		return ErrHealthCheck
	}

	return nil
}

func (c *moladinEvo) UserDetail(ctx context.Context, token string) (actLogModel.UserDetail, error) {
	const logCtx = "common.client.moladin_evo.UserDetail"

	type Response struct {
		Success      bool                   `json:"success"`
		Message      string                 `json:"message"`
		MessageTitle string                 `json:"messageTitle"`
		Data         actLogModel.UserDetail `json:"data"`
		ResponseTime string                 `json:"responseTime"`
	}

	var res Response
	_, _, err := c.client.Clone().Get(fmt.Sprintf("%s/%s", c.baseURL, "crm/account/user-management/detail")).
		Set(constant.XRequestIdHeader, ctx.Value(constant.XRequestIdHeader).(string)).
		Set(constant.XServiceNameHeader, c.xServicesName).
		Set("authorization", token).
		EndStruct(&res)

	if len(err) > 0 {
		logger.Error(ctx, "fail get user", err[0], logger.Tag{Key: "logCtx", Value: logCtx})
	}

	if !res.Success {
		return res.Data, errors.New(res.Message)
	}

	return res.Data, nil
}
