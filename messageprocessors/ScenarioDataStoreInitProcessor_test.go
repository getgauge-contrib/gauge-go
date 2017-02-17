package messageprocessors

import (
	"testing"

	m "github.com/getgauge-contrib/gauge-go/gauge_messages"
	t "github.com/getgauge-contrib/gauge-go/testsuit"
	"github.com/stretchr/testify/assert"
)

func TestShouldReturnExecutionStatusResponseWithSameIdForScenarioDataStoreInit(tst *testing.T) {
	msgId := int64(12345)
	context := &t.GaugeContext{
		Steps: make([]t.Step, 0),
	}
	msg := &m.Message{
		MessageType: m.Message_ScenarioDataStoreInit,
		MessageId:   msgId,
	}

	p := ScenarioDataStoreInitProcessor{}

	result := p.Process(msg, context)

	assert.Equal(tst, result.MessageType, m.Message_ExecutionStatusResponse)
	assert.Equal(tst, result.MessageId, msgId)
}

func TestShouldResetScenarioDataStore(tst *testing.T) {
	msgId := int64(12345)
	context := &t.GaugeContext{
		ScenarioStore: make(map[string]interface{}),
	}
	msg := &m.Message{
		MessageType: m.Message_ScenarioDataStoreInit,
		MessageId:   msgId,
	}

	context.ScenarioStore["foo"] = "bar"

	p := ScenarioDataStoreInitProcessor{}

	result := p.Process(msg, context)

	assert.Equal(tst, result.MessageType, m.Message_ExecutionStatusResponse)
	assert.Equal(tst, result.MessageId, msgId)
	assert.Nil(tst, context.ScenarioStore["foo"])
}
