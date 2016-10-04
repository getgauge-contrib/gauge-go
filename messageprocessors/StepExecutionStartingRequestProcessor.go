package messageprocessors

import (
	m "github.com/getgauge-contrib/gauge-go/gauge_messages"
	t "github.com/getgauge-contrib/gauge-go/testsuit"
)

type StepExecutionStartingRequestProcessor struct{}

func (r *StepExecutionStartingRequestProcessor) Process(msg *m.Message, context *t.GaugeContext) *m.Message {
	tags := msg.GetStepExecutionStartingRequest().GetCurrentExecutionInfo().GetCurrentSpec().GetTags()
	hooks := context.GetHooks(t.BEFORESTEP, tags)
	context.ClearCustomMessages()
	return executeHooks(hooks, msg)
}
