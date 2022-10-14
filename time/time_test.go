package time_test

import (
	"testing"
	"time"

	commonTime "bitbucket.org/moladinTech/go-lib-common/time"
	"github.com/stretchr/testify/require"
)

func TestNewTime_ShouldSucceedWithNow(t *testing.T) {
	t.Run("Should Succeed New Time", func(t *testing.T) {
		commonTime.LoadTimeZoneFromEnv()
		tm := commonTime.InitTime()
		tmDay := tm.Now().Day()
		timeDay := time.Now().Day()
		require.Equal(t, timeDay, tmDay)
	})
}

func TestNewTime_ShouldSucceedWithToDateTime(t *testing.T) {
	t.Run("Should Succeed New Time", func(t *testing.T) {
		commonTime.LoadTimeZoneFromEnv()
		tm := commonTime.InitTime()
		tmDay := tm.ToDateTime()
		require.NotEmpty(t, tmDay)
	})
}
