package messageprocessors

import (
	m "github.com/manuviswam/gauge-go/gauge_messages"
	t "github.com/manuviswam/gauge-go/testsuit"
)

type ScenarioExecutionEndingProcessor struct{}

func (r *ScenarioExecutionEndingProcessor) Process(msg *m.Message, context *t.GaugeContext) *m.Message {
	tags := msg.GetScenarioExecutionEndingRequest().GetCurrentExecutionInfo().GetCurrentScenario().GetTags()
	hooks := context.GetHooks(t.AFTERSCENARIO, tags)

	executionTime, err := executeHooks(hooks)
	return createResponseMessage(msg.MessageId, executionTime, err)
}
