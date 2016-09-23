package messageprocessors

import (
	"errors"
	"testing"

	m "github.com/manuviswam/gauge-go/gauge_messages"
	t "github.com/manuviswam/gauge-go/testsuit"
	"github.com/stretchr/testify/assert"
)

func TestShouldReturnExecutionStatusResponseWithSameIdForStepExecutionEnding(tst *testing.T) {
	msgId := int64(12345)
	context := &t.GaugeContext{
		Steps: make([]t.Step, 0),
	}

	msg := &m.Message{
		MessageType: m.Message_StepExecutionEnding.Enum(),
		MessageId:   &msgId,
	}

	p := StepExecutionEndingProcessor{}

	result := p.Process(msg, context)

	assert.Equal(tst, result.MessageType, m.Message_ExecutionStatusResponse.Enum())
	assert.Equal(tst, *result.MessageId, msgId)
}

func TestExecutesHooksForTheTagsForStepExecutionEnding(tst *testing.T) {
	called1 := false
	called2 := false
	context := &t.GaugeContext{
		Hooks: []t.Hook{
			t.Hook{
				Type: t.AFTERSTEP,
				Impl: func() {
					called1 = true
				},
				Tags:     []string{"foo", "bar"},
				Operator: t.AND,
			},
			t.Hook{
				Type: t.AFTERSTEP,
				Impl: func() {
					called2 = true
				},
				Tags:     []string{"notfoo", "bar"},
				Operator: t.OR,
			},
		},
	}
	msgId := int64(12345)
	msg := &m.Message{
		MessageType: m.Message_StepExecutionEnding.Enum(),
		MessageId:   &msgId,
		StepExecutionEndingRequest: &m.StepExecutionEndingRequest{
			CurrentExecutionInfo: &m.ExecutionInfo{
				CurrentSpec: &m.SpecInfo{
					Tags: []string{"foo", "bar"},
				},
			},
		},
	}

	p := StepExecutionEndingProcessor{}

	result := p.Process(msg, context)

	assert.Equal(tst, result.MessageType, m.Message_ExecutionStatusResponse.Enum())
	assert.Equal(tst, *result.MessageId, msgId)
	assert.True(tst, called1)
	assert.True(tst, called2)

}

func TestReportErrorIfHookFailsForStepExecutionEnding(tst *testing.T) {
	called1 := false
	called2 := false
	context := &t.GaugeContext{
		Hooks: []t.Hook{
			t.Hook{
				Type: t.AFTERSTEP,
				Impl: func() {
					called1 = true
				},
				Tags:     []string{"foo", "bar"},
				Operator: t.AND,
			},
			t.Hook{
				Type: t.AFTERSTEP,
				Impl: func() {
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
		MessageType: m.Message_StepExecutionEnding.Enum(),
		MessageId:   &msgId,
		StepExecutionEndingRequest: &m.StepExecutionEndingRequest{
			CurrentExecutionInfo: &m.ExecutionInfo{
				CurrentSpec: &m.SpecInfo{
					Tags: []string{"foo", "bar"},
				},
			},
		},
	}

	p := StepExecutionEndingProcessor{}

	result := p.Process(msg, context)

	assert.True(tst, called1)
	assert.True(tst, called2)
	assert.Equal(tst, result.MessageType, m.Message_ExecutionStatusResponse.Enum())
	assert.Equal(tst, *result.MessageId, msgId)
	assert.True(tst, *result.ExecutionStatusResponse.ExecutionResult.Failed)
	assert.Equal(tst, *result.ExecutionStatusResponse.ExecutionResult.ErrorMessage, "Execution failed")

}

func TestShouldReturnCustomMessagesInResult(tst *testing.T) {
	customMessages := []string{"my custom message"}
	called := false
	context := &t.GaugeContext{
		Hooks: []t.Hook{
			t.Hook{
				Type: t.AFTERSTEP,
				Impl: func() {
					called = true
				},
				Operator: t.AND,
			},
		},
		CustomMessageRegistry: customMessages,
	}
	msgId := int64(12345)
	msg := &m.Message{
		MessageType: m.Message_StepExecutionEnding.Enum(),
		MessageId:   &msgId,
		StepExecutionEndingRequest: &m.StepExecutionEndingRequest{
			CurrentExecutionInfo: &m.ExecutionInfo{
				CurrentSpec: &m.SpecInfo{},
			},
		},
	}

	p := StepExecutionEndingProcessor{}

	resMsg := p.Process(msg, context)

	assert.True(tst, called)
	assert.Equal(tst, resMsg.GetExecutionStatusResponse().GetExecutionResult().GetMessage(), customMessages)
}
