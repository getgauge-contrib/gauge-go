package messageprocessors

import (
	"errors"
	"testing"

	m "github.com/getgauge-contrib/gauge-go/gauge_messages"
	t "github.com/getgauge-contrib/gauge-go/testsuit"
	"github.com/stretchr/testify/assert"
)

func TestShouldReturnExecutionStatusResponseWithSameIdForScenarioExecutionEnding(tst *testing.T) {
	msgId := int64(12345)
	context := &t.GaugeContext{
		Steps: make([]t.Step, 0),
	}

	msg := &m.Message{
		MessageType: m.Message_ScenarioExecutionEnding,
		MessageId:   msgId,
	}

	p := ScenarioExecutionEndingProcessor{}

	result := p.Process(msg, context)

	assert.Equal(tst, result.MessageType, m.Message_ExecutionStatusResponse)
	assert.Equal(tst, result.MessageId, msgId)
}

func TestExecutesHooksForScenarioTagsForScenarioExecutionEnding(tst *testing.T) {
	called1 := false
	called2 := false
	context := &t.GaugeContext{
		Hooks: []t.Hook{
			t.Hook{
				Type: t.AFTERSCENARIO,
				Impl: func(*m.ExecutionInfo) {
					called1 = true
				},
				Tags:     []string{"foo", "bar"},
				Operator: t.AND,
			},
			t.Hook{
				Type: t.AFTERSCENARIO,
				Impl: func(*m.ExecutionInfo) {
					called2 = true
				},
				Tags:     []string{"notfoo", "bar"},
				Operator: t.OR,
			},
		},
	}
	msgId := int64(12345)
	msg := &m.Message{
		MessageType: m.Message_ScenarioExecutionEnding,
		MessageId:   msgId,
		ScenarioExecutionEndingRequest: &m.ScenarioExecutionEndingRequest{
			CurrentExecutionInfo: &m.ExecutionInfo{
				CurrentScenario: &m.ScenarioInfo{
					Tags: []string{"foo", "bar"},
				},
			},
		},
	}

	p := ScenarioExecutionEndingProcessor{}

	result := p.Process(msg, context)

	assert.Equal(tst, result.MessageType, m.Message_ExecutionStatusResponse)
	assert.Equal(tst, result.MessageId, msgId)
	assert.True(tst, called1)
	assert.True(tst, called2)

}

func TestExecutesHooksForSpecTagsForScenarioExecutionEnding(tst *testing.T) {
	called1 := false
	called2 := false
	context := &t.GaugeContext{
		Hooks: []t.Hook{
			t.Hook{
				Type: t.AFTERSCENARIO,
				Impl: func(*m.ExecutionInfo) {
					called1 = true
				},
				Tags:     []string{"foo", "bar"},
				Operator: t.AND,
			},
			t.Hook{
				Type: t.AFTERSCENARIO,
				Impl: func(*m.ExecutionInfo) {
					called2 = true
				},
				Tags:     []string{"notfoo", "bar"},
				Operator: t.OR,
			},
		},
	}
	msgId := int64(12345)
	msg := &m.Message{
		MessageType: m.Message_ScenarioExecutionEnding,
		MessageId:   msgId,
		ScenarioExecutionEndingRequest: &m.ScenarioExecutionEndingRequest{
			CurrentExecutionInfo: &m.ExecutionInfo{
				CurrentSpec: &m.SpecInfo{
					Tags: []string{"foo", "bar"},
				},
			},
		},
	}

	p := ScenarioExecutionEndingProcessor{}

	result := p.Process(msg, context)

	assert.Equal(tst, result.MessageType, m.Message_ExecutionStatusResponse)
	assert.Equal(tst, result.MessageId, msgId)
	assert.True(tst, called1)
	assert.True(tst, called2)

}

func TestReportErrorIfHookFailsForScenarioExecutionEnding(tst *testing.T) {
	called1 := false
	called2 := false
	context := &t.GaugeContext{
		Hooks: []t.Hook{
			t.Hook{
				Type: t.AFTERSCENARIO,
				Impl: func(*m.ExecutionInfo) {
					called1 = true
				},
				Tags:     []string{"foo", "bar"},
				Operator: t.AND,
			},
			t.Hook{
				Type: t.AFTERSCENARIO,
				Impl: func(*m.ExecutionInfo) {
					called2 = true
					if 1 == 1 {
						panic(errors.New("Execution failed"))
					}
				},
				Tags:     []string{"notfoo", "bar"},
				Operator: t.OR,
			},
		},
	}
	msgId := int64(12345)
	msg := &m.Message{
		MessageType: m.Message_ScenarioExecutionEnding,
		MessageId:   msgId,
		ScenarioExecutionEndingRequest: &m.ScenarioExecutionEndingRequest{
			CurrentExecutionInfo: &m.ExecutionInfo{
				CurrentScenario: &m.ScenarioInfo{
					Tags: []string{"foo", "bar"},
				},
			},
		},
	}

	p := ScenarioExecutionEndingProcessor{}

	result := p.Process(msg, context)

	assert.True(tst, called1)
	assert.True(tst, called2)
	assert.Equal(tst, result.MessageType, m.Message_ExecutionStatusResponse)
	assert.Equal(tst, result.MessageId, msgId)
	assert.True(tst, result.ExecutionStatusResponse.ExecutionResult.Failed)
	assert.Equal(tst, result.ExecutionStatusResponse.ExecutionResult.ErrorMessage, "Execution failed")

}
