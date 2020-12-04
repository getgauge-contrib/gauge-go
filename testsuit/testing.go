package testsuit

import (
	"fmt"
	"runtime/debug"
	"strings"
	"testing"
)

// T is a wrapper around *testing.T
var T *testingT

type testingT struct {
	testing.T
	errors            []testError
	continueOnFailure bool
}

type testError struct {
	errMsg     string
	stacktrace string
}

func (t *testingT) getErrors() string {
	var errs []string
	for _, e := range t.errors {
		errs = append(errs, e.errMsg)
	}
	return strings.Join(errs, "\n")
}

func (t *testingT) getStacktraces() string {
	var stacktraces []string
	for _, e := range t.errors {
		stacktraces = append(stacktraces, e.stacktrace)
	}
	return strings.Join(stacktraces, "\n\n")
}

func (t *testingT) getContinueOnFailure() bool {
	return t.continueOnFailure
}

// Fail fails the step execution with the given error
func (t *testingT) Fail(err error) {
	panic(err)
}

// Any errors added through the Errorf method will now let the step continue to the next one.
func (t *testingT) ContinueOnFailure() {
	t.continueOnFailure = true
}

// Errorf records the error given, but step execution continues. However, the step is marked as failure.
func (t *testingT) Errorf(format string, args ...interface{}) {
	t.errors = append(t.errors, testError{errMsg: fmt.Sprintf(format, args...), stacktrace: string(debug.Stack())})
}
