package messageprocessors

import (
	m "github.com/manuviswam/gauge-go/gauge_messages"
	t "github.com/manuviswam/gauge-go/testsuit"
)

type SpecExecutionStartingRequestProcessor struct{}

func (r *SpecExecutionStartingRequestProcessor) Process(msg *m.Message, steps []t.Step) *m.Message {
	//TODO do the intended operation here. Right now I am focused on getting the first test running.
	//TODO So I am replying with whatever this function is supposed to do is a success.
	failed := false
	executionTime := int64(1)
	return &m.Message{
		MessageType: m.Message_ExecutionStatusResponse.Enum(),
		MessageId:   msg.MessageId,
		ExecutionStatusResponse: &m.ExecutionStatusResponse{
			ExecutionResult: &m.ProtoExecutionResult{
				Failed:        &failed,
				ExecutionTime: &executionTime,
			},
		},
	}
}
