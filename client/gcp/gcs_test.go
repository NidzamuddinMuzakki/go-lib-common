package gcp_test

import (
	"context"
	"testing"

	"bitbucket.org/moladinTech/go-lib-common/client/gcp"
	sentryMock "bitbucket.org/moladinTech/go-lib-common/sentry/mocks"
	"bitbucket.org/moladinTech/go-lib-common/validator"
	"github.com/getsentry/sentry-go"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewGCSCLient_ShouldSucceedWithValidation(t *testing.T) {
	t.Parallel()
	t.Run("Should Succeed New GCS", func(t *testing.T) {
		sentry := sentryMock.NewISentry(t)
		dummy := "dummy"
		serviceAccount := gcp.ServiceAccountKeyJSON{
			Type:                    "service_account",
			ProjectId:               dummy,
			PrivateKeyId:            dummy,
			PrivateKey:              dummy,
			ClientEmail:             dummy,
			ClientId:                dummy,
			AuthUri:                 dummy,
			TokenUri:                dummy,
			AuthProviderX509CertUrl: dummy,
			ClientX509CertUrl:       dummy,
		}
		g := gcp.NewGCS(
			context.TODO(),
			validator.New(),
			gcp.WithSentry(sentry),
			gcp.WithServiceAccountKeyJSON(serviceAccount),
			gcp.WithSignedUrlTimeInMinutes(10),
			gcp.WithTimeoutInSeconds(30),
			gcp.WithBucketName(dummy),
		)
		require.NotNil(t, g)
	})
}

func TestNewGCS_ErrorOnValidation(t *testing.T) {
	t.Parallel()
	t.Run("Error On Validation New GCS", func(t *testing.T) {
		sentry := sentryMock.NewISentry(t)
		dummy := "dummy"
		serviceAccount := gcp.ServiceAccountKeyJSON{
			Type:                    "service_account",
			ProjectId:               dummy,
			PrivateKeyId:            dummy,
			PrivateKey:              dummy,
			ClientEmail:             dummy,
			ClientId:                dummy,
			AuthUri:                 dummy,
			TokenUri:                dummy,
			AuthProviderX509CertUrl: dummy,
			ClientX509CertUrl:       dummy,
		}
		require.Panics(t, func() {
			gcp.NewGCS(
				context.TODO(),
				validator.New(),
				gcp.WithSentry(sentry),
				gcp.WithServiceAccountKeyJSON(serviceAccount),
				gcp.WithBucketName(dummy),
			)
		})
	})
}

func TestUploadFileInByte_ErrorOnSendPutObjectRequest(t *testing.T) {
	t.Parallel()
	t.Run("Should Succeed UploadFileInByte GCS", func(t *testing.T) {
		span := sentry.Span{}
		sentry := sentryMock.NewISentry(t)
		sentry.On("StartSpan", mock.Anything, mock.Anything).
			Return(&span).
			Once()
		dummy := "dummy"
		serviceAccount := gcp.ServiceAccountKeyJSON{
			Type:                    "service_account",
			ProjectId:               dummy,
			PrivateKeyId:            dummy,
			PrivateKey:              dummy,
			ClientEmail:             dummy,
			ClientId:                dummy,
			AuthUri:                 dummy,
			TokenUri:                dummy,
			AuthProviderX509CertUrl: dummy,
			ClientX509CertUrl:       dummy,
		}
		g := gcp.NewGCS(
			context.TODO(),
			validator.New(),
			gcp.WithSentry(sentry),
			gcp.WithServiceAccountKeyJSON(serviceAccount),
			gcp.WithSignedUrlTimeInMinutes(10),
			gcp.WithBucketName(dummy),
		)
		require.NotNil(t, g)

		url, err := g.UploadFileInByte(context.Background(), dummy, []byte{})
		require.Error(t, err)
		require.Empty(t, url)
	})
}
