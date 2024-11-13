package slices_test

import (
	"fmt"
	"testing"

	"bitbucket.org/moladinTech/go-lib-common/slices"
	"github.com/stretchr/testify/suite"
)

type SlicesTestSuite struct {
	suite.Suite
}

func (suite *SlicesTestSuite) SetupSuite() {
	fmt.Println("SetupSuite: SlicesTestSuite")
}

func (suite *SlicesTestSuite) TearDownSuite() {
	fmt.Println("TearDownSuite: SlicesTestSuite")
}

func (suite *SlicesTestSuite) SetupTest() {
	fmt.Println("SlicesTestSuite: SetupTest")
}

func (suite *SlicesTestSuite) TearDownTest() {
	fmt.Println("SlicesTestSuite: TearDownTest")
}

func (suite *SlicesTestSuite) TestFirst() {
	type (
		args struct {
			sliceVal []any
		}

		want struct {
			element any
			ok      bool
		}

		testCase struct {
			name string
			args args
			want want
		}

		testElement struct {
			element1 string
			element2 int
			element3 float32
		}
	)

	testCases := []testCase{
		{
			name: "success positive int slice",
			args: args{
				sliceVal: []any{1, 2, 4, 5, 6, 7},
			},
			want: want{
				element: 1,
				ok:      true,
			},
		},
		{
			name: "success negative int slice",
			args: args{
				sliceVal: []any{-1, -2, -4, -5, -6, -7},
			},
			want: want{
				element: -1,
				ok:      true,
			},
		},
		{
			name: "success float slice",
			args: args{
				sliceVal: []any{1.4, 2.5, 4.1, 5.2, 6.8, 7.3},
			},
			want: want{
				element: 1.4,
				ok:      true,
			},
		},
		{
			name: "success string slice",
			args: args{
				sliceVal: []any{"a", "b", "c", "d", "e", "f"},
			},
			want: want{
				element: "a",
				ok:      true,
			},
		},
		{
			name: "success struct slice",
			args: args{
				sliceVal: []any{
					testElement{
						element1: "test1",
						element2: 1,
						element3: 1.2,
					},
					testElement{
						element1: "test2",
						element2: 2,
						element3: 2.3,
					},
					testElement{
						element1: "test3",
						element2: 3,
						element3: 3.4,
					},
					testElement{
						element1: "test4",
						element2: 4,
						element3: 4.5,
					},
					testElement{
						element1: "test6",
						element2: 6,
						element3: 6.7,
					},
					testElement{
						element1: "test7",
						element2: 7,
						element3: 7.8,
					},
				},
			},
			want: want{
				element: testElement{
					element1: "test1",
					element2: 1,
					element3: 1.2,
				},
				ok: true,
			},
		},
		{
			name: "failed empty slice",
			args: args{
				sliceVal: []any{},
			},
			want: want{
				element: nil,
				ok:      false,
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			el, ok := slices.First(tc.args.sliceVal)
			suite.Equal(tc.want.element, el)
			suite.Equal(tc.want.ok, ok)
		})
	}
}

func (suite *SlicesTestSuite) TestLast() {
	type (
		args struct {
			sliceVal []any
		}

		want struct {
			element any
			ok      bool
		}

		testCase struct {
			name string
			args args
			want want
		}

		testElement struct {
			element1 string
			element2 int
			element3 float32
		}
	)

	testCases := []testCase{
		{
			name: "success positive int slice",
			args: args{
				sliceVal: []any{1, 2, 4, 5, 6, 7},
			},
			want: want{
				element: 7,
				ok:      true,
			},
		},
		{
			name: "success negative int slice",
			args: args{
				sliceVal: []any{-1, -2, -4, -5, -6, -7},
			},
			want: want{
				element: -7,
				ok:      true,
			},
		},
		{
			name: "success float slice",
			args: args{
				sliceVal: []any{1.4, 2.5, 4.1, 5.2, 6.8, 7.3},
			},
			want: want{
				element: 7.3,
				ok:      true,
			},
		},
		{
			name: "success string slice",
			args: args{
				sliceVal: []any{"a", "b", "c", "d", "e", "f"},
			},
			want: want{
				element: "f",
				ok:      true,
			},
		},
		{
			name: "success struct slice",
			args: args{
				sliceVal: []any{
					testElement{
						element1: "test1",
						element2: 1,
						element3: 1.2,
					},
					testElement{
						element1: "test2",
						element2: 2,
						element3: 2.3,
					},
					testElement{
						element1: "test3",
						element2: 3,
						element3: 3.4,
					},
					testElement{
						element1: "test4",
						element2: 4,
						element3: 4.5,
					},
					testElement{
						element1: "test6",
						element2: 6,
						element3: 6.7,
					},
					testElement{
						element1: "test7",
						element2: 7,
						element3: 7.8,
					},
				},
			},
			want: want{
				element: testElement{
					element1: "test7",
					element2: 7,
					element3: 7.8,
				},
				ok: true,
			},
		},
		{
			name: "failed empty slice",
			args: args{
				sliceVal: []any{},
			},
			want: want{
				element: nil,
				ok:      false,
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			el, ok := slices.Last(tc.args.sliceVal)
			suite.Equal(tc.want.element, el)
			suite.Equal(tc.want.ok, ok)
		})
	}
}

func TestSlicesTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(SlicesTestSuite))
}
