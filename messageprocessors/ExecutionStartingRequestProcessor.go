package messageprocessors

import (
	m "github.com/getgauge-contrib/gauge-go/gauge_messages"
	t "github.com/getgauge-contrib/gauge-go/testsuit"
)

type ExecutionStartingRequestProcessor struct{}

func (r *ExecutionStartingRequestProcessor) Process(msg *m.Message, context *t.GaugeContext) *m.Message {

	tags := msg.GetExecutionStartingRequest().GetCurrentExecutionInfo().GetCurrentScenario().GetTags()
	hooks := context.GetHooks(t.BEFORESUITE, tags)

	return executeHooks(hooks, msg)
}
