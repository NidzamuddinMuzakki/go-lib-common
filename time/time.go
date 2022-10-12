//go:generate mockery --name=TimeItf
package time

import (
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
