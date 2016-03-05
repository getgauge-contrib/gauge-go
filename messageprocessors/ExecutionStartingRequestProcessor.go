package messageprocessors

import (
	m "github.com/manuviswam/gauge-go/gauge_messages"
	t "github.com/manuviswam/gauge-go/testsuit"
	"time"
)

type ExecutionStartingRequestProcessor struct{}

func (r *ExecutionStartingRequestProcessor) Process(msg *m.Message, context t.GaugeContext) *m.Message {
	var failed = false
	var executionTime int64
	var errorMsg string

	tags := msg.GetExecutionStartingRequest().GetCurrentExecutionInfo().GetCurrentScenario().GetTags()
	beforeSuiteHooks := context.GetHooks(t.BEFORESUITE, tags)

	start := time.Now()
	for _, hook := range beforeSuiteHooks   {
		err := hook.Impl()
		if err != nil {
			failed = true
			errorMsg = err.Error()
		}
	}
	executionTime = time.Since(start).Nanoseconds()

	return &m.Message{
		MessageType: m.Message_ExecutionStatusResponse.Enum(),
		MessageId:   msg.MessageId,
		ExecutionStatusResponse: &m.ExecutionStatusResponse{
			ExecutionResult: &m.ProtoExecutionResult{
				Failed:        &failed,
				ExecutionTime: &executionTime,
				ErrorMessage:  &errorMsg,
			},
		},
	}
}
