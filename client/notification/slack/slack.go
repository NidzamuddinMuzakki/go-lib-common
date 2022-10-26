//go:generate mockery --name=ISlack
package slack

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"bitbucket.org/moladinTech/go-lib-common/logger"
	"bitbucket.org/moladinTech/go-lib-common/sentry"
	commonValidator "bitbucket.org/moladinTech/go-lib-common/validator"

	"github.com/go-playground/validator/v10"
	"github.com/parnurzeal/gorequest"
	"github.com/pkg/errors"
)

type ISlack interface {
	Send(ctx context.Context, message string) error
	Health(ctx context.Context) error
}

type SlackPackage struct {
	Sentry                                       sentry.ISentry `validate:"required"`
	SlackConfigNotificationSlackTimeoutInSeconds int            `validate:"required"`
	SlackConfigURL                               string         `validate:"required"`
	SlackConfigChannel                           string         `validate:"required"`
	client                                       *gorequest.SuperAgent
}

func WithSentry(sentry sentry.ISentry) Option {
	return func(s *SlackPackage) {
		s.Sentry = sentry
	}
}
func WithSlackConfigNotificationSlackTimeoutInSeconds(slackConfigNotificationSlackTimeoutInSeconds int) Option {
	return func(s *SlackPackage) {
		s.SlackConfigNotificationSlackTimeoutInSeconds = slackConfigNotificationSlackTimeoutInSeconds
	}
}
func WithSlackConfigURL(slackConfigURL string) Option {
	return func(s *SlackPackage) {
		s.SlackConfigURL = slackConfigURL
	}
}
func WithSlackConfigChannel(slackConfigChannel string) Option {
	return func(s *SlackPackage) {
		s.SlackConfigChannel = slackConfigChannel
	}
}

type Payload struct {
	Channel string `json:"channel"`
	Message string `json:"text"`
}

type Response struct {
	Status string `json:"status"`
	Data   any    `json:"data"`
}

type Config struct {
	NotificationSlackTimeoutInSeconds int    `json:"notificationSlackTimeoutInSeconds" yaml:"notificationSlackTimeoutInSeconds"`
	Channel                           string `json:"channel" yaml:"channel"`
	URL                               string `json:"url" yaml:"url"`
}

type LogCtx struct {
	ClientNotificationSlackSend string
	ClientNotificationSlackPing string
}

var (
	ErrSendNotification = errors.New("failed to send slack notification")
	ErrHealthCheck      = errors.New("failed health check notification")
	LogCtxName          = LogCtx{
		ClientNotificationSlackSend: "client.notification.slack.Send",
		ClientNotificationSlackPing: "client.notification.slack.Ping",
	}
)

type Option func(*SlackPackage)

func NewSlack(
	validator *validator.Validate,
	options ...Option,
) ISlack {
	slackPkg := &SlackPackage{}

	for _, option := range options {
		option(slackPkg)
	}

	client := gorequest.New()
	client.Timeout(time.Duration(slackPkg.SlackConfigNotificationSlackTimeoutInSeconds) * time.Second)
	// TODO: Set auth (if any)
	client.AppendHeader("accept", "application/json")
	client.AppendHeader("Content-Type", "application/json")
	slackPkg.client = client

	err := validator.Struct(slackPkg)
	if err != nil {
		panic(commonValidator.ToErrResponse(err))
	}
	return slackPkg
}

func (c *SlackPackage) Health(ctx context.Context) error {
	var span = c.Sentry.StartSpan(ctx, LogCtxName.ClientNotificationSlackPing)
	defer span.Finish()

	resp, _, err := c.client.Clone().Get(fmt.Sprintf("%s/%s", c.SlackConfigURL, "health")).End()
	if len(err) > 0 {
		logger.Error(ctx, ErrHealthCheck.Error(), err[0], logger.Tag{Key: "logCtx", Value: LogCtxName.ClientNotificationSlackPing})
		return ErrHealthCheck
	}

	if resp.StatusCode != http.StatusOK {
		return ErrHealthCheck
	}

	return nil
}

func (c *SlackPackage) Send(ctx context.Context, message string) error {

	var (
		response Response
		span     = c.Sentry.StartSpan(ctx, LogCtxName.ClientNotificationSlackSend)
	)
	defer span.Finish()

	resp, _, err := c.client.Clone().Post(fmt.Sprintf("%s/%s", c.SlackConfigURL, "slack")).
		SendStruct(Payload{
			Channel: c.SlackConfigChannel,
			Message: message,
		}).
		EndStruct(&response)

	if len(err) > 0 {
		logger.Error(ctx, ErrSendNotification.Error(), err[0], logger.Tag{Key: "logCtx", Value: LogCtxName.ClientNotificationSlackSend})
		return ErrSendNotification
	}

	if resp.StatusCode != http.StatusOK {
		return ErrSendNotification
	}

	return nil
}