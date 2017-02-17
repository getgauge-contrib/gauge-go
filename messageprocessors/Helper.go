package messageprocessors

import (
	m "github.com/getgauge-contrib/gauge-go/gauge_messages"
	t "github.com/getgauge-contrib/gauge-go/testsuit"
)

func executeHooks(hooks []t.Hook, msg *m.Message) *m.Message {
	var res *m.ProtoExecutionResult
	var totalExecutionTime int64
	for _, hook := range hooks {
		res = hook.Execute()
		totalExecutionTime += res.GetExecutionTime()
		if res.GetFailed() {
			return &m.Message{
				MessageType:             m.Message_ExecutionStatusResponse,
				MessageId:               msg.MessageId,
				ExecutionStatusResponse: &m.ExecutionStatusResponse{ExecutionResult: res},
			}
		}
	}
	return createResponseMessage(msg.MessageId, totalExecutionTime, nil)
}

func createResponseMessage(msgId int64, executionTime int64, err error) *m.Message {
	failed := false
	errorMsg := ""
	if err != nil {
		failed = true
		errorMsg = err.Error()
	}
	return &m.Message{
		MessageType: m.Message_ExecutionStatusResponse,
		MessageId:   msgId,
		ExecutionStatusResponse: &m.ExecutionStatusResponse{
			ExecutionResult: &m.ProtoExecutionResult{
				Failed:        failed,
				ExecutionTime: executionTime,
				ErrorMessage:  errorMsg,
			},
		},
	}
}
