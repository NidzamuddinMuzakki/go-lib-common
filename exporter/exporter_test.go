//go:build unit
// +build unit

package exporter_test

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"reflect"
	"testing"
	"time"

	"bitbucket.org/moladinTech/go-lib-common/cast"
	commonError "bitbucket.org/moladinTech/go-lib-common/errors"
	"bitbucket.org/moladinTech/go-lib-common/exporter"
	commonSentryMock "bitbucket.org/moladinTech/go-lib-common/sentry/mocks"
	commonTime "bitbucket.org/moladinTech/go-lib-common/time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/getsentry/sentry-go"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type mockStruct struct {
	sentry *commonSentryMock.ISentry
}

type ExporterSuite struct {
	suite.Suite
	mock     mockStruct
	Exporter exporter.Exporter
}

type TestStruct struct {
	Name string               `exporter:"Name"`
	When *commonTime.DateTime `exporter:"When"`
}

type TestStruct2 struct {
	Name string              `exporter:"Name"`
	When commonTime.DateTime `exporter:"When"`
}

type TestEmbeddedStruct struct {
	Word string `exporter:"Word"`
	TestStruct
}

type TestEmbeddedStructExcel struct {
	TestStruct2
}

func (t TestStruct) GenerateConverter() exporter.MapFuncConvert {
	mapFuncConvert := make(exporter.MapFuncConvert)
	mapFuncConvert.Add("When", func(v interface{}) string {
		time_ := (time.Time)(v.(commonTime.DateTime))
		return time_.Format(time.Kitchen)
	})

	return mapFuncConvert
}

func (e *ExporterSuite) SetupTest() {
	e.mock.sentry = commonSentryMock.NewISentry(e.T())
	e.Exporter = exporter.Newexporter(
		validator.New(),
		exporter.WithSentry(e.mock.sentry),
		exporter.WithExporterType(exporter.ExcelType),
		exporter.AddConverter(TestStruct{}, make(map[string]exporter.FuncConvert)))
}

func (e *ExporterSuite) SetupTestCSV() {
	e.mock.sentry = commonSentryMock.NewISentry(e.T())
	e.Exporter = exporter.Newexporter(
		validator.New(),
		exporter.WithSentry(e.mock.sentry),
		exporter.WithExporterType(exporter.CSVType),
		exporter.AddConverter(TestStruct{}, make(map[string]exporter.FuncConvert)))
}

func (e *ExporterSuite) SetupTestCSVUsingConverter() {
	e.mock.sentry = commonSentryMock.NewISentry(e.T())
	e.Exporter = exporter.Newexporter(
		validator.New(),
		exporter.WithSentry(e.mock.sentry),
		exporter.WithExporterType(exporter.CSVType),
		exporter.AddConverter(TestStruct{}, TestStruct{}.GenerateConverter()))
}

func (e *ExporterSuite) TearDownTest() {
	e.mock.sentry.AssertExpectations(e.T())
}

func mockingSentryConf(mock_ *commonSentryMock.ISentry, ctx interface{}, logCtx string, returnContext context.Context, callSpan bool, callFinish bool) {

	var newSpan sentry.Span
	mock_.On("StartSpan", mock.Anything, logCtx).Run(func(args mock.Arguments) {
		ctx := args.Get(0).(context.Context)

		newSpan = *sentry.StartSpan(ctx, logCtx)
		if callSpan {
			mock_.On("SpanContext", newSpan).Return(newSpan.Context()).Once()
		}
		if callFinish {
			mock_.On("Finish", &newSpan).Once()

		}
	}).Return(&newSpan).Once()

}

func (e *ExporterSuite) TestExportExcel() {
	type (
		args struct {
			ctx context.Context
			v   interface{}
		}

		funcMock func(a args)

		output struct {
			err  error
			file string
		}

		testCase struct {
			args
			funcMock
			output
			description string
			no          int
		}
	)

	testCases := []testCase{
		{
			args: args{
				ctx: context.Background(),
				v: []TestStruct{
					{
						Name: "Test Name",
						When: cast.NewPointer(commonTime.DateTime(time.Date(2022, 2, 30, 9, 0, 0, 0, time.UTC))),
					},
				},
			},
			funcMock: func(a args) {
				mockingSentryConf(e.mock.sentry, a.ctx, "common.exporter.Export."+"EXCEL", context.Background(), true, true)
			},
			output: output{
				err:  nil,
				file: "TestStruct.xlsx",
			},
			description: "should success",
			no:          1,
		},
		{
			args: args{
				ctx: context.Background(),
				v: []*TestStruct{
					{
						Name: "Test Name",
						When: cast.NewPointer(commonTime.DateTime(time.Date(2022, 2, 30, 9, 0, 0, 0, time.UTC))),
					},
				},
			},
			funcMock: func(a args) {
				mockingSentryConf(e.mock.sentry, a.ctx, "common.exporter.Export."+"EXCEL", context.Background(), true, true)
			},
			output: output{
				err:  nil,
				file: "TestStruct.xlsx",
			},
			description: "should success with pointer struct",
			no:          2,
		},
	}

	for _, testCase_ := range testCases {
		e.Run(fmt.Sprintf("%d %s", testCase_.no, testCase_.description), func() {
			testCase_.funcMock(testCase_.args)

			res, err := e.Exporter.Export(testCase_.ctx, testCase_.v)

			assert.Equal(e.T(), testCase_.err, commonError.GetErrKey(err))

			excelRes, err := excelize.OpenReader(bytes.NewReader(res.ExcelRaw))
			if err != nil {
				e.T().Fatal(err)
			}

			ExcelExp, err := excelize.OpenFile(testCase_.file)
			if err != nil {
				e.T().Fatal(err)
			}

			assert.Equal(e.T(), ExcelExp.XLSX, excelRes.XLSX)

		})
	}
}
func (e *ExporterSuite) TestExportCSV() {
	type (
		args struct {
			ctx context.Context
			v   interface{}
		}

		funcMock func(a args)

		output struct {
			err  error
			file [][]string
		}

		testCase struct {
			args
			funcMock
			output
			description string
			no          int
		}
	)

	testCases := []testCase{
		{
			args: args{
				ctx: context.Background(),
				v: []TestStruct{
					{
						Name: "Test Name",
						When: cast.NewPointer(commonTime.DateTime(time.Date(2022, 3, 28, 9, 0, 0, 0, time.UTC))),
					},
				},
			},
			funcMock: func(a args) {
				mockingSentryConf(e.mock.sentry, a.ctx, "common.exporter.Export."+"CSV", context.Background(), true, true)
			},
			output: output{
				err: nil,
				file: [][]string{
					{"Name", "When"},
					{"Test Name", "2022-03-28 09:00:00"},
				},
			},
			description: "should success",
			no:          1,
		},
		{
			args: args{
				ctx: context.Background(),
				v: []*TestStruct{
					{
						Name: "Test Name",
						When: cast.NewPointer(commonTime.DateTime(time.Date(2022, 3, 28, 9, 0, 0, 0, time.UTC))),
					},
				},
			},
			funcMock: func(a args) {
				mockingSentryConf(e.mock.sentry, a.ctx, "common.exporter.Export."+"CSV", context.Background(), true, true)
			},
			output: output{
				err: nil,
				file: [][]string{
					{"Name", "When"},
					{"Test Name", "2022-03-28 09:00:00"},
				},
			},
			description: "should success with pointer struct",
			no:          2,
		},
		{
			args: args{
				ctx: context.Background(),
				v: []*TestStruct{
					{
						Name: `Test Name with Quote ""`,
						When: cast.NewPointer(commonTime.DateTime(time.Date(2022, 3, 28, 9, 0, 0, 0, time.UTC))),
					},
				},
			},
			funcMock: func(a args) {
				mockingSentryConf(e.mock.sentry, a.ctx, "common.exporter.Export."+"CSV", context.Background(), true, true)
			},
			output: output{
				err: nil,
				file: [][]string{
					{"Name", "When"},
					{`Test Name with Quote ""`, "2022-03-28 09:00:00"},
				},
			},
			description: "should success with string quote",
			no:          3,
		},
	}

	for _, testCase_ := range testCases {
		e.Run(fmt.Sprintf("%d %s", testCase_.no, testCase_.description), func() {
			e.SetupTestCSV()
			testCase_.funcMock(testCase_.args)

			res, err := e.Exporter.Export(testCase_.ctx, testCase_.v)

			assert.Equal(e.T(), testCase_.err, commonError.GetErrKey(err))

			// read csv values using csv.Reader
			csvReader := csv.NewReader(bytes.NewReader(res.CSVRaw))
			data, err := csvReader.ReadAll()
			if err != nil {
				log.Fatal(err)
			}

			e.EqualValues(testCase_.output.file, data)

		})
	}
}

func (e *ExporterSuite) TestExportUsingConverter() {
	type (
		args struct {
			ctx context.Context
			v   interface{}
		}

		funcMock func(a args)

		output struct {
			err  error
			file [][]string
		}

		testCase struct {
			args
			funcMock
			output
			description string
			no          int
		}
	)

	testCases := []testCase{
		{
			args: args{
				ctx: context.Background(),
				v: []TestStruct{
					{
						Name: "Test Name",
						When: cast.NewPointer(commonTime.DateTime(time.Date(2022, 3, 28, 9, 0, 0, 0, time.UTC))),
					},
				},
			},
			funcMock: func(a args) {
				mockingSentryConf(e.mock.sentry, a.ctx, "common.exporter.Export."+"CSV", context.Background(), true, true)
			},
			output: output{
				err: nil,
				file: [][]string{
					{"Name", "When"},
					{"Test Name", "9:00AM"},
				},
			},
			description: "should success",
			no:          1,
		},
		{
			args: args{
				ctx: context.Background(),
				v: []*TestStruct{
					{
						Name: "Test Name",
						When: cast.NewPointer(commonTime.DateTime(time.Date(2022, 3, 28, 9, 0, 0, 0, time.UTC))),
					},
				},
			},
			funcMock: func(a args) {
				mockingSentryConf(e.mock.sentry, a.ctx, "common.exporter.Export."+"CSV", context.Background(), true, true)
			},
			output: output{
				err: nil,
				file: [][]string{
					{"Name", "When"},
					{"Test Name", "9:00AM"},
				},
			},
			description: "should success with pointer struct",
			no:          2,
		},
	}

	for _, testCase_ := range testCases {
		e.Run(fmt.Sprintf("%d %s", testCase_.no, testCase_.description), func() {
			e.SetupTestCSVUsingConverter()
			testCase_.funcMock(testCase_.args)

			res, err := e.Exporter.Export(testCase_.ctx, testCase_.v)

			assert.Equal(e.T(), testCase_.err, commonError.GetErrKey(err))

			// read csv values using csv.Reader
			csvReader := csv.NewReader(bytes.NewReader(res.CSVRaw))
			data, err := csvReader.ReadAll()
			if err != nil {
				log.Fatal(err)
			}

			e.EqualValues(testCase_.output.file, data)

		})
	}
}
func (e *ExporterSuite) TestExportNestedStruct() {
	type (
		args struct {
			ctx context.Context
			v   interface{}
		}

		funcMock func(a args)

		output struct {
			err  error
			file [][]string
		}

		testCase struct {
			args
			funcMock
			output
			description string
			no          int
		}
	)

	testCases := []testCase{
		{
			args: args{
				ctx: context.Background(),
				v: []TestEmbeddedStruct{
					{
						TestStruct: TestStruct{
							Name: "Test Name",
							When: cast.NewPointer(commonTime.DateTime(time.Date(2022, 3, 28, 9, 0, 0, 0, time.UTC))),
						},
						Word: "Test Word",
					},
				},
			},
			funcMock: func(a args) {
				mockingSentryConf(e.mock.sentry, a.ctx, "common.exporter.Export."+"CSV", context.Background(), true, true)
			},
			output: output{
				err: nil,
				file: [][]string{
					{"Word", "Name", "When"},
					{"Test Word", "Test Name", "2022-03-28 09:00:00"},
				},
			},
			description: "should success",
			no:          1,
		},
		{
			args: args{
				ctx: context.Background(),
				v: []*TestEmbeddedStruct{
					{
						TestStruct: TestStruct{
							Name: "Test Name",
							When: cast.NewPointer(commonTime.DateTime(time.Date(2022, 3, 28, 9, 0, 0, 0, time.UTC))),
						},
						Word: `Test Word`,
					},
				},
			},
			funcMock: func(a args) {
				mockingSentryConf(e.mock.sentry, a.ctx, "common.exporter.Export."+"CSV", context.Background(), true, true)
			},
			output: output{
				err: nil,
				file: [][]string{
					{"Word", "Name", "When"},
					{"Test Word", "Test Name", "2022-03-28 09:00:00"},
				},
			},
			description: "should success with pointer struct",
			no:          2,
		},
	}

	for _, testCase_ := range testCases {
		e.Run(fmt.Sprintf("%d %s", testCase_.no, testCase_.description), func() {
			e.SetupTestCSVUsingConverter()
			testCase_.funcMock(testCase_.args)

			res, err := e.Exporter.Export(testCase_.ctx, testCase_.v)

			assert.Equal(e.T(), testCase_.err, commonError.GetErrKey(err))

			// read csv values using csv.Reader
			csvReader := csv.NewReader(bytes.NewReader(res.CSVRaw))
			data, err := csvReader.ReadAll()
			if err != nil {
				log.Fatal(err)
			}

			e.EqualValues(testCase_.output.file, data)
			e.TearDownTest()

		})
	}
}

func (e *ExporterSuite) TestEmbeddedStructExportExcel() {
	type (
		args struct {
			ctx context.Context
			v   interface{}
		}

		funcMock func(a args)

		output struct {
			err  error
			file string
			raw  [][]string
		}

		testCase struct {
			args
			funcMock
			output
			description string
			no          int
		}
	)

	testCases := []testCase{
		{
			args: args{
				ctx: context.Background(),
				v: []TestEmbeddedStructExcel{
					{
						TestStruct2: TestStruct2{
							Name: "Test Name",
							When: commonTime.DateTime(time.Date(2022, 2, 30, 9, 0, 0, 0, time.UTC)),
						},
					},
				},
			},
			funcMock: func(a args) {
				mockingSentryConf(e.mock.sentry, a.ctx, "common.exporter.Export."+"EXCEL", context.Background(), true, true)
			},
			output: output{
				err:  nil,
				file: "TestStruct.xlsx",
				raw: [][]string{
					{"Name", "When"},
					{"Test Name", "3/2/22 09:00"},
				},
			},
			description: "should success",
			no:          1,
		},
		{
			args: args{
				ctx: context.Background(),
				v: []*TestEmbeddedStructExcel{
					{
						TestStruct2: TestStruct2{
							Name: "Test Name",
							When: commonTime.DateTime(time.Date(2022, 2, 30, 9, 0, 0, 0, time.UTC)),
						},
					},
				},
			},
			funcMock: func(a args) {
				mockingSentryConf(e.mock.sentry, a.ctx, "common.exporter.Export."+"EXCEL", context.Background(), true, true)
			},
			output: output{
				err:  nil,
				file: "TestStruct.xlsx",
				raw: [][]string{
					{"Name", "When"},
					{"Test Name", "3/2/22 09:00"},
				},
			},
			description: "should success with pointer struct",
			no:          2,
		},
	}

	for _, testCase_ := range testCases {
		e.Run(fmt.Sprintf("%d %s", testCase_.no, testCase_.description), func() {
			e.SetupTest()
			testCase_.funcMock(testCase_.args)

			res, err := e.Exporter.Export(testCase_.ctx, testCase_.v)

			assert.Equal(e.T(), testCase_.err, commonError.GetErrKey(err))

			excelRes, err := excelize.OpenReader(bytes.NewReader(res.ExcelRaw))
			if err != nil {
				e.T().Fatal(err)
			}

			e.T().Log("Name struct", reflect.ValueOf(TestEmbeddedStructExcel{}))

			rows := excelRes.GetRows(reflect.ValueOf(TestEmbeddedStructExcel{}).Type().Name())

			e.EqualValues(testCase_.output.raw, rows)

		})
	}
}

func TestExporter(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ExporterSuite))
}
