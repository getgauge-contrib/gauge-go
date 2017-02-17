package messageprocessors

import (
	"testing"

	m "github.com/getgauge-contrib/gauge-go/gauge_messages"
	t "github.com/getgauge-contrib/gauge-go/testsuit"
	"github.com/stretchr/testify/assert"
)

func TestShouldReturnStepNamesResponseWithSameIdForStepValidateRequest(tst *testing.T) {
	stepText := "Step description"
	msgId := int64(12345)
	context := &t.GaugeContext{
		Steps: make([]t.Step, 0),
	}

	msg := &m.Message{
		MessageType: m.Message_StepNamesRequest,
		MessageId:   msgId,
		StepValidateRequest: &m.StepValidateRequest{
			StepText: stepText,
		},
	}

	p := StepValidateRequestProcessor{}

	result := p.Process(msg, context)

	assert.Equal(tst, result.MessageType, m.Message_StepValidateResponse)
	assert.Equal(tst, result.MessageId, msgId)
}

func TestShouldValidateStep(tst *testing.T) {
	stepText := "Step description"
	msgId := int64(12345)
	context := &t.GaugeContext{
		Steps: []t.Step{t.Step{
			Description: stepText,
			Impl:        func(args ...interface{}) {},
		},
		},
	}

	msg := &m.Message{
		MessageType: m.Message_StepNamesRequest,
		MessageId:   msgId,
		StepValidateRequest: &m.StepValidateRequest{
			StepText: stepText,
		},
	}

	p := StepValidateRequestProcessor{}

	result := p.Process(msg, context)

	assert.True(tst, result.StepValidateResponse.IsValid)
	assert.Equal(tst, result.StepValidateResponse.GetErrorMessage(), "")
}

func TestShouldValidateStepWhenNotFound(tst *testing.T) {
	stepText := "Step description"
	requiredStep := "hello"
	msgId := int64(12345)
	context := &t.GaugeContext{
		Steps: []t.Step{t.Step{
			Description: stepText,
			Impl:        func(args ...interface{}) {},
		},
		},
	}

	msg := &m.Message{
		MessageType: m.Message_StepNamesRequest,
		MessageId:   msgId,
		StepValidateRequest: &m.StepValidateRequest{
			StepText: requiredStep,
		},
	}

	p := StepValidateRequestProcessor{}

	result := p.Process(msg, context)

	assert.False(tst, result.StepValidateResponse.IsValid)
	assert.Equal(tst, result.StepValidateResponse.GetErrorMessage(), "No implementation found for : "+requiredStep)
}
