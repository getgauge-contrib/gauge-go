package messageprocessors

import (
	"errors"
	m "github.com/manuviswam/gauge-go/gauge_messages"
	t "github.com/manuviswam/gauge-go/testsuit"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestShouldReturnExecutionStatusResponseWithSameIdForSpecExecutionStarting(tst *testing.T) {
	msgId := int64(12345)
	context := t.GaugeContext{
		Steps: make([]t.Step, 0),
	}

	msg := &m.Message{
		MessageType: m.Message_SpecExecutionStarting.Enum(),
		MessageId:   &msgId,
	}

	p := SpecExecutionStartingRequestProcessor{}

	result := p.Process(msg, context)

	assert.Equal(tst, result.MessageType, m.Message_ExecutionStatusResponse.Enum())
	assert.Equal(tst, *result.MessageId, msgId)
}

func TestExecutesHooksForTheTagsForSpecExecutionStarting(tst *testing.T) {
	called1 := false
	called2 := false
	context := t.GaugeContext{
		Hooks: []t.Hook{
			t.Hook{
				Type: t.BEFORESPEC,
				Impl: func() error {
					called1 = true
					return nil
				},
				Tags:     []string{"foo", "bar"},
				Operator: t.AND,
			},
			t.Hook{
				Type: t.BEFORESPEC,
				Impl: func() error {
					called2 = true
					return nil
				},
				Tags:     []string{"notfoo", "bar"},
				Operator: t.OR,
			},
		},
	}
	msgId := int64(12345)
	msg := &m.Message{
		MessageType: m.Message_SpecExecutionStarting.Enum(),
		MessageId:   &msgId,
		SpecExecutionStartingRequest: &m.SpecExecutionStartingRequest{
			CurrentExecutionInfo: &m.ExecutionInfo{
				CurrentSpec: &m.SpecInfo{
					Tags: []string{"foo", "bar"},
				},
			},
		},
	}

	p := SpecExecutionStartingRequestProcessor{}

	result := p.Process(msg, context)

	assert.Equal(tst, result.MessageType, m.Message_ExecutionStatusResponse.Enum())
	assert.Equal(tst, *result.MessageId, msgId)
	assert.True(tst, called1)
	assert.True(tst, called2)

}

func TestReportErrorIfHookFailsForSpecExecutionStarting(tst *testing.T) {
	called1 := false
	called2 := false
	context := t.GaugeContext{
		Hooks: []t.Hook{
			t.Hook{
				Type: t.BEFORESPEC,
				Impl: func() error {
					called1 = true
					return nil
				},
				Tags:     []string{"foo", "bar"},
				Operator: t.AND,
			},
			t.Hook{
				Type: t.BEFORESPEC,
				Impl: func() error {
					called2 = true
					return errors.New("Execution failed")
				},
				Tags:     []string{"notfoo", "bar"},
				Operator: t.OR,
			},
		},
	}
	msgId := int64(12345)
	msg := &m.Message{
		MessageType: m.Message_SpecExecutionStarting.Enum(),
		MessageId:   &msgId,
		SpecExecutionStartingRequest: &m.SpecExecutionStartingRequest{
			CurrentExecutionInfo: &m.ExecutionInfo{
				CurrentSpec: &m.SpecInfo{
					Tags: []string{"foo", "bar"},
				},
			},
		},
	}

	p := SpecExecutionStartingRequestProcessor{}

	result := p.Process(msg, context)

	assert.True(tst, called1)
	assert.True(tst, called2)
	assert.Equal(tst, result.MessageType, m.Message_ExecutionStatusResponse.Enum())
	assert.Equal(tst, *result.MessageId, msgId)
	assert.True(tst, *result.ExecutionStatusResponse.ExecutionResult.Failed)
	assert.Equal(tst, *result.ExecutionStatusResponse.ExecutionResult.ErrorMessage, "Execution failed")

}
