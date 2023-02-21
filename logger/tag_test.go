package logger_test

import (
	"errors"
	"testing"

	"bitbucket.org/moladinTech/go-lib-common/logger"
	commonLogger "bitbucket.org/moladinTech/go-lib-common/logger"
	"github.com/stretchr/testify/assert"
)

func Test_Err(t *testing.T) {
	t.Parallel()
	err := errors.New("error tag")
	errTag := commonLogger.Err(err)

	assert.Equal(t,
		logger.Tag{
			Key:   "error",
			Value: err.Error(),
		},
		errTag,
	)
}
