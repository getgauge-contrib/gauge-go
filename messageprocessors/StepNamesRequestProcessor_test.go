package messageprocessors

import (
	"testing"

	m "github.com/getgauge-contrib/gauge-go/gauge_messages"
	t "github.com/getgauge-contrib/gauge-go/testsuit"
	"github.com/stretchr/testify/assert"
)

func TestShouldReturnStepNamesResponseWithSameIdForStepnamesRequest(tst *testing.T) {
	msgId := int64(12345)
	context := &t.GaugeContext{
		Steps: make([]t.Step, 0),
	}

	msg := &m.Message{
		MessageType: m.Message_StepNamesRequest,
		MessageId:   msgId,
	}

	p := StepNamesRequestProcessor{}

	result := p.Process(msg, context)

	assert.Equal(tst, result.MessageType, m.Message_StepNamesResponse)
	assert.Equal(tst, result.MessageId, msgId)
}

func TestShouldReturnAllStepNames(tst *testing.T) {
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
	}

	p := StepNamesRequestProcessor{}

	result := p.Process(msg, context)

	assert.Contains(tst, result.StepNamesResponse.Steps, stepText)
}
