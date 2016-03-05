package messageprocessors

import (
	m "github.com/manuviswam/gauge-go/gauge_messages"
	t "github.com/manuviswam/gauge-go/testsuit"
	"time"
)

func executeHooks(hooks []t.Hook) (int64, error) {
	var err error
	start := time.Now()
	for _, hook := range hooks {
		err = hook.Impl()
		if err != nil {
			break
		}
	}
	return time.Since(start).Nanoseconds(), err
}

func createResponseMessage(msgId *int64, executionTime int64, err error) *m.Message {
	failed := false
	errorMsg := ""
	if err != nil {
		failed = true
		errorMsg = err.Error()
	}
	return &m.Message{
		MessageType: m.Message_ExecutionStatusResponse.Enum(),
		MessageId:   msgId,
		ExecutionStatusResponse: &m.ExecutionStatusResponse{
			ExecutionResult: &m.ProtoExecutionResult{
				Failed:        &failed,
				ExecutionTime: &executionTime,
				ErrorMessage:  &errorMsg,
			},
		},
	}
}
