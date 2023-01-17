package moladin_evo_test

import (
	"context"
	"testing"

	"bitbucket.org/moladinTech/go-lib-common/client/moladin_evo"
	sentryMock "bitbucket.org/moladinTech/go-lib-common/sentry/mocks"
	"bitbucket.org/moladinTech/go-lib-common/validator"
	"github.com/getsentry/sentry-go"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewMoladinEvo_ShouldSucceedWithValidation(t *testing.T) {
	t.Run("Should Succeed New MoladinEvo", func(t *testing.T) {
		dummy := "dummy"
		mockSentry := sentryMock.NewISentry(t)
		moladinEvo := moladin_evo.NewMoladinEvo(
			validator.New(),
			moladin_evo.WithSentry(mockSentry),
			moladin_evo.WithBaseUrl(dummy),
			moladin_evo.WithServicesName(dummy),
		)
		require.NotNil(t, moladinEvo)
	})
}

func TestNewMoladinEvo_ErrorOnValidation(t *testing.T) {
	t.Run("Error On Validation New Moladin Evo", func(t *testing.T) {
		dummy := "dummy"
		mockSentry := sentryMock.NewISentry(t)
		require.Panics(t, func() {
			moladin_evo.NewMoladinEvo(
				validator.New(),
				moladin_evo.WithSentry(mockSentry),
				moladin_evo.WithBaseUrl(dummy),
			)
		})
	})
}

func TestHealth_ErrorOnClientGetUrl(t *testing.T) {
	t.Run("Error On Client Get Url", func(t *testing.T) {
		mockSentry := sentryMock.NewISentry(t)
		dummy := "dummy"

		moladinEvo := moladin_evo.NewMoladinEvo(
			validator.New(),
			moladin_evo.WithSentry(mockSentry),
			moladin_evo.WithBaseUrl(dummy),
			moladin_evo.WithServicesName(dummy),
		)
		require.NotNil(t, moladinEvo)

		err := moladinEvo.Health(context.TODO())
		require.Error(t, err)

	})
}

func TestUserDetail_ErrorOnClientGetUrlUser(t *testing.T) {
	t.Run("Error On Client Get Url User", func(t *testing.T) {
		span := sentry.Span{}
		mockSentry := sentryMock.NewISentry(t)
		mockSentry.On("StartSpan", mock.Anything, mock.Anything).
			Return(&span).
			Once()
		dummy := "dummy"

		moladinEvo := moladin_evo.NewMoladinEvo(
			validator.New(),
			moladin_evo.WithSentry(mockSentry),
			moladin_evo.WithBaseUrl(dummy),
			moladin_evo.WithServicesName(dummy),
		)
		require.NotNil(t, moladinEvo)

		userDetail, err := moladinEvo.UserDetail(
			context.TODO(),
			dummy,
		)
		require.Error(t, err)
		require.Empty(t, userDetail)

	})
}
