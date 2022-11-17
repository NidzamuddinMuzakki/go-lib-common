//go:generate mockery --name=TimeItf
package time

import (
	"database/sql/driver"
	"encoding/json"
	"strings"
	"time"
)

type TimeItf interface {
	Now() time.Time
	ToDateTime() string
}

func InitTime() *timeStruct {
	return &timeStruct{}
}

type timeStruct struct {
}

func (t *timeStruct) Now() time.Time {
	return time.Now()
}

func (t *timeStruct) ToDateTime() string {
	return t.Now().In(time.Local).Format("2006-01-02 15:04:05")
}

type DateTime time.Time

func (b *DateTime) UnmarshalJSON(bs []byte) error {
	s := strings.Trim(string(bs), "\"")
	t, err := time.ParseInLocation("2006-01-02 15:04:05", s, LoadTimeZoneAsiaJakarta())
	if err != nil {
		return err
	}
	*b = DateTime(t.UTC())
	return nil
}

func (b DateTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(b).In(LoadTimeZoneAsiaJakarta()))
}

func (b DateTime) Value() (driver.Value, error) {
	return json.Marshal(b)
}
