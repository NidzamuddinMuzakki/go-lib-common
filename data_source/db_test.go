package data_source_test

import (
	"testing"

	commonDataSource "bitbucket.org/moladinTech/go-lib-common/data_source"
	"github.com/stretchr/testify/assert"
)

func Test_GetDSN(t *testing.T) {

	type (
		args struct {
			config *commonDataSource.Config
		}
		want struct {
			dsn string
		}
		testCase struct {
			name string
			args
			want
		}
	)

	testCases := []testCase{
		{
			name: "Mysql case",
			args: args{
				config: &commonDataSource.Config{
					Driver:   "mysql",
					Host:     "localhost",
					Port:     1234,
					DBName:   "dummyDB",
					User:     "user1",
					Password: "password1",
				},
			},
			want: want{
				dsn: "user1:password1@(localhost:1234)/dummyDB?parseTime=true",
			},
		},
		{
			name: "Postgres case",
			args: args{
				config: &commonDataSource.Config{
					Driver:   "postgres",
					Host:     "localhost",
					Port:     1234,
					DBName:   "dummyDB",
					User:     "user1",
					Password: "password1",
					SSLMode:  "disabled",
				},
			},
			want: want{
				dsn: "host=localhost port=1234 user=user1 password=password1 dbname=dummyDB sslmode=disabled",
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			dsn := commonDataSource.GetDsn(
				testCase.config,
			)

			assert.Equal(t, testCase.want.dsn, dsn)
		})
	}
}

func Test_GetDbColumnsAndValue(t *testing.T) {
	var dummyTableAndValue = struct {
		FirstName   string `db:"firstName"`
		LastName    string `db:"lastName"`
		Age         int    `db:"age"`
		IsForeigner bool   `db:"isForeigner"`
	}{
		FirstName:   "John",
		LastName:    "Doe",
		Age:         25,
		IsForeigner: true,
	}

	var dummyResult = map[string]interface{}{
		"firstName":   "John",
		"lastName":    "Doe",
		"age":         25,
		"isForeigner": true,
	}

	type (
		args struct {
			data any
		}
		want struct {
			result map[string]interface{}
		}
		testCase struct {
			name string
			args
			want
		}
	)

	testCases := []testCase{
		{
			name: "Normal Case",
			args: args{
				data: dummyTableAndValue,
			},
			want: want{
				result: dummyResult,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result := commonDataSource.GetDbColumnsAndValue(
				testCase.args.data,
			)

			assert.Equal(t, testCase.want.result, result)
		})
	}
}

func Test_GetDbColumns(t *testing.T) {
	var dummyTableAndValue = struct {
		FirstName   string `db:"firstName"`
		LastName    string `db:"lastName"`
		Age         int    `db:"age"`
		IsForeigner bool   `db:"isForeigner"`
	}{
		FirstName:   "John",
		LastName:    "Doe",
		Age:         25,
		IsForeigner: true,
	}

	var dummyResult = []string{
		"firstName",
		"lastName",
		"age",
		"isForeigner",
	}

	type (
		args struct {
			data any
		}
		want struct {
			result []string
		}
		testCase struct {
			name string
			args
			want
		}
	)

	testCases := []testCase{
		{
			name: "Normal Case",
			args: args{
				data: dummyTableAndValue,
			},
			want: want{
				result: dummyResult,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result := commonDataSource.GetDbColumns(
				testCase.args.data,
			)

			assert.Equal(t, testCase.want.result, result)
		})
	}
}
