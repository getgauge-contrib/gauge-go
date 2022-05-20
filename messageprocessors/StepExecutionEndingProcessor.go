package messageprocessors

import (
	m "github.com/getgauge-contrib/gauge-go/gauge_messages"
	t "github.com/getgauge-contrib/gauge-go/testsuit"
)

type StepExecutionEndingProcessor struct{}

func (r *StepExecutionEndingProcessor) Process(msg *m.Message, context *t.GaugeContext) *m.Message {
	tags := mergeSpecAndScenarioTags(msg.GetStepExecutionEndingRequest().GetCurrentExecutionInfo())
	hooks := context.GetHooks(t.AFTERSTEP, tags)
	exInfo := msg.GetStepExecutionEndingRequest().GetCurrentExecutionInfo()

	res := executeHooks(hooks, msg, exInfo)
	res.GetExecutionStatusResponse().GetExecutionResult().Message = context.CustomMessageRegistry
	context.ClearCustomMessages()

	return res
}
