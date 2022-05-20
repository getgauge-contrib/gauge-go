package messageprocessors

import (
	m "github.com/getgauge-contrib/gauge-go/gauge_messages"
	t "github.com/getgauge-contrib/gauge-go/testsuit"
)

type ExecutionStartingRequestProcessor struct{}

func (r *ExecutionStartingRequestProcessor) Process(msg *m.Message, context *t.GaugeContext) *m.Message {

	tags := []string{}
	hooks := context.GetHooks(t.BEFORESUITE, tags)
	exInfo := msg.GetExecutionStartingRequest().GetCurrentExecutionInfo()

	return executeHooks(hooks, msg, exInfo)
}
