package messageprocessors

import (
	m "github.com/manuviswam/gauge-go/gauge_messages"
	t "github.com/manuviswam/gauge-go/testsuit"
)

type SpecDataStoreInitProcessor struct{}

func (r *SpecDataStoreInitProcessor) Process(msg *m.Message, context *t.GaugeContext) *m.Message {
	context.SpecStore = make(map[string]interface{})
	return createResponseMessage(msg.MessageId, int64(0), nil)
}
