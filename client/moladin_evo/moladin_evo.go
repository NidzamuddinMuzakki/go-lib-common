//go:generate mockery --name=IMoladinEvo
package moladin_evo

import (
	"context"
	"fmt"
	"net/http"

	actLogModel "bitbucket.org/moladinTech/go-lib-activity-log/model"
	"bitbucket.org/moladinTech/go-lib-common/constant"
	commonContext "bitbucket.org/moladinTech/go-lib-common/context"
	"bitbucket.org/moladinTech/go-lib-common/logger"
	"bitbucket.org/moladinTech/go-lib-common/sentry"
	commonValidator "bitbucket.org/moladinTech/go-lib-common/validator"
	"github.com/go-playground/validator/v10"
	"github.com/parnurzeal/gorequest"
	"github.com/pkg/errors"
)

var ErrHealthCheck = errors.New("failed health check moladin-evo")

type IMoladinEvo interface {
	Health(ctx context.Context) error
	UserDetail(ctx context.Context, token string) (actLogModel.UserDetail, error)
}

type MoladinEvoPackage struct {
	client        *gorequest.SuperAgent
	Sentry        sentry.ISentry `validate:"required"`
	BaseURL       string         `validate:"required"`
	XServicesName string         `validate:"required"`
}

func WithBaseUrl(baseUrl string) Option {
	return func(s *MoladinEvoPackage) {
		s.BaseURL = baseUrl
	}
}

func WithServicesName(servicesName string) Option {
	return func(s *MoladinEvoPackage) {
		s.XServicesName = servicesName
	}
}

func WithSentry(sentry sentry.ISentry) Option {
	return func(s *MoladinEvoPackage) {
		s.Sentry = sentry
	}
}

type Option func(*MoladinEvoPackage)

func NewMoladinEvo(
	validator *validator.Validate,
	options ...Option,
) *MoladinEvoPackage {
	moladinEvoPkg := &MoladinEvoPackage{
		client: gorequest.New(),
	}

	for _, option := range options {
		option(moladinEvoPkg)
	}

	err := validator.Struct(moladinEvoPkg)
	if err != nil {
		panic(commonValidator.ToErrResponse(err))
	}

	return moladinEvoPkg
}

func (c *MoladinEvoPackage) Health(ctx context.Context) error {
	const logCtx = "common.client.moladin_evo.Health"

	resp, _, err := c.client.Clone().Get(fmt.Sprintf("%s", c.BaseURL)).
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

func (c *MoladinEvoPackage) UserDetail(ctx context.Context, token string) (actLogModel.UserDetail, error) {
	const logCtx = "common.client.moladin_evo.UserDetail"
	span := c.Sentry.StartSpan(ctx, logCtx)
	defer span.Finish()

	type Response struct {
		Success      bool                   `json:"success"`
		Message      string                 `json:"message"`
		MessageTitle string                 `json:"messageTitle"`
		Data         actLogModel.UserDetail `json:"data"`
		ResponseTime string                 `json:"responseTime"`
	}

	var res Response
	_, _, err := c.client.Clone().Get(fmt.Sprintf("%s/%s", c.BaseURL, "crm/account/user-management/detail")).
		Set(constant.XRequestIdHeader, commonContext.GetValueAsString(ctx, constant.XRequestIdHeader)).
		Set(constant.XServiceNameHeader, c.XServicesName).
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
