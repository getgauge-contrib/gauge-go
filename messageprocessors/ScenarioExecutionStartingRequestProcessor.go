package messageprocessors

import (
	m "github.com/getgauge-contrib/gauge-go/gauge_messages"
	t "github.com/getgauge-contrib/gauge-go/testsuit"
)

type ScenarioExecutionStartingRequestProcessor struct{}

func (r *ScenarioExecutionStartingRequestProcessor) Process(msg *m.Message, context *t.GaugeContext) *m.Message {
	tags := mergeSpecAndScenarioTags(msg.GetScenarioExecutionStartingRequest().GetCurrentExecutionInfo())
	hooks := context.GetHooks(t.BEFORESCENARIO, tags)
	exInfo := msg.GetScenarioExecutionStartingRequest().GetCurrentExecutionInfo()

	return executeHooks(hooks, msg, exInfo)
}
