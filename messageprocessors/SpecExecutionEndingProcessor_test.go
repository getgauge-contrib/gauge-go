package messageprocessors

import (
	"errors"
	"testing"

	m "github.com/getgauge-contrib/gauge-go/gauge_messages"
	t "github.com/getgauge-contrib/gauge-go/testsuit"
	"github.com/stretchr/testify/assert"
)

func TestShouldReturnExecutionStatusResponseWithSameIdForSpecExecutionEnding(tst *testing.T) {
	msgId := int64(12345)
	context := &t.GaugeContext{
		Steps: make([]t.Step, 0),
	}

	msg := &m.Message{
		MessageType: m.Message_SpecExecutionEnding,
		MessageId:   msgId,
	}

	p := SpecExecutionEndingProcessor{}

	result := p.Process(msg, context)

	assert.Equal(tst, result.MessageType, m.Message_ExecutionStatusResponse)
	assert.Equal(tst, result.MessageId, msgId)
}

func TestExecutesHooksForTheTagsForSpecExecutionEnding(tst *testing.T) {
	called1 := false
	called2 := false
	context := &t.GaugeContext{
		Hooks: []t.Hook{
			t.Hook{
				Type: t.AFTERSPEC,
				Impl: func(*m.ExecutionInfo) {
					called1 = true
				},
				Tags:     []string{"foo", "bar"},
				Operator: t.AND,
			},
			t.Hook{
				Type: t.AFTERSPEC,
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
		MessageType: m.Message_SpecExecutionEnding,
		MessageId:   msgId,
		SpecExecutionEndingRequest: &m.SpecExecutionEndingRequest{
			CurrentExecutionInfo: &m.ExecutionInfo{
				CurrentSpec: &m.SpecInfo{
					Tags: []string{"foo", "bar"},
				},
			},
		},
	}

	p := SpecExecutionEndingProcessor{}

	result := p.Process(msg, context)

	assert.Equal(tst, result.MessageType, m.Message_ExecutionStatusResponse)
	assert.Equal(tst, result.MessageId, msgId)
	assert.True(tst, called1)
	assert.True(tst, called2)

}

func TestReportErrorIfHookFailsForSpecExecutionEnding(tst *testing.T) {
	called1 := false
	called2 := false
	context := &t.GaugeContext{
		Hooks: []t.Hook{
			t.Hook{
				Type: t.AFTERSPEC,
				Impl: func(*m.ExecutionInfo) {
					called1 = true
				},
				Tags:     []string{"foo", "bar"},
				Operator: t.AND,
			},
			t.Hook{
				Type: t.AFTERSPEC,
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
		MessageType: m.Message_SpecExecutionEnding,
		MessageId:   msgId,
		SpecExecutionEndingRequest: &m.SpecExecutionEndingRequest{
			CurrentExecutionInfo: &m.ExecutionInfo{
				CurrentSpec: &m.SpecInfo{
					Tags: []string{"foo", "bar"},
				},
			},
		},
	}

	p := SpecExecutionEndingProcessor{}

	result := p.Process(msg, context)

	assert.True(tst, called1)
	assert.True(tst, called2)
	assert.Equal(tst, result.MessageType, m.Message_ExecutionStatusResponse)
	assert.Equal(tst, result.MessageId, msgId)
	assert.True(tst, result.ExecutionStatusResponse.ExecutionResult.Failed)
	assert.Equal(tst, result.ExecutionStatusResponse.ExecutionResult.ErrorMessage, "Execution failed")

}
