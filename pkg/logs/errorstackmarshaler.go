package logs

import (
	"regexp"
	"runtime"
	"strings"

	"github.com/pkg/errors"
)

type stackTracer interface {
	StackTrace() errors.StackTrace
}

type errorWrapper interface {
	Unwrap() error
}

type stackTrace struct {
	Frames []frame     `json:"frames"`
	Cause  *stackTrace `json:"cause,omitempty"`
}

type frame struct {
	Source string `json:"source"`
	Line   int    `json:"line"`
	Func   string `json:"func"`
}

func removePackageName(funcName string) string {
	return funcName[strings.Index(funcName, ".")+1:]
}

func isGoSource(source string) bool {
	return strings.HasPrefix(source, goroot)
}

var (
	goPathRegex  = regexp.MustCompile("/go/(pkg/mod/)?")
	removeGoPath = func(source string) string {
		sourceParts := goPathRegex.Split(source, 2)
		return sourceParts[len(sourceParts)-1]
	}
)

func getErrorFrames(err error) *stackTrace {
	if err == nil {
		return nil
	}
	var stackErr stackTracer
	var ok bool
	if stackErr, ok = err.(stackTracer); !ok {
		return nil
	}

	stack := stackErr.StackTrace()
	result := &stackTrace{
		Frames: make([]frame, 0, len(stack)),
	}
	for _, stackFrame := range stack {
		pc := uintptr(stackFrame) - 1
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)

		if isGoSource(file) {
			continue
		}

		funcName := fn.Name()
		result.Frames = append(result.Frames, frame{
			Source: removeGoPath(file),
			Line:   line,
			Func:   removePackageName(funcName),
		})
	}
	return result
}

func getErrorStack(err error) *stackTrace {
	if err == nil {
		return nil
	}

	stack := getErrorFrames(err)

	if errWrapper, ok := err.(errorWrapper); ok {
		causeErr := errWrapper.Unwrap()
		causeStack := getErrorStack(causeErr)
		if stack == nil {
			return causeStack
		}
		stack.Cause = causeStack
	}

	return stack
}

func ErrorStackMarshaler(err error) interface{} {
	stack := getErrorStack(err)
	if stack == nil {
		return nil
	}
	return stack
}
