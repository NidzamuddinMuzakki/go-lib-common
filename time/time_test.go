package time_test

import (
	"encoding/json"
	"testing"
	"time"

	commonTime "bitbucket.org/moladinTech/go-lib-common/time"
	"github.com/stretchr/testify/assert"
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

func Test_UnmarshalJSON(t *testing.T) {
	var dateTime = struct {
		CommonTime commonTime.DateTime `json:"time"`
	}{}
	wantDateTime := commonTime.DateTime(time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC))

	b := []byte(`{"time": "2006-01-02 22:04:05"}`)
	err := json.Unmarshal(b, &dateTime)

	assert.Nil(t, err)
	assert.Equal(t, wantDateTime, dateTime.CommonTime)
}

func Test_GetValue(t *testing.T) {
	timeNow := time.Now()
	dateTime := commonTime.GetValue(&timeNow, commonTime.LoadTimeZoneAsiaJakarta())

	assert.NotNil(t, dateTime)
}
