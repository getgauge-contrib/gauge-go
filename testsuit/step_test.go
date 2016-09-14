package testsuit

import (
	"testing"

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
