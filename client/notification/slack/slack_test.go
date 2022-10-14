package slack_test

import (
	"context"
	"testing"

	"bitbucket.org/moladinTech/go-lib-common/client/notification/slack"
	sentryMock "bitbucket.org/moladinTech/go-lib-common/sentry/mocks"
	"bitbucket.org/moladinTech/go-lib-common/validator"
	"github.com/getsentry/sentry-go"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewSlack_ShouldSucceedWithValidation(t *testing.T) {
	t.Run("Should Succeed New Slack", func(t *testing.T) {
		dummy := "dummy"
		sentry := sentryMock.NewISentry(t)
		slack := slack.NewSlack(
			validator.New(),
			slack.WithSentry(sentry),
			slack.WithSlackConfigChannel(dummy),
			slack.WithSlackConfigNotificationSlackTimeoutInSeconds(
				10,
			),
			slack.WithSlackConfigURL(dummy),
		)

		require.NotNil(t, slack)
	})
}

func TestNewSlack_ErrorOnValidation(t *testing.T) {
	t.Run("Error On Validation New Slack", func(t *testing.T) {
		dummy := "dummy"
		sentry := sentryMock.NewISentry(t)

		require.Panics(t, func() {
			slack.NewSlack(
				validator.New(),
				slack.WithSentry(sentry),
				slack.WithSlackConfigNotificationSlackTimeoutInSeconds(
					10,
				),
				slack.WithSlackConfigURL(dummy),
			)
		})
	})
}

func TestHealth_ErrorOnClientGetUrl(t *testing.T) {
	t.Run("Error On Client Get Url", func(t *testing.T) {
		span := sentry.Span{}
		sentry := sentryMock.NewISentry(t)
		sentry.On("StartSpan", mock.Anything, mock.Anything).
			Return(&span).
			Once()
		dummy := "dummy"

		slack := slack.NewSlack(
			validator.New(),
			slack.WithSentry(sentry),
			slack.WithSlackConfigChannel(dummy),
			slack.WithSlackConfigNotificationSlackTimeoutInSeconds(
				10,
			),
			slack.WithSlackConfigURL(dummy),
		)

		require.NotNil(t, slack)

		err := slack.Health(context.TODO())
		require.Error(t, err)

	})
}

func TestSend_ErrorOnClientPostMessage(t *testing.T) {
	t.Run("Error On Client Post Message", func(t *testing.T) {
		span := sentry.Span{}
		sentry := sentryMock.NewISentry(t)
		sentry.On("StartSpan", mock.Anything, mock.Anything).
			Return(&span).
			Once()
		dummy := "dummy"

		slack := slack.NewSlack(
			validator.New(),
			slack.WithSentry(sentry),
			slack.WithSlackConfigChannel(dummy),
			slack.WithSlackConfigNotificationSlackTimeoutInSeconds(
				10,
			),
			slack.WithSlackConfigURL(dummy),
		)

		require.NotNil(t, slack)

		err := slack.Send(
			context.TODO(),
			dummy,
		)
		require.Error(t, err)

	})
}
