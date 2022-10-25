//go:generate mockery --name=ISlack
package aws

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"bitbucket.org/moladinTech/go-lib-common/logger"
	"bitbucket.org/moladinTech/go-lib-common/sentry"
	commonValidator "bitbucket.org/moladinTech/go-lib-common/validator"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/go-playground/validator/v10"
)

type S3 interface {
	UploadFileInByte(ctx context.Context, path, fileName string, data []byte) (string, error)
}

type S3Package struct {
	S3                        *s3.S3
	Sentry                    sentry.ISentry `validate:"required"`
	AwsS3Region               string         `validate:"required"`
	AwsS3AccessKeyID          string         `validate:"required"`
	AwsS3SecretAccessKey      string         `validate:"required"`
	AwsS3ACL                  string         `validate:"required"`
	AwsS3BucketName           string         `validate:"required"`
	AwsS3PresignTimeInMinutes uint           `validate:"required"`
	AwsS3ARN                  string
}

type Option func(*S3Package)

func withS3(s3 *s3.S3) Option {
	return func(s *S3Package) {
		s.S3 = s3
	}
}
func WithSentry(sentry sentry.ISentry) Option {
	return func(s *S3Package) {
		s.Sentry = sentry
	}
}
func WithAwsS3Region(awsS3Region string) Option {
	return func(s *S3Package) {
		s.AwsS3Region = awsS3Region
	}
}
func WithAwsS3AccessKeyID(awsS3AccessKeyID string) Option {
	return func(s *S3Package) {
		s.AwsS3AccessKeyID = awsS3AccessKeyID
	}
}
func WithAwsS3SecretAccessKey(awsS3SecretAccessKey string) Option {
	return func(s *S3Package) {
		s.AwsS3SecretAccessKey = awsS3SecretAccessKey
	}
}
func WithAwsS3Arn(awsS3ARN string) Option {
	return func(s *S3Package) {
		s.AwsS3ARN = awsS3ARN
	}
}
func WithAwsS3ACL(awsS3ACL string) Option {
	return func(s *S3Package) {
		s.AwsS3ACL = awsS3ACL
	}
}
func WithAwsS3BucketName(awsS3BucketName string) Option {
	return func(s *S3Package) {
		s.AwsS3BucketName = awsS3BucketName
	}
}
func WithAwsS3PresignTimeInMinutes(awsS3PresignTimeInMinutes uint) Option {
	return func(s *S3Package) {
		s.AwsS3PresignTimeInMinutes = awsS3PresignTimeInMinutes
	}
}

func NewS3(
	ctx context.Context,
	validator *validator.Validate,
	options ...Option,
) *S3Package {

	s3Pkg := &S3Package{}

	for _, option := range options {
		option(s3Pkg)
	}

	s3, err := s3Pkg.createClient(ctx)
	if err != nil {
		panic(err)
	}

	optionS3 := withS3(s3)
	optionS3(s3Pkg)

	err = validator.Struct(s3Pkg)
	if err != nil {
		panic(commonValidator.ToErrResponse(err))
	}

	return s3Pkg
}

func (c *S3Package) createClient(ctx context.Context) (*s3.S3, error) {
	const logCtx = "common.client.aws.s3.createClient"

	sess, err := session.NewSession()
	if err != nil {
		logger.Error(ctx, "failed to create s3 session", err, logger.Tag{Key: "logCtx", Value: logCtx})
		return nil, err
	}

	awsConfig := &aws.Config{
		Region:                        aws.String(c.AwsS3Region),
		CredentialsChainVerboseErrors: aws.Bool(true),
	}

	if c.AwsS3AccessKeyID != "" && c.AwsS3SecretAccessKey != "" {
		sess.Config.Credentials = credentials.NewStaticCredentials(
			c.AwsS3AccessKeyID,
			c.AwsS3SecretAccessKey,
			"", // a token will be created when the session it's used.
		)
	}

	if c.AwsS3ARN != "" {
		sess.Config.Credentials = stscreds.NewCredentials(sess, c.AwsS3ARN)
	}

	return s3.New(sess, awsConfig), nil
}

// UploadFileInByte to upload file in byte to S3
func (c *S3Package) UploadFileInByte(ctx context.Context, path, fileName string, data []byte) (string, error) {
	const logCtx = "common.client.aws.s3.UploadFileInByte"

	var (
		span        = c.Sentry.StartSpan(ctx, logCtx)
		contentType = "application/octet-stream"
		size        = int64(len(data))
	)
	defer span.Finish()
	reqPut, _ := c.S3.PutObjectRequest(&s3.PutObjectInput{
		Bucket:        aws.String(c.AwsS3BucketName),
		ACL:           aws.String(c.AwsS3ACL),
		Key:           aws.String(fmt.Sprintf("%s/%s", path, fileName)),
		ContentLength: aws.Int64(size),
		ContentType:   aws.String(contentType),
		Body:          bytes.NewReader(data),
	})
	err := reqPut.Send()
	if err != nil {
		logger.Error(ctx, "failed to upload file to s3", err, logger.Tag{Key: "logCtx", Value: logCtx})
		return "", err
	}

	reqGet, _ := c.S3.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(c.AwsS3BucketName),
		Key:    aws.String(fmt.Sprintf("%s/%s", path, fileName)),
	})
	urlStr, err := reqGet.Presign(time.Duration(c.AwsS3PresignTimeInMinutes) * time.Minute)
	if err != nil {
		logger.Error(ctx, "failed to presign url", err, logger.Tag{Key: "logCtx", Value: logCtx})
		return "", err
	}

	return urlStr, nil
}
