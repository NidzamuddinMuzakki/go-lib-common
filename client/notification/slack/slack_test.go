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
		mockSentry := sentryMock.NewISentry(t)
		slackClient := slack.NewSlack(
			validator.New(),
			slack.WithSentry(mockSentry),
			slack.WithSlackConfigChannel(dummy),
			slack.WithSlackConfigNotificationSlackTimeoutInSeconds(
				10,
			),
			slack.WithSlackConfigURL(dummy),
		)

		require.NotNil(t, slackClient)
	})
}

func TestNewSlack_ErrorOnValidation(t *testing.T) {
	t.Run("Error On Validation New Slack", func(t *testing.T) {
		dummy := "dummy"
		mockSentry := sentryMock.NewISentry(t)

		require.Panics(t, func() {
			slack.NewSlack(
				validator.New(),
				slack.WithSentry(mockSentry),
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
		mockSentry := sentryMock.NewISentry(t)
		dummy := "dummy"

		slackClient := slack.NewSlack(
			validator.New(),
			slack.WithSentry(mockSentry),
			slack.WithSlackConfigChannel(dummy),
			slack.WithSlackConfigNotificationSlackTimeoutInSeconds(
				10,
			),
			slack.WithSlackConfigURL(dummy),
		)

		require.NotNil(t, slackClient)

		err := slackClient.Health(context.TODO())
		require.Error(t, err)

	})
}

func TestSend_ErrorOnClientPostMessage(t *testing.T) {
	t.Run("Error On Client Post Message", func(t *testing.T) {
		span := sentry.Span{}
		mockSentry := sentryMock.NewISentry(t)
		mockSentry.On("StartSpan", mock.Anything, mock.Anything).
			Return(&span).
			Once()
		dummy := "dummy"

		slackClient := slack.NewSlack(
			validator.New(),
			slack.WithSentry(mockSentry),
			slack.WithSlackConfigChannel(dummy),
			slack.WithSlackConfigNotificationSlackTimeoutInSeconds(
				10,
			),
			slack.WithSlackConfigURL(dummy),
		)

		require.NotNil(t, slackClient)

		err := slackClient.Send(
			context.TODO(),
			dummy,
		)
		require.Error(t, err)

	})
}
