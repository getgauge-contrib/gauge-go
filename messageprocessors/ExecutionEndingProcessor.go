package messageprocessors

import (
	m "github.com/getgauge-contrib/gauge-go/gauge_messages"
	t "github.com/getgauge-contrib/gauge-go/testsuit"
)

type ExecutionEndingProcessor struct{}

func (r *ExecutionEndingProcessor) Process(msg *m.Message, context *t.GaugeContext) *m.Message {
	tags := []string{}
	hooks := context.GetHooks(t.AFTERSUITE, tags)
	exInfo := msg.GetExecutionEndingRequest().GetCurrentExecutionInfo()

	return executeHooks(hooks, msg, exInfo)
}
