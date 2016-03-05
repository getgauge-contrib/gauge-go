package messageprocessors

import (
	m "github.com/manuviswam/gauge-go/gauge_messages"
	t "github.com/manuviswam/gauge-go/testsuit"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestShouldRunStep(tst *testing.T) {
	stepText := "Step description"
	msgId := int64(12345)
	called := false
	context := t.GaugeContext{
		Steps : []t.Step{t.Step{
				Description: stepText,
				Impl:        func(args ...interface{}) { called = true },
			},
		},
	}

	msg := &m.Message{
		MessageType: m.Message_ExecuteStep.Enum(),
		MessageId:   &msgId,
		ExecuteStepRequest: &m.ExecuteStepRequest{
			ParsedStepText: &stepText,
		},
	}

	p := ExecuteStepProcessor{}

	p.Process(msg, context)

	assert.True(tst, called)

}

func TestShouldRunReturnExecutionStatusResponseWithSameId(tst *testing.T) {
	stepText := "Step description"
	msgId := int64(12345)
	called := false
	context := t.GaugeContext{
		Steps : []t.Step{t.Step{
				Description: stepText,
				Impl:        func(args ...interface{}) { called = true },
			},
		},
	}

	msg := &m.Message{
		MessageType: m.Message_ExecuteStep.Enum(),
		MessageId:   &msgId,
		ExecuteStepRequest: &m.ExecuteStepRequest{
			ParsedStepText: &stepText,
		},
	}

	p := ExecuteStepProcessor{}

	result := p.Process(msg, context)

	assert.Equal(tst, result.MessageType, m.Message_ExecutionStatusResponse.Enum())
	assert.Equal(tst, *result.MessageId, msgId)
}
