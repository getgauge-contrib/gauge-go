package testsuit

import (
	"testing"

	"github.com/stretchr/testify/assert"
	m "github.com/getgauge-contrib/gauge-go/gauge_messages"
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
				Impl:     func(*m.ExecutionInfo) {},
				Tags:     []string{},
				Operator: AND,
			},
			Hook{
				Type:     AFTERSUITE,
				Impl:     func(*m.ExecutionInfo) {},
				Tags:     []string{},
				Operator: AND,
			},
		},
	}

	hooks := context.GetHooks(BEFORESUITE, []string{"foo", "bar"})

	assert.Equal(t, 1, len(hooks))
}

// Suite hooks should be not filtered by tags

func TestBeforeSuiteShouldGetHooksIfTagsDontMatch(t *testing.T) {
	context := &GaugeContext{
		Hooks: []Hook{
			Hook{
				Type:     BEFORESUITE,
				Impl:     func(*m.ExecutionInfo) {},
				Tags:     []string{"foo", "bar"},
				Operator: OR,
			},
			Hook{
				Type:     BEFORESUITE,
				Impl:     func(*m.ExecutionInfo) {},
				Tags:     []string{"notfoo", "bar"},
				Operator: OR,
			},
		},
	}

	hooks := context.GetHooks(BEFORESUITE, []string{"foobar"})

	assert.Equal(t, 2, len(hooks))
	assert.NotContains(t, hooks[0].Tags, "foobar")
	assert.NotContains(t, hooks[1].Tags, "foobar")
}

func TestAfterSuiteShouldGetHooksIfTagsDontMatch(t *testing.T) {
	context := &GaugeContext{
		Hooks: []Hook{
			Hook{
				Type:     AFTERSUITE,
				Impl:     func(*m.ExecutionInfo) {},
				Tags:     []string{"foo", "bar"},
				Operator: OR,
			},
			Hook{
				Type:     AFTERSUITE,
				Impl:     func(*m.ExecutionInfo) {},
				Tags:     []string{"notfoo", "bar"},
				Operator: OR,
			},
		},
	}

	hooks := context.GetHooks(AFTERSUITE, []string{"foobar"})

	assert.Equal(t, 2, len(hooks))
	assert.NotContains(t, hooks[0].Tags, "foobar")
	assert.NotContains(t, hooks[1].Tags, "foobar")
}

// Test hook filtering by tags

func TestShouldGetHooksWithAnyOfTheseTags(t *testing.T) {
	context := &GaugeContext{
		Hooks: []Hook{
			Hook{
				Type:     BEFORESPEC,
				Impl:     func(*m.ExecutionInfo) {},
				Tags:     []string{"foo", "bar"},
				Operator: OR,
			},
			Hook{
				Type:     BEFORESPEC,
				Impl:     func(*m.ExecutionInfo) {},
				Tags:     []string{"notfoo", "bar"},
				Operator: OR,
			},
		},
	}

	hooks := context.GetHooks(BEFORESPEC, []string{"foo", "foobar"})

	assert.Equal(t, 1, len(hooks))
	assert.Contains(t, hooks[0].Tags, "foo")
}

func TestShouldNotGetHooksIfTagsDontMatch(t *testing.T) {
	context := &GaugeContext{
		Hooks: []Hook{
			Hook{
				Type:     BEFORESPEC,
				Impl:     func(*m.ExecutionInfo) {},
				Tags:     []string{"foo", "bar"},
				Operator: OR,
			},
			Hook{
				Type:     BEFORESPEC,
				Impl:     func(*m.ExecutionInfo) {},
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
				Type:     BEFORESPEC,
				Impl:     func(*m.ExecutionInfo) {},
				Tags:     []string{"foo", "bar"},
				Operator: AND,
			},
			Hook{
				Type:     BEFORESPEC,
				Impl:     func(*m.ExecutionInfo) {},
				Tags:     []string{"notfoo", "bar"},
				Operator: AND,
			},
		},
	}

	hooks := context.GetHooks(BEFORESPEC, []string{"foo", "bar"})

	assert.Equal(t, 1, len(hooks))
	assert.Contains(t, hooks[0].Tags, "foo")
	assert.Contains(t, hooks[0].Tags, "bar")
}
