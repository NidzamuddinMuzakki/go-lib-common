package aws_test

import (
	"context"
	"testing"

	"bitbucket.org/moladinTech/go-lib-common/client/aws"
	sentryMock "bitbucket.org/moladinTech/go-lib-common/sentry/mocks"
	"bitbucket.org/moladinTech/go-lib-common/validator"
	"github.com/getsentry/sentry-go"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewS3_ShouldSucceedWithValidation(t *testing.T) {
	t.Run("Should Succeed New S3", func(t *testing.T) {
		sentry := sentryMock.NewISentry(t)
		dummy := "dummy"
		s3 := aws.NewS3(
			context.TODO(),
			validator.New(),
			aws.WithSentry(sentry),
			aws.WithAwsS3Region(dummy),
			aws.WithAwsS3AccessKeyID(dummy),
			aws.WithAwsS3SecretAccessKey(dummy),
			aws.WithAwsS3Arn(dummy),
			aws.WithAwsS3ACL(dummy),
			aws.WithAwsS3BucketName(dummy),
			aws.WithAwsS3PresignTimeInMinutes(10),
		)
		require.NotNil(t, s3)
	})
}

func TestNewS3_ErrorOnValidation(t *testing.T) {
	t.Run("Error On Validation New S3", func(t *testing.T) {
		sentry := sentryMock.NewISentry(t)
		dummy := "dummy"
		require.Panics(t, func() {
			aws.NewS3(
				context.TODO(),
				validator.New(),
				aws.WithSentry(sentry),
				aws.WithAwsS3Region(dummy),
				aws.WithAwsS3SecretAccessKey(dummy),
				aws.WithAwsS3Arn(dummy),
				aws.WithAwsS3ACL(dummy),
				aws.WithAwsS3BucketName(dummy),
				aws.WithAwsS3PresignTimeInMinutes(10),
			)
		})
	})
}

func TestUploadFileInByte_ErrorOnSendPutObjectRequest(t *testing.T) {
	t.Run("Should Succeed UploadFileInByte S3", func(t *testing.T) {
		span := sentry.Span{}
		sentry := sentryMock.NewISentry(t)
		sentry.On("StartSpan", mock.Anything, mock.Anything).
			Return(&span).
			Once()
		dummy := "dummy"
		s3 := aws.NewS3(
			context.TODO(),
			validator.New(),
			aws.WithSentry(sentry),
			aws.WithAwsS3Region(dummy),
			aws.WithAwsS3AccessKeyID(dummy),
			aws.WithAwsS3SecretAccessKey(dummy),
			aws.WithAwsS3Arn(dummy),
			aws.WithAwsS3ACL(dummy),
			aws.WithAwsS3BucketName(dummy),
			aws.WithAwsS3PresignTimeInMinutes(10),
		)

		require.NotNil(t, s3)

		url, err := s3.UploadFileInByte(context.Background(), dummy, dummy, []byte{})
		require.Error(t, err)
		require.Empty(t, url)
	})
}
