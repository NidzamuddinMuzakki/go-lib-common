package gcp

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"time"

	"cloud.google.com/go/storage"
	"github.com/pkg/errors"
	"google.golang.org/api/option"

	"bitbucket.org/moladinTech/go-lib-common/logger"
	"bitbucket.org/moladinTech/go-lib-common/sentry"
	commonValidator "bitbucket.org/moladinTech/go-lib-common/validator"
	"github.com/go-playground/validator/v10"
)

const (
	coldLineStorageClass = "COLDLINE"
)

type GCSClient interface {
	UploadFileInByte(ctx context.Context, fileName string, data []byte) (string, error)
	GetSignedURL(context.Context, string, time.Duration) (string, error)
	Upload(context.Context, *UploadOptions) error
	Delete(ctx context.Context, object string, hardDelete bool, timeout time.Duration) error
}

type GCSPackage struct {
	GCSClient              *storage.Client
	Sentry                 sentry.ISentry        `validate:"required"`
	ServiceAccountKeyJSON  ServiceAccountKeyJSON `validate:"required"`
	SignedUrlTimeInMinutes uint                  `validate:"required"`
	BucketName             string                `validate:"required"`
	TimeoutInSeconds       uint
}

type ServiceAccountKeyJSON struct {
	Type                    string `json:"type"`
	ProjectId               string `json:"project_id"`
	PrivateKeyId            string `json:"private_key_id"`
	PrivateKey              string `json:"private_key"`
	ClientEmail             string `json:"client_email"`
	ClientId                string `json:"client_id"`
	AuthUri                 string `json:"auth_uri"`
	TokenUri                string `json:"token_uri"`
	AuthProviderX509CertUrl string `json:"auth_provider_x509_cert_url"`
	ClientX509CertUrl       string `json:"client_x509_cert_url"`
}
type Option func(*GCSPackage)

func withGCSClient(gcsClient *storage.Client) Option {
	return func(s *GCSPackage) {
		s.GCSClient = gcsClient
	}
}
func WithSentry(sentry sentry.ISentry) Option {
	return func(s *GCSPackage) {
		s.Sentry = sentry
	}
}
func WithServiceAccountKeyJSON(serviceAccountKeyJSON ServiceAccountKeyJSON) Option {
	return func(s *GCSPackage) {
		s.ServiceAccountKeyJSON = serviceAccountKeyJSON
	}
}
func WithSignedUrlTimeInMinutes(signedUrlTimeInMinutes uint) Option {
	return func(s *GCSPackage) {
		s.SignedUrlTimeInMinutes = signedUrlTimeInMinutes
	}
}
func WithBucketName(bucketName string) Option {
	return func(s *GCSPackage) {
		s.BucketName = bucketName
	}
}
func WithTimeoutInSeconds(timeoutInSeconds uint) Option {
	return func(s *GCSPackage) {
		s.TimeoutInSeconds = timeoutInSeconds
	}
}
func NewGCS(ctx context.Context,
	validator *validator.Validate,
	options ...Option) *GCSPackage {

	gcsPkg := &GCSPackage{}

	for _, option := range options {
		option(gcsPkg)
	}

	client, err := gcsPkg.createClient(ctx)
	if err != nil {
		panic(err)
	}

	optionGCS := withGCSClient(client)
	optionGCS(gcsPkg)
	err = validator.Struct(gcsPkg)
	if err != nil {
		panic(commonValidator.ToErrResponse(err))
	}
	return gcsPkg
}

func (c *GCSPackage) createClient(ctx context.Context) (*storage.Client, error) {
	const logCtx = "common.client.gcp.storage.createClient"

	reqBodyBytes := new(bytes.Buffer)

	err := json.NewEncoder(reqBodyBytes).Encode(c.ServiceAccountKeyJSON)
	if err != nil {
		logger.Error(ctx, "Failed to encode gcs struct to json", err, logger.Tag{Key: "logCtx", Value: logCtx})
	}
	b := reqBodyBytes.Bytes()

	client, err := storage.NewClient(ctx, option.WithCredentialsJSON(b))
	if err != nil {
		logger.Error(ctx, "Failed to create gcs client", err, logger.Tag{Key: "logCtx", Value: logCtx})
	}

	return client, nil
}

func (g *GCSPackage) UploadFileInByte(ctx context.Context, fileName string, data []byte) (string, error) {
	const logCtx = "common.client.gcp.storage.UploadFileInByte"

	var (
		span             = g.Sentry.StartSpan(ctx, logCtx)
		contentType      = "application/octet-stream"
		timeoutInSeconds uint
	)
	defer span.Finish()
	defer g.GCSClient.Close()

	if g.TimeoutInSeconds > 0 {
		timeoutInSeconds = g.TimeoutInSeconds
	} else {
		timeoutInSeconds = 30
	}
	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeoutInSeconds)*time.Second)
	defer cancel()

	buc := g.GCSClient.Bucket(g.BucketName)
	obj := buc.Object(fileName)
	buf := bytes.NewBuffer(data)

	writer := obj.NewWriter(ctx)
	writer.ChunkSize = 0

	if _, err := io.Copy(writer, buf); err != nil {
		logger.Error(ctx, "Failed to io.Copy", err, logger.Tag{Key: "logCtx", Value: logCtx})
		return "", err
	}
	if err := writer.Close(); err != nil {
		logger.Error(ctx, "Failed to Writer.Close", err, logger.Tag{Key: "logCtx", Value: logCtx})
		return "", err
	}

	if _, err := obj.Update(ctx, storage.ObjectAttrsToUpdate{ContentType: contentType}); err != nil {
		logger.Error(ctx, "Failed to Update Obj Content-Type", err, logger.Tag{Key: "logCtx", Value: logCtx})
	}

	url, err := buc.SignedURL(fileName, &storage.SignedURLOptions{
		GoogleAccessID: g.ServiceAccountKeyJSON.ClientEmail,
		PrivateKey:     []byte(g.ServiceAccountKeyJSON.PrivateKey),
		Method:         "GET",
		Expires:        time.Now().Add(time.Duration(g.SignedUrlTimeInMinutes) * time.Minute),
	})

	if err != nil {
		logger.Error(ctx, "Failed to get SignedUrl", err, logger.Tag{Key: "logCtx", Value: logCtx})
		return "", err
	}

	return url, nil
}

func (g *GCSPackage) GetSignedURL(ctx context.Context, object string, timout time.Duration) (string, error) {
	const logCtx = "common.client.gcp.storage.GetSignedURL"

	span := g.Sentry.StartSpan(ctx, string(logCtx))
	defer g.Sentry.Finish(span)
	ctx = g.Sentry.SpanContext(*span)

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}

	url, err := g.GCSClient.Bucket(g.BucketName).SignedURL(object, &storage.SignedURLOptions{
		GoogleAccessID: g.ServiceAccountKeyJSON.ClientEmail,
		PrivateKey:     []byte(g.ServiceAccountKeyJSON.PrivateKey),
		Method:         "GET",
		Expires:        time.Now().Add(time.Duration(g.SignedUrlTimeInMinutes) * time.Minute),
	})
	if err != nil {
		return "", err
	}

	return url, nil

}

type UploadOptions struct {
	Object  string
	Timeout time.Duration
	File    io.Reader
	// making the object to public
	// the user can directly access the file without authentication
	Public bool
}

func (g *GCSPackage) Upload(ctx context.Context, opt *UploadOptions) error {
	const logCtx = "common.client.gcp.storage.Upload"

	span := g.Sentry.StartSpan(ctx, string(logCtx))
	defer g.Sentry.Finish(span)
	ctx = g.Sentry.SpanContext(*span)

	ctx, cancel := context.WithTimeout(ctx, opt.Timeout)
	defer cancel()

	objWriter := g.GCSClient.Bucket(g.BucketName).Object(opt.Object).NewWriter(ctx)

	if _, err := io.Copy(objWriter, opt.File); err != nil {
		return err
	}

	if err := objWriter.Close(); err != nil {
		return err
	}

	if opt.Public {
		acl := g.GCSClient.Bucket(g.BucketName).Object(opt.Object).ACL()
		acl.Set(ctx, storage.AllUsers, storage.RoleReader)
	}

	return nil
}

func (g *GCSPackage) Delete(ctx context.Context, object string, hardDelete bool, timeout time.Duration) error {
	const logCtx = "common.client.gcp.storage.Delete"

	span := g.Sentry.StartSpan(ctx, string(logCtx))
	defer g.Sentry.Finish(span)
	ctx = g.Sentry.SpanContext(*span)

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	o := g.GCSClient.Bucket(g.BucketName).Object(object)

	// Optional: set a generation-match precondition to avoid potential race
	// conditions and data corruptions. The request to copy is aborted if the
	// object's generation number does not match your precondition.
	attrs, err := o.Attrs(ctx)
	if err != nil {
		return errors.Wrap(err, "object.Attrs")
	}
	o = o.If(storage.Conditions{GenerationMatch: attrs.Generation})

	if hardDelete {
		if err := o.Delete(ctx); err != nil {
			return errors.Wrap(err, "object.HardDelete")
		}
		return nil
	}

	// You can't change an object's storage class directly, the only way is
	// to rewrite the object with the desired storage class.
	copier := o.CopierFrom(o)
	copier.StorageClass = coldLineStorageClass
	if _, err := copier.Run(ctx); err != nil {
		return errors.Wrap(err, "object.Copy.Run")
	}

	return nil
}
