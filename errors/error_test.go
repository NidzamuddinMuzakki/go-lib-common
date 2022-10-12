package errors

import (
	"context"
	"errors"
	stderrors "errors"
	"os"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TestErrorSuite struct {
	suite.Suite
	getPath func(lineNUm string) string
	path    string
}

func (t *TestErrorSuite) SetupSuite() {
	t.path, _ = os.Getwd()
	t.getPath = func(lineNum string) string {
		return " -- At : " + t.path + "/error_test.go:" + lineNum
	}
}

func (t *TestErrorSuite) TestUnwrap() {
	originalErr := errors.New("test")
	wrappedErr := stderrors.New("text")
	err := WrapWithErr(originalErr, wrappedErr)
	unwrappedErr := stderrors.Unwrap(err)
	assert.Equal(t.T(), "test", unwrappedErr.Error())
}

func (t *TestErrorSuite) TestIs() {
	testErr1 := errors.New("test err 1")
	testErr2 := errors.New("test err 2")
	testErr3 := errors.New("test err 3")
	testErr4 := errors.New("test err 4")
	testFail := errors.New("test fail")

	err := WrapWithErr(errors.New("origin"), testErr1)
	err = WrapWithErr(err, testErr2)
	err = WrapWithErr(err, testErr3)
	err = WrapWithErr(err, testErr4)

	assert.Equal(t.T(), true, stderrors.Is(err, testErr1))
	assert.Equal(t.T(), true, stderrors.Is(err, testErr2))
	assert.Equal(t.T(), true, stderrors.Is(err, testErr3))
	assert.Equal(t.T(), true, stderrors.Is(err, testErr4))
	assert.Equal(t.T(), true, stderrors.Is(err, err))
	assert.Equal(t.T(), false, stderrors.Is(err, testFail))
}

func (t *TestErrorSuite) TestWrappedError() {
	testErr1 := errors.New("test err 1")
	testErr2 := errors.New("test err 2")
	testErr3 := errors.New("test err 3")
	testErr4 := errors.New("test err 4")

	err := WrapWithErr(errors.New("origin"), testErr1)
	err = WrapWithErr(err, testErr2)
	err = WrapWithErr(err, testErr3)
	err = WrapWithErr(err, testErr4)
	expected := "test err 4" + t.getPath("65") + ": test err 3" + t.getPath("64") + ": test err 2" + t.getPath("63") + ": test err 1" + t.getPath("62") + ": root cause : origin"

	assert.Equal(t.T(), expected, err.Error())
}

func (t *TestErrorSuite) TestWrap() {
	testErr1 := context.DeadlineExceeded

	err := Wrap(testErr1)
	err = Wrap(err)
	expected := "context deadline exceeded" + t.getPath("75") + t.getPath("74")

	errKey := GetErrKey(err)

	assert.Equal(t.T(), expected, err.Error())
	assert.Equal(t.T(), testErr1, errKey)
}

func (t *TestErrorSuite) TestRootErr() {
	err1 := errors.New("TEST")
	err2 := errors.New("test 2")
	err := WrapWithErr(err1, err2)

	rootErr := RootErr(err)

	assert.Equal(t.T(), err1, rootErr)

}

func (t *TestErrorSuite) TestGetErrKey() {
	testErr1 := errors.New("test err 1")
	testErr2 := errors.New("test err 2")

	err := WrapWithErr(testErr1, testErr2)

	errKey := GetErrKey(err)

	assert.Equal(t.T(), testErr2, errKey)

	err = context.DeadlineExceeded

	errKey = GetErrKey(err)

	assert.Equal(t.T(), err, errKey)

}

func (t *TestErrorSuite) TestGetStackTrace() {
	err1 := errors.New("test err 1")
	err2 := errors.New("test err 2")

	errWrap := WrapWithErr(err1, err2)
	errWrap = AnotherWrapFunc(t.T(), errWrap)
	var stackPtr []uintptr
	if err_, ok := errWrap.(*err); ok {
		stackPtr = err_.StackTrace()
	}
	frames := runtime.CallersFrames(stackPtr)
	funcName := make([]string, 0)
	for {
		frame, ok := frames.Next()

		funcName = append(funcName, frame.Function)

		if !ok {
			break
		}
	}

	assert.Equal(t.T(), []string{"bitbucket.org/moladinTech/go-lib-common/errors.AnotherWrapFunc",
		"bitbucket.org/moladinTech/go-lib-common/errors.(*TestErrorSuite).TestGetStackTrace"}, funcName)
}

func AnotherWrapFunc(t *testing.T, err error) error {
	return Wrap(err)
}

func TestSuiteErrorPackage(t *testing.T) {
	suite.Run(t, new(TestErrorSuite))
}
