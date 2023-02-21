package logger_test

import (
	"context"
	"testing"

	commonLogger "bitbucket.org/moladinTech/go-lib-common/logger"
	"github.com/stretchr/testify/assert"
)

func Test_AddLoggingTag(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctx = context.WithValue(ctx, commonLogger.LoggingTagKey, "tag wrapper")

	newCtx := commonLogger.AddLoggingTag(ctx,
		commonLogger.Tag{
			Key:   "tag2",
			Value: "this is tag 2",
		},
	)

	expectedCtxValue := map[string]string{
		"tag2": "this is tag 2",
	}
	assert.Equal(t, expectedCtxValue, newCtx.Value(commonLogger.LoggingTagKey))
}

func Test_AddRequestID(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	newCtx := commonLogger.AddRequestID(
		ctx,
		"12345",
	)
	expectedCtxValue := map[string]string{
		commonLogger.RequestIDKey: "12345",
	}

	assert.Equal(
		t,
		expectedCtxValue,
		newCtx.Value(commonLogger.LoggingTagKey),
	)
}

func Test_GetAllLoggingTagInTagStr(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	newCtx := commonLogger.AddRequestID(
		ctx,
		"12345",
	)

	tags := commonLogger.GetAllLoggingTagInTagStr(newCtx)
	assert.Equal(
		t,
		[]commonLogger.Tag{
			{
				Key:   commonLogger.RequestIDKey,
				Value: "12345",
			},
		},
		tags)
}
