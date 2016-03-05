package messageprocessors

import (
	m "github.com/manuviswam/gauge-go/gauge_messages"
	t "github.com/manuviswam/gauge-go/testsuit"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestShouldReturnStepNamesResponseWithSameIdForStepValidateRequest(tst *testing.T) {
	stepText := "Step description"
	msgId := int64(12345)
	context := t.GaugeContext{
		Steps: make([]t.Step, 0),
	}

	msg := &m.Message{
		MessageType: m.Message_StepNamesRequest.Enum(),
		MessageId:   &msgId,
		StepValidateRequest: &m.StepValidateRequest{
			StepText: &stepText,
		},
	}

	p := StepValidateRequestProcessor{}

	result := p.Process(msg, context)

	assert.Equal(tst, result.MessageType, m.Message_StepValidateResponse.Enum())
	assert.Equal(tst, *result.MessageId, msgId)
}

func TestShouldValidateStep(tst *testing.T) {
	stepText := "Step description"
	msgId := int64(12345)
	context := t.GaugeContext{
		Steps: []t.Step{t.Step{
			Description: stepText,
			Impl:        func(args ...interface{}) {},
		},
		},
	}

	msg := &m.Message{
		MessageType: m.Message_StepNamesRequest.Enum(),
		MessageId:   &msgId,
		StepValidateRequest: &m.StepValidateRequest{
			StepText: &stepText,
		},
	}

	p := StepValidateRequestProcessor{}

	result := p.Process(msg, context)

	assert.True(tst, *result.StepValidateResponse.IsValid)
}
