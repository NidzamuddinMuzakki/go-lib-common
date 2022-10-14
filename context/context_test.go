package context_test

import (
	"context"
	"reflect"
	"testing"

	commonContext "bitbucket.org/moladinTech/go-lib-common/context"
	"github.com/stretchr/testify/require"
)

func TestContext_ShouldSucceed(t *testing.T) {
	t.Run("Should Succeed Context", func(t *testing.T) {
		key := "myKey"
		val := "myValue"
		ctx := context.WithValue(context.TODO(), key, val)
		actualVal := commonContext.GetValueAsString(ctx, key)
		require.Equal(t, val, actualVal)
	})
}

func TestContext_ErrorOnTypeValue(t *testing.T) {
	t.Run("Error on type value", func(t *testing.T) {
		key := "myKey"
		val := 1
		ctx := context.WithValue(context.TODO(), key, val)
		actualVal := commonContext.GetValueAsString(ctx, key)
		require.NotEqual(t, reflect.TypeOf(val), reflect.TypeOf(actualVal))
	})
}
