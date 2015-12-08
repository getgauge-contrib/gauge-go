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
	step := t.Step{
		Description: stepText,
		Impl:        func(args ...interface{}) { called = true },
	}
	steps := make([]t.Step, 0)
	steps = append(steps, step)

	msg := &m.Message{
		MessageType: m.Message_ExecuteStep.Enum(),
		MessageId:   &msgId,
		ExecuteStepRequest: &m.ExecuteStepRequest{
			ParsedStepText: &stepText,
		},
	}

	p := ExecuteStepProcessor{}

	p.Process(msg, steps)

	assert.True(tst, called)

}

func TestShouldRunReturnExecutionStatusResponseWithSameId(tst *testing.T) {
	stepText := "Step description"
	msgId := int64(12345)
	called := false
	step := t.Step{
		Description: stepText,
		Impl:        func(args ...interface{}) { called = true },
	}
	steps := make([]t.Step, 0)
	steps = append(steps, step)

	msg := &m.Message{
		MessageType: m.Message_ExecuteStep.Enum(),
		MessageId:   &msgId,
		ExecuteStepRequest: &m.ExecuteStepRequest{
			ParsedStepText: &stepText,
		},
	}

	p := ExecuteStepProcessor{}

	result := p.Process(msg, steps)

	assert.Equal(tst, result.MessageType, m.Message_ExecutionStatusResponse.Enum())
	assert.Equal(tst, *result.MessageId, msgId)
}
