package strings_test

import (
	"reflect"
	"testing"

	commonString "bitbucket.org/moladinTech/go-lib-common/strings"
	"github.com/stretchr/testify/require"
)

func TestToUint64_ShouldSucceed(t *testing.T) {
	t.Parallel()

	t.Run("Should Succeed ToUint64", func(t *testing.T) {
		actualVal := commonString.ToUint64("1")
		require.Equal(t, uint64(1), actualVal)
		require.Equal(t, "uint64", reflect.TypeOf(actualVal).String())
	})
}

func TestToUint64_ShouldSucceedWithEmpty(t *testing.T) {
	t.Parallel()

	t.Run("Should Succeed ToUint64", func(t *testing.T) {
		actualVal := commonString.ToUint64("")
		require.Equal(t, uint64(0), actualVal)
		require.Equal(t, "uint64", reflect.TypeOf(actualVal).String())
	})
}

func TestToFloat64_ShouldSucceed(t *testing.T) {
	t.Parallel()

	t.Run("Should Succeed ToFloat64", func(t *testing.T) {
		actualVal := commonString.ToFloat64("1")
		require.Equal(t, float64(1), actualVal)
		require.Equal(t, "float64", reflect.TypeOf(actualVal).String())
	})
}

func TestToFloat64_ShouldSucceedWithEmpty(t *testing.T) {
	t.Parallel()

	t.Run("Should Succeed ToFloat64", func(t *testing.T) {
		actualVal := commonString.ToFloat64("")
		require.Equal(t, float64(0), actualVal)
		require.Equal(t, "float64", reflect.TypeOf(actualVal).String())
	})
}
