package testsuit

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestShouldRunImplementation(t *testing.T){
	called := false
	step := Step{
		Description: "Test description",
		Impl: func(){ called = true },
	}

	step.Execute()
	
	assert.True(t, called)
}
