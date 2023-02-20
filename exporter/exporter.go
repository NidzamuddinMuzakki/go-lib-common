//go:generate mockery --name=Exporter
package exporter

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"

	"reflect"
	"time"

	commonError "bitbucket.org/moladinTech/go-lib-common/errors"
	"bitbucket.org/moladinTech/go-lib-common/sentry"
	commonTime "bitbucket.org/moladinTech/go-lib-common/time"
	commonValidator "bitbucket.org/moladinTech/go-lib-common/validator"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/go-playground/validator/v10"
)

const (
	ExcelType = "EXCEL"
	CSVType   = "CSV"

	TagExporter = "exporter"

	CSVDelimiter = ","

	DefaultTimeFormat = "2006-01-02 15:04:05"
)

type ExportType string

type FuncConvert func(interface{}) string

type Exporter interface {
	Export(context.Context, interface{}) (ResultExport, error)
}

type MapFuncConvert map[string]FuncConvert

func (m *MapFuncConvert) Add(tagName string, func_ FuncConvert) {
	(*m)[tagName] = func_
}

type ResultExport struct {
	ExcelObj *excelize.File
	ExcelRaw []byte
	CSVRaw   []byte
}

type exporter struct {
	converter    map[string]map[string]FuncConvert
	exporterType ExportType `validate:"required"`
	Exporter     Exporter
	Sentry       sentry.ISentry `validate:"required"`
}

type exporterExcel map[string]map[string]FuncConvert

func (e *exporterExcel) Export(ctx context.Context, v interface{}) (ResultExport, error) {
	switch reflect.TypeOf(v).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(v)

		// create new excel object
		file := excelize.NewFile()

		//set sheet name with name of struct
		structElem := s.Type().Elem()
		if structElem.Kind() == reflect.Pointer {
			structElem = structElem.Elem()
		}
		structName := structElem.Name()
		file.SetSheetName(file.GetSheetName(1), structName)

		// set excel header from struct
		e.GetHeaderFromStruct(reflect.New(structElem).Elem(), file, -1, 1, structName)

		//get converter map function
		convFunc := (*e)[structName]

		for i := 0; i < s.Len(); i++ {
			col := -1
			t := s.Index(i)
			if t.Type().Kind() == reflect.Pointer {
				t = t.Elem()
			}
			e.SetValueExcel(t, file, col, i+2, structName, convFunc)
		}

		return e.ReturnResult(file)
	}

	return ResultExport{}, commonError.ErrorExporterNotSupportedType
}

func (e *exporterExcel) ReturnResult(file *excelize.File) (ResultExport, error) {
	buffer, err := file.WriteToBuffer()
	if err != nil {
		return ResultExport{}, err
	}

	return ResultExport{
		ExcelObj: file,
		ExcelRaw: buffer.Bytes(),
	}, nil
}

func (e *exporterExcel) GetHeaderFromStruct(t reflect.Value, exc *excelize.File, col int, row int, sheetName string) {
	for i := 0; i < t.NumField(); i++ {
		field := t.Type().Field(i)
		tag := field.Tag.Get(TagExporter)

		if len(tag) > 0 && tag != "-" {
			exc.SetCellValue(sheetName, incrColumnAndRow(&col, row), tag)
		}
	}
}

func (e *exporterExcel) SetValueExcel(t reflect.Value, exc *excelize.File, col int, row int, sheetName string, convFunc map[string]FuncConvert) {
	for i := 0; i < t.NumField(); i++ {
		field := t.Type().Field(i)
		tag := field.Tag.Get(TagExporter)

		if len(tag) > 0 && tag != "-" {

			type_ := t.Field(i)              // get field[idx] on the struct
			if type_.Kind() == reflect.Ptr { // if pointer then get the pointed element
				if !type_.IsNil() {
					type_ = type_.Elem()
				}
			}

			interface_ := type_.Interface()

			// call preprocessing function
			if func_, ok := convFunc[tag]; ok {
				exc.SetCellValue(sheetName, incrColumnAndRow(&col, row), func_(interface_))
				continue
			}

			// for type casted

			switch interface_ := interface_.(type) {
			case commonTime.DateTime:
				exc.SetCellValue(sheetName, incrColumnAndRow(&col, row), time.Time(interface_))
			default:
				exc.SetCellValue(sheetName, incrColumnAndRow(&col, row), interface_)
			}
		}
	}
}

func incrColumnAndRow(col *int, row int) string {
	*col++
	return fmt.Sprintf("%s%d", excelize.ToAlphaString(*col), row)
}

type exporterCSV map[string]map[string]FuncConvert

func (e *exporterCSV) Export(ctx context.Context, v interface{}) (ResultExport, error) {
	switch reflect.TypeOf(v).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(v)

		// create new csv writer
		var buffer bytes.Buffer
		writer := csv.NewWriter(&buffer)

		structElem := s.Type().Elem()
		if structElem.Kind() == reflect.Pointer {
			structElem = structElem.Elem()
		}
		structName := structElem.Name()

		// set header
		header := e.ExtractHeader(reflect.New(structElem).Elem())
		lengthRow := len(header)
		writer.Write(header)

		//get converter map function
		convFunc := (*e)[structName]

		for i := 0; i < s.Len(); i++ {
			t := s.Index(i)
			if t.Type().Kind() == reflect.Pointer {
				t = t.Elem()
			}
			writer.Write(e.ExtractRow(t, lengthRow, convFunc))
		}

		return e.ReturnResult(writer, &buffer)
	}

	return ResultExport{}, commonError.ErrorExporterNotSupportedType
}

func (e *exporterCSV) ExtractRow(t reflect.Value, length int, convFunc map[string]FuncConvert) []string {
	row := make([]string, 0, length)
	for i := 0; i < t.NumField(); i++ {
		field := t.Type().Field(i)
		tag := field.Tag.Get(TagExporter)

		if len(tag) > 0 && tag != "-" {

			type_ := t.Field(i)              // get field[idx] on the struct
			if type_.Kind() == reflect.Ptr { // if pointer then get the pointed element
				if !type_.IsNil() {
					type_ = type_.Elem()
				}
			}

			interface_ := type_.Interface()

			// call preprocessing function
			if func_, ok := convFunc[tag]; ok {
				row = append(row, func_(interface_))
				continue
			}

			// for type casted
			switch interface_ := interface_.(type) {
			case commonTime.DateTime:
				row = append(row, time.Time(interface_).Format(DefaultTimeFormat))
			case string:
				row = append(row, fmt.Sprintf("\"%s\"", interface_))
			case time.Time:
				row = append(row, interface_.Format(DefaultTimeFormat))
			default:
				row = append(row, fmt.Sprint(interface_))
			}
		}
	}
	return row
}

func (e *exporterCSV) ExtractHeader(t reflect.Value) []string {
	row := make([]string, 0)
	for i := 0; i < t.NumField(); i++ {
		field := t.Type().Field(i)
		tag := field.Tag.Get(TagExporter)

		if len(tag) > 0 && tag != "-" {
			row = append(row, tag)
		}
	}

	return row
}

func (e *exporterCSV) ReturnResult(writer *csv.Writer, buffer *bytes.Buffer) (ResultExport, error) {
	writer.Flush()
	if err := writer.Error(); err != nil {
		return ResultExport{}, err
	}
	return ResultExport{CSVRaw: buffer.Bytes()}, nil
}

type Option func(*exporter)

func AddConverter(strct interface{}, funcConv map[string]FuncConvert) Option {
	return func(exp *exporter) {
		if reflect.ValueOf(strct).Kind() != reflect.Struct {
			panic(`the registered inteface is not a struct`)
		}

		nameStruct := reflect.TypeOf(strct).Name()

		exp.converter[nameStruct] = funcConv
	}
}

func WithExporterType(expType ExportType) Option {
	return func(exp *exporter) {
		exp.exporterType = expType
	}
}

func WithSentry(sentry sentry.ISentry) Option {
	return func(exp *exporter) {
		exp.Sentry = sentry
	}
}

func Newexporter(
	validator *validator.Validate,
	opt ...Option,
) Exporter {
	exporter := &exporter{
		converter: make(map[string]map[string]FuncConvert),
	}

	for _, option := range opt {
		option(exporter)
	}

	err := validator.Struct(exporter)
	if err != nil {
		panic(commonValidator.ToErrResponse(err))
	}

	switch exporter.exporterType {
	case ExcelType:
		ex := exporterExcel(exporter.converter)
		exporter.Exporter = &ex
	case CSVType:
		ex := exporterCSV(exporter.converter)
		exporter.Exporter = &ex
	default:
		panic(`invalid exporter type`)
	}

	return exporter
}

func (e *exporter) Export(ctx context.Context, v interface{}) (ResultExport, error) {
	logCtx := `common.exporter.Export.` + e.exporterType
	span := e.Sentry.StartSpan(ctx, string(logCtx))
	defer e.Sentry.Finish(span)
	ctx = e.Sentry.SpanContext(*span)

	return e.Exporter.Export(ctx, v)
}
