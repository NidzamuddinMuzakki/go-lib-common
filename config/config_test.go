package config_test

import (
	"errors"
	"testing"

	"bitbucket.org/moladinTech/go-lib-common/config"
	"github.com/stretchr/testify/assert"
)

func Test_BindFromFile(t *testing.T) {

	type (
		fileTest struct {
			Title       string  `json:"title"`
			Description string  `json:"description"`
			ParamNumber int     `json:"paramNumber"`
			ParamBool   bool    `json:"paramBool"`
			ParamFloat  float64 `json:"paramFloat"`
		}
		args struct {
			dest     any
			filename string
		}
		want struct {
			err error
		}
		testCase struct {
			name string
			args args
			want want
		}
	)

	dummyFileTest := fileTest{}
	invalidDest := ""

	testCases := []testCase{
		{
			name: "Normal case",
			args: args{
				dest:     &dummyFileTest,
				filename: "config.test",
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "Error ReadInConfig",
			args: args{
				dest:     &dummyFileTest,
				filename: "invalidFilename",
			},
			want: want{
				err: errors.New("Not Found in"),
			},
		},
		{
			name: "Error Unmarshal",
			args: args{
				dest:     invalidDest,
				filename: "config.test",
			},
			want: want{
				err: errors.New("result must be a pointer"),
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			err := config.BindFromFile(
				testCase.args.dest,
				testCase.args.filename,
				".",
			)

			if testCase.want.err != nil {
				assert.ErrorContains(t, err, testCase.want.err.Error())
			} else {
				assert.NoError(t, err)
			}

		})
	}
}
