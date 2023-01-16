package errors

import (
	stderrors "errors"
	"fmt"
	"runtime"
)

type err struct {
	original   error
	wrapped    error
	stacktrace string
	keyerr     error
	stack      []uintptr
	logCtx     string
}

func Wrap(err_ error) error {
	pc, file, no, _ := runtime.Caller(1)
	retErr := err{
		original:   err_,
		wrapped:    nil,
		stacktrace: " -- At : " + fmt.Sprintf("%s:%d", file, no),
		keyerr:     GetErrKey(err_),
		stack:      []uintptr{pc},
	}

	// since its error is in DFS, we want the message to be like bfs, to be bfs we implement when we call this func
	if val, ok := err_.(*err); ok {
		val.stacktrace = retErr.stacktrace + val.stacktrace
		retErr.stacktrace = ""

		retErr.stack = append(retErr.stack, val.stack...)
	}

	return &retErr

}

func (e *err) StackTrace() []uintptr {
	return e.stack
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

// Wrap new error into an existing error
func WrapWithErr(original error, wrapped error) error {
	pc, file, no, _ := runtime.Caller(1)
	return &err{
		original:   original,
		wrapped:    wrapped,
		keyerr:     GetErrKey(wrapped),
		stacktrace: " -- At : " + fmt.Sprintf("%s:%d", file, no),
		stack:      []uintptr{pc},
	}
}

// get error as key to compare what the output response will be
func GetErrKey(err_ error) error {
	if val, ok := err_.(*err); ok {
		return val.keyerr
	}

	return err_
}

// SetLogCtx will set the logCtx of the error
func SetLogCtx(err_ error, logCtx string) error {
	if val, ok := err_.(*err); ok {
		val.logCtx = logCtx
	}

	return err_
}

// getLogCtx will return the logCtx of the error
func getLogCtx(err_ error) string {
	if val, ok := err_.(*err); ok {
		return val.logCtx
	}

	return ""
}
