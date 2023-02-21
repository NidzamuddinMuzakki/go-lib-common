package data_source_test

import (
	"testing"

	commonDataSource "bitbucket.org/moladinTech/go-lib-common/data_source"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func Test_IsErrDuplicateKey(t *testing.T) {
	t.Parallel()
	type (
		args struct {
			err error
		}
		want struct {
			isDuplicate bool
		}
		testCase struct {
			name string
			args
			want
		}
	)

	testCases := []testCase{
		{
			name: "Duplicate Case",
			args: args{
				err: &pq.Error{
					Code: "1062",
				},
			},
			want: want{
				isDuplicate: true,
			},
		},
		{
			name: "Not Duplicate Case",
			args: args{
				err: &pq.Error{
					Code: "1234",
				},
			},
			want: want{
				isDuplicate: false,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			isDuplicate := commonDataSource.IsErrDuplicateKey(
				testCase.args.err,
			)

			assert.Equal(t, testCase.want.isDuplicate, isDuplicate)
		})
	}
}
