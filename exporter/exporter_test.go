//go:build unit
// +build unit

package exporter_test

import (
	"bytes"
	"context"
	"fmt"
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

func (e *ExporterSuite) SetupTest() {
	e.mock.sentry = commonSentryMock.NewISentry(e.T())
	e.Exporter = exporter.Newexporter(
		validator.New(),
		exporter.WithSentry(e.mock.sentry),
		exporter.WithExporterType(exporter.ExcelType),
		exporter.AddConverter(TestStruct{}, make(map[string]exporter.FuncConvert)))
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

func (e *ExporterSuite) TestExport() {
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

func TestExporter(t *testing.T) {
	suite.Run(t, new(ExporterSuite))
}
