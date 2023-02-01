package errors

import (
	"context"
	stderrors "errors"
	"fmt"
	"os"
	"runtime"
	"strings"

	"bitbucket.org/moladinTech/go-lib-common/logger"
	"bitbucket.org/moladinTech/go-lib-common/registry"
)

type err struct {
	original      error
	wrapped       error
	stacktrace    string
	keyerr        error
	stack         []uintptr
	logCtx        string
	isNotify      bool
	isSuccessResp bool
}

func (e *err) StackTrace() []uintptr {
	return e.stack
}

func (e *err) WithNotify(ctx context.Context, rgs registry.IRegistry) *err {
	rgs.GetSentry().CaptureException(e)
	// send notif
	slackMessage := rgs.GetNotif().GetFormattedMessage(e.logCtx, ctx, e)
	errSlack := rgs.GetNotif().Send(ctx, slackMessage)
	if errSlack != nil {
		logger.Error(ctx, "Error sending notif to slack", errSlack)
	}

	e.isNotify = true
	return e
}

func (e *err) WithSuccessResp() *err {
	e.isSuccessResp = true
	return e
}

func (e *err) Error() string {
	var original string
	var wrapped string

	if e.wrapped == nil {
		return e.original.Error() + e.stacktrace
	}

	wrapped = e.wrapped.Error() + e.stacktrace

	if e.original != nil {
		if _, ok := e.original.(*err); !ok {
			original = "root cause : " + e.original.Error()
		} else {
			original = e.original.Error()
		}
	}

	return wrapped + ": " + original
}

func (e *err) Is(target error) bool {
	if e == target {
		return true
	}
	if stderrors.Is(e.original, target) {
		return true
	}
	return stderrors.Is(e.wrapped, target)
}

func (e *err) Unwrap() error {
	return e.original
}

// will return the root cause and return itself if the type of struct isnt err type
func RootErr(err_ error) error {
	if val, ok := err_.(*err); ok {
		return RootErr(val.original)
	}

	return err_
}

func Wrap(err_ error) *err {
	pc, file, no, _ := runtime.Caller(1)

	// Current working directory
	dir, _ := os.Getwd()
	splitDir := strings.Split(dir, "/")
	rootDir := splitDir[len(splitDir)-1]

	logCtx := generateLogCtx(pc, file, rootDir)

	retErr := err{
		original:   err_,
		wrapped:    nil,
		stacktrace: " -- At : " + fmt.Sprintf("%s:%d", file, no),
		keyerr:     GetErrKey(err_),
		stack:      []uintptr{pc},
		logCtx:     logCtx,
	}

	// since its error is in DFS, we want the message to be like bfs, to be bfs we implement when we call this func
	if val, ok := err_.(*err); ok {
		val.stacktrace = retErr.stacktrace + val.stacktrace
		retErr.stacktrace = ""

		retErr.stack = append(retErr.stack, val.stack...)
	}

	return &retErr

}

// WrapWithErr wrap new error into an existing error
func WrapWithErr(original error, wrapped error) *err {
	pc, file, no, _ := runtime.Caller(1)

	// Current working directory
	dir, _ := os.Getwd()
	splitDir := strings.Split(dir, "/")
	rootDir := splitDir[len(splitDir)-1]

	logCtx := generateLogCtx(pc, file, rootDir)

	return &err{
		original:   original,
		wrapped:    wrapped,
		keyerr:     GetErrKey(wrapped),
		stacktrace: " -- At : " + fmt.Sprintf("%s:%d", file, no),
		stack:      []uintptr{pc},
		logCtx:     logCtx,
	}
}

// get error as key to compare what the output response will be
func GetErrKey(err_ error) error {
	if val, ok := err_.(*err); ok {
		return val.keyerr
	}

	return err_
}

// generateLogCtx will generate the log context
func generateLogCtx(pc uintptr, file, rootDir string) string {
	var filePath, funcName string
	filePath = strings.Split(file, fmt.Sprintf("/%s/", rootDir))[1]
	filePath = strings.Split(filePath, ".")[0]
	filePath = strings.ReplaceAll(filePath, "/", ".")
	funcName = runtime.FuncForPC(pc).Name()
	i := strings.LastIndex(funcName, ".")
	funcName = funcName[i+1:]

	return fmt.Sprintf("%s.%s", filePath, funcName)
}
