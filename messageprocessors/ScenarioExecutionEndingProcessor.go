package messageprocessors

import (
	m "github.com/getgauge-contrib/gauge-go/gauge_messages"
	t "github.com/getgauge-contrib/gauge-go/testsuit"
)

type ScenarioExecutionEndingProcessor struct{}

func (r *ScenarioExecutionEndingProcessor) Process(msg *m.Message, context *t.GaugeContext) *m.Message {
	tags := mergeSpecAndScenarioTags(msg.GetScenarioExecutionEndingRequest().GetCurrentExecutionInfo())
	hooks := context.GetHooks(t.AFTERSCENARIO, tags)
	exInfo := msg.GetScenarioExecutionEndingRequest().GetCurrentExecutionInfo()

	return executeHooks(hooks, msg, exInfo)
}
