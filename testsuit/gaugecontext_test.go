package testsuit

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestShouldGetStepsWithDescription(t *testing.T) {
	context := &GaugeContext{
		Steps: []Step{
			Step{
				Description: "I have an implementation",
				Impl:        func() {},
			},
		},
	}

	step, err := context.GetStepByDesc("I have an implementation")

	assert.NoError(t, err)
	assert.Equal(t, "I have an implementation", step.Description)
}

func TestShouldReturnErrorForUnImplementedStep(t *testing.T) {
	context := &GaugeContext{
		Steps: []Step{
			Step{
				Description: "I have an implementation",
				Impl:        func() {},
			},
		},
	}

	step, err := context.GetStepByDesc("I don't have an implementation")

	assert.Error(t, err)
	assert.Nil(t, step)
}

func TestShouldGetHooksOfGivenType(t *testing.T) {
	context := &GaugeContext{
		Hooks: []Hook{
			Hook{
				Type:     BEFORESUITE,
				Impl:     func() error { return nil },
				Tags:     []string{},
				Operator: AND,
			},
		},
	}

	hooks := context.GetHooks(BEFORESUITE, []string{"foo", "bar"})

	assert.Equal(t, 1, len(hooks))
}

func TestShouldGetHooksWithAnyOfTheseTags(t *testing.T) {
	context := &GaugeContext{
		Hooks: []Hook{
			Hook{
				Type:     BEFORESUITE,
				Impl:     func() error { return nil },
				Tags:     []string{"foo", "bar"},
				Operator: OR,
			},
			Hook{
				Type:     BEFORESUITE,
				Impl:     func() error { return nil },
				Tags:     []string{"notfoo", "bar"},
				Operator: OR,
			},
		},
	}

	hooks := context.GetHooks(BEFORESUITE, []string{"foo", "foobar"})

	assert.Equal(t, 1, len(hooks))
	assert.Contains(t, hooks[0].Tags, "foo")
}

func TestShouldNotGetHooksIfTagsDontMatch(t *testing.T) {
	context := &GaugeContext{
		Hooks: []Hook{
			Hook{
				Type:     BEFORESUITE,
				Impl:     func() error { return nil },
				Tags:     []string{"foo", "bar"},
				Operator: OR,
			},
			Hook{
				Type:     BEFORESUITE,
				Impl:     func() error { return nil },
				Tags:     []string{"notfoo", "bar"},
				Operator: OR,
			},
		},
	}

	hooks := context.GetHooks(BEFORESUITE, []string{"foobar"})

	assert.Equal(t, 0, len(hooks))
}
func TestShouldGetHooksWithBothTags(t *testing.T) {
	context := &GaugeContext{
		Hooks: []Hook{
			Hook{
				Type:     BEFORESUITE,
				Impl:     func() error { return nil },
				Tags:     []string{"foo", "bar"},
				Operator: AND,
			},
			Hook{
				Type:     BEFORESUITE,
				Impl:     func() error { return nil },
				Tags:     []string{"notfoo", "bar"},
				Operator: AND,
			},
		},
	}

	hooks := context.GetHooks(BEFORESUITE, []string{"foo", "bar"})

	assert.Equal(t, 1, len(hooks))
	assert.Contains(t, hooks[0].Tags, "foo")
	assert.Contains(t, hooks[0].Tags, "bar")
}
