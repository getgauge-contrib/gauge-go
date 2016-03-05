package messageprocessors

import (
	m "github.com/manuviswam/gauge-go/gauge_messages"
	t "github.com/manuviswam/gauge-go/testsuit"
)

type ExecutionStartingRequestProcessor struct{}

func (r *ExecutionStartingRequestProcessor) Process(msg *m.Message, context t.GaugeContext) *m.Message {

	tags := msg.GetExecutionStartingRequest().GetCurrentExecutionInfo().GetCurrentScenario().GetTags()
	hooks := context.GetHooks(t.BEFORESUITE, tags)

	executionTime, err := executeHooks(hooks)
	return createResponseMessage(msg.MessageId, executionTime, err)
}
