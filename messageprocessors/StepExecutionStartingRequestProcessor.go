package messageprocessors

import (
	m "github.com/manuviswam/gauge-go/gauge_messages"
	t "github.com/manuviswam/gauge-go/testsuit"
)

type StepExecutionStartingRequestProcessor struct{}

func (r *StepExecutionStartingRequestProcessor) Process(msg *m.Message, context *t.GaugeContext) *m.Message {
	tags := msg.GetStepExecutionStartingRequest().GetCurrentExecutionInfo().GetCurrentSpec().GetTags()
	hooks := context.GetHooks(t.BEFORESTEP, tags)

	return executeHooks(hooks, msg)
}
