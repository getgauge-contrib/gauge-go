package messageprocessors

import (
	m "github.com/manuviswam/gauge-go/gauge_messages"
	t "github.com/manuviswam/gauge-go/testsuit"
)

type StepExecutionEndingProcessor struct{}

func (r *StepExecutionEndingProcessor) Process(msg *m.Message, context t.GaugeContext) *m.Message {
	tags := msg.GetStepExecutionEndingRequest().GetCurrentExecutionInfo().GetCurrentSpec().GetTags()
	hooks := context.GetHooks(t.AFTERSTEP, tags)

	executionTime, err := executeHooks(hooks)
	return createResponseMessage(msg.MessageId, executionTime, err)
}
