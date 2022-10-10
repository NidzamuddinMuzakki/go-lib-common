package time

import (
	"os"

	"bitbucket.org/moladinTech/go-lib-common/constant"
)

// LoadTimeZoneFromEnv load timezone from env
// if no default time from env it will return constant.DefaultTimeZone
func LoadTimeZoneFromEnv() string {
	tz := os.Getenv(constant.Timezone)
	if len(tz) <= 0 {
		return constant.DefaultTimeZone
	}
	return tz
}
