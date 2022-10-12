//go:generate mockery --name=ISlack
package aws

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"bitbucket.org/moladinTech/go-lib-common/logger"
	"bitbucket.org/moladinTech/go-lib-common/sentry"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3 interface {
	UploadFileInByte(ctx context.Context, path, fileName string, data []byte) (string, error)
}

type S3Package struct {
	s3                        *s3.S3
	sentry                    sentry.ISentry
	awsS3Region               string
	awsS3AccessKeyID          string
	awsS3SecretAccessKey      string
	awsS3ARN                  string
	awsS3ACL                  string
	awsS3BucketName           string
	awsS3PresignTimeInMinutes uint
}

type Option func(*S3Package)

func WithS3(s3 *s3.S3) Option {
	return func(s *S3Package) {
		s.s3 = s3
	}
}
func WithSentry(sentry sentry.ISentry) Option {
	return func(s *S3Package) {
		s.sentry = sentry
	}
}
func WithAwsS3Region(awsS3Region string) Option {
	return func(s *S3Package) {
		s.awsS3Region = awsS3Region
	}
}
func WithAwsS3AccessKeyID(awsS3AccessKeyID string) Option {
	return func(s *S3Package) {
		s.awsS3AccessKeyID = awsS3AccessKeyID
	}
}
func WithAwsS3SecretAccessKey(awsS3SecretAccessKey string) Option {
	return func(s *S3Package) {
		s.awsS3SecretAccessKey = awsS3SecretAccessKey
	}
}
func WithAwsS3Arn(awsS3ARN string) Option {
	return func(s *S3Package) {
		s.awsS3ARN = awsS3ARN
	}
}
func WithAwsS3ACL(awsS3ACL string) Option {
	return func(s *S3Package) {
		s.awsS3ACL = awsS3ACL
	}
}
func WithAwsS3BucketName(awsS3BucketName string) Option {
	return func(s *S3Package) {
		s.awsS3BucketName = awsS3BucketName
	}
}
func WithAwsS3PresignTimeInMinutes(awsS3PresignTimeInMinutes uint) Option {
	return func(s *S3Package) {
		s.awsS3PresignTimeInMinutes = awsS3PresignTimeInMinutes
	}
}

func New(
	ctx context.Context,
	options ...Option,
) (*S3Package, error) {
	const logCtx = "S3Package.aws.s3.New"

	s3Pkg := &S3Package{}

	for _, option := range options {
		option(s3Pkg)
	}

	s3, err := s3Pkg.createClient(ctx)
	if err != nil {
		logger.Error(ctx, "failed to create new aws s3 S3Package", err, logger.Tag{Key: "logCtx", Value: logCtx})
		return nil, err
	}
	optionS3 := WithS3(s3)
	optionS3(s3Pkg)

	return s3Pkg, nil
}

func (c *S3Package) createClient(ctx context.Context) (*s3.S3, error) {
	const logCtx = "S3Package.aws.s3.createClient"

	sess, err := session.NewSession()
	if err != nil {
		logger.Error(ctx, "failed to create s3 session", err, logger.Tag{Key: "logCtx", Value: logCtx})
		return nil, err
	}

	awsConfig := &aws.Config{
		Region:                        aws.String(c.awsS3Region),
		CredentialsChainVerboseErrors: aws.Bool(true),
	}

	if c.awsS3AccessKeyID != "" && c.awsS3SecretAccessKey != "" {
		sess.Config.Credentials = credentials.NewStaticCredentials(
			c.awsS3AccessKeyID,
			c.awsS3SecretAccessKey,
			"", // a token will be created when the session it's used.
		)
	}

	if c.awsS3ARN != "" {
		sess.Config.Credentials = stscreds.NewCredentials(sess, c.awsS3ARN)
	}

	return s3.New(sess, awsConfig), nil
}

// UploadFileInByte to upload file in byte to S3
func (c *S3Package) UploadFileInByte(ctx context.Context, path, fileName string, data []byte) (string, error) {
	const logCtx = "S3Package.aws.s3.UploadFileInByte"

	var (
		span        = c.sentry.StartSpan(ctx, logCtx)
		contentType = "application/octet-stream"
		size        = int64(len(data))
	)
	defer span.Finish()

	reqPut, _ := c.s3.PutObjectRequest(&s3.PutObjectInput{
		Bucket:        aws.String(c.awsS3BucketName),
		ACL:           aws.String(c.awsS3ACL),
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

	reqGet, _ := c.s3.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(c.awsS3BucketName),
		Key:    aws.String(fmt.Sprintf("%s/%s", path, fileName)),
	})
	urlStr, err := reqGet.Presign(time.Duration(c.awsS3PresignTimeInMinutes) * time.Minute)
	if err != nil {
		logger.Error(ctx, "failed to presign url", err, logger.Tag{Key: "logCtx", Value: logCtx})
		return "", err
	}

	return urlStr, nil
}
