package testsuit

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestShouldRunImplementation(t *testing.T) {
	called := false
	var calledWith interface{}
	step := Step{
		Description: "Test description",
		Impl: func(args ...interface{}) {
			calledWith = args
			called = true
		},
	}

	step.Execute(1, true, "foo")

	assert.True(t, called)
	assert.Contains(t, calledWith, 1)
	assert.Contains(t, calledWith, true)
	assert.Contains(t, calledWith, "foo")
}

func TestShouldReturnPassedMethodExecutionResult(t *testing.T) {
	var called bool
	step := Step{
		Description: "Test description",
		Impl: func(args ...interface{}) {
			doWork()
			called = true
		},
	}

	res := step.Execute("foo")

	assert.True(t, called)
	assert.False(t, res.GetFailed())
	assert.NotZero(t, res.GetExecutionTime())
	assert.Zero(t, res.GetErrorMessage())
	assert.Zero(t, res.GetStackTrace())
	assert.False(t, res.GetRecoverableError())
}

func TestShouldReturnFailedMethodExecutionResult(t *testing.T) {
	var called bool
	step := Step{
		Description: "Test description",
		Impl: func(args ...interface{}) {
			doWork()
			called = true
			var a []string
			fmt.Println(a[7])

		},
	}

	res := step.Execute("foo")

	assert.True(t, called)
	assert.True(t, res.GetFailed())
	assert.NotZero(t, res.GetExecutionTime())
	assert.Equal(t, "runtime error: index out of range [7] with length 0", res.GetErrorMessage())
	assert.NotZero(t, res.GetStackTrace())
	assert.False(t, res.GetRecoverableError())
}

func TestShouldReturnFailedButContinuableMethodExecutionResult(t *testing.T) {
	var called bool
	step := Step{
		Description: "Test description",
		Impl: func(args ...interface{}) {
			T.ContinueOnFailure()
			doWork()
			called = true
			var a []string
			fmt.Println(a[7])
		},
	}

	res := step.Execute("foo")

	assert.True(t, called)
	assert.True(t, res.GetFailed())
	assert.NotZero(t, res.GetExecutionTime())
	assert.Equal(t, "runtime error: index out of range [7] with length 0", res.GetErrorMessage())
	assert.NotZero(t, res.GetStackTrace())
	assert.True(t, res.GetRecoverableError())
}

func doWork() {
	// test execution time resolution is in milliseconds. So we need to wait at least that long for any
	// assertions on it to be meaningful.
	time.Sleep(1 * time.Millisecond)
}
