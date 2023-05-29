//go:generate mockery --name=INotificationService
package notification_service

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"bitbucket.org/moladinTech/go-lib-common/logger"
	"bitbucket.org/moladinTech/go-lib-common/sentry"
	commonValidator "bitbucket.org/moladinTech/go-lib-common/validator"

	"github.com/go-playground/validator/v10"
	"github.com/parnurzeal/gorequest"
	"github.com/pkg/errors"
)

type INotificationService interface {
	SendSlack(
		ctx context.Context,
		channel string,
		message string,
	) error
	SendEmail(
		ctx context.Context,
		to []string,
		subject string,
		message string,
	) error
	Health(ctx context.Context) error
}

type notificationServicePackage struct {
	sentry           sentry.ISentry `validate:"required"`
	timeoutInSeconds int            `validate:"required"`
	url              string         `validate:"required"`
	serviceName      string
	serviceEnv       string
	client           *gorequest.SuperAgent
}

func WithSentry(sentry sentry.ISentry) Option {
	return func(s *notificationServicePackage) {
		s.sentry = sentry
	}
}

func WithTimeoutInSeconds(timeoutInSeconds int) Option {
	return func(s *notificationServicePackage) {
		s.timeoutInSeconds = timeoutInSeconds
	}
}

func WithURL(url string) Option {
	return func(s *notificationServicePackage) {
		s.url = url
	}
}

func WithServiceName(serviceName string) Option {
	return func(sp *notificationServicePackage) {
		sp.serviceName = serviceName
	}
}

func WithServiceEnv(serviceEnv string) Option {
	return func(sp *notificationServicePackage) {
		sp.serviceEnv = serviceEnv
	}
}

type Response struct {
	Status string `json:"status"`
	Data   any    `json:"data"`
}

var (
	ErrSendNotification = errors.New("failed to send notificationService notification")
	ErrHealthCheck      = errors.New("failed health check notification")
)

type Option func(*notificationServicePackage)

func NewNotificationService(
	validator *validator.Validate,
	options ...Option,
) *notificationServicePackage {
	notificationServicePkg := &notificationServicePackage{}

	for _, option := range options {
		option(notificationServicePkg)
	}

	client := gorequest.New()
	client.Timeout(time.Duration(notificationServicePkg.timeoutInSeconds) * time.Second)
	// TODO: Set auth (if any)
	client.AppendHeader("accept", "application/json")
	client.AppendHeader("Content-Type", "application/json")
	notificationServicePkg.client = client

	err := validator.Struct(notificationServicePkg)
	if err != nil {
		panic(commonValidator.ToErrResponse(err))
	}
	return notificationServicePkg
}

func (c *notificationServicePackage) Health(ctx context.Context) error {
	logCtx := "client.notification_service.notification_service.Ping"
	resp, _, err := c.client.Clone().Get(fmt.Sprintf("%s/%s", c.url, "health")).End()
	if len(err) > 0 {
		logger.Error(ctx, ErrHealthCheck.Error(), err[0], logger.Tag{Key: "logCtx", Value: logCtx})
		return ErrHealthCheck
	}

	if resp.StatusCode != http.StatusOK {
		return ErrHealthCheck
	}

	return nil
}

func (c *notificationServicePackage) SendSlack(
	ctx context.Context,
	channel string,
	message string,
) error {

	var (
		logCtx   = "client.notification_service.notification_service.SendSlack"
		response Response
		span     = c.sentry.StartSpan(ctx, logCtx)
	)
	defer span.Finish()

	resp, _, err := c.client.Clone().Post(fmt.Sprintf("%s/%s", c.url, "slack")).
		SendStruct(struct {
			Channel string `json:"channel"`
			Message string `json:"text"`
		}{
			Channel: channel,
			Message: message,
		}).
		EndStruct(&response)

	if len(err) > 0 {
		logger.Error(ctx, ErrSendNotification.Error(), err[0], logger.Tag{Key: "logCtx", Value: logCtx})
		return ErrSendNotification
	}

	if resp.StatusCode != http.StatusOK {
		return ErrSendNotification
	}

	return nil
}

func (c *notificationServicePackage) SendEmail(
	ctx context.Context,
	to []string,
	subject string,
	message string,
) error {

	var (
		logCtx   = "client.notification_service.notification_service.SendEmail"
		response Response
		span     = c.sentry.StartSpan(ctx, logCtx)
	)
	defer span.Finish()

	resp, _, err := c.client.Clone().Post(fmt.Sprintf("%s/%s", c.url, "sendEmail")).
		Query(struct {
			To      string `json:"to"`
			Subject string `json:"subject"`
			Message string `json:"text"`
		}{
			To:      strings.Join(to, ","),
			Subject: subject,
			Message: message,
		}).
		EndStruct(&response)

	if len(err) > 0 {
		logger.Error(ctx, ErrSendNotification.Error(), err[0], logger.Tag{Key: "logCtx", Value: logCtx})
		return ErrSendNotification
	}

	if resp.StatusCode != http.StatusOK {
		return ErrSendNotification
	}

	return nil
}
