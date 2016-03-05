package messageprocessors

import (
	m "github.com/manuviswam/gauge-go/gauge_messages"
	t "github.com/manuviswam/gauge-go/testsuit"
)

type SpecExecutionEndingProcessor struct{}

func (r *SpecExecutionEndingProcessor) Process(msg *m.Message, context *t.GaugeContext) *m.Message {
	tags := msg.GetSpecExecutionEndingRequest().GetCurrentExecutionInfo().GetCurrentSpec().GetTags()
	hooks := context.GetHooks(t.AFTERSPEC, tags)

	executionTime, err := executeHooks(hooks)
	return createResponseMessage(msg.MessageId, executionTime, err)
}
