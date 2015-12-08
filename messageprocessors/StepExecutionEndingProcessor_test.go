package messageprocessors

import (
	m "github.com/manuviswam/gauge-go/gauge_messages"
	t "github.com/manuviswam/gauge-go/testsuit"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestShouldReturnExecutionStatusResponseWithSameIdForStepExecutionEnding(tst *testing.T) {
	msgId := int64(12345)
	steps := make([]t.Step, 0)

	msg := &m.Message{
		MessageType: m.Message_StepExecutionEnding.Enum(),
		MessageId:   &msgId,
	}

	p := StepExecutionEndingProcessor{}

	result := p.Process(msg, steps)

	assert.Equal(tst, result.MessageType, m.Message_ExecutionStatusResponse.Enum())
	assert.Equal(tst, *result.MessageId, msgId)
}
