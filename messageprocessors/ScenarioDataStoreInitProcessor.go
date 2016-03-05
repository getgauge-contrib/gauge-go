package messageprocessors

import (
	m "github.com/manuviswam/gauge-go/gauge_messages"
	t "github.com/manuviswam/gauge-go/testsuit"
)

type ScenarioDataStoreInitProcessor struct{}

func (r *ScenarioDataStoreInitProcessor) Process(msg *m.Message, context *t.GaugeContext) *m.Message {
	context.ScenarioStore = make(map[string]interface{})
	return createResponseMessage(msg.MessageId, int64(0), nil)
}
