package messageprocessors

import (
	"fmt"

	m "github.com/manuviswam/gauge-go/gauge_messages"
	"github.com/manuviswam/gauge-go/models"
	t "github.com/manuviswam/gauge-go/testsuit"
)

type ExecuteStepProcessor struct{}

func (r *ExecuteStepProcessor) Process(msg *m.Message, context *t.GaugeContext) *m.Message {
	step, err := context.GetStepByDesc(*msg.ExecuteStepRequest.ParsedStepText)
	if err != nil {
		// if step implementation not found
		fmt.Println(err.Error())
	}
	args := getArgs(msg.ExecuteStepRequest)
	exeRes := step.Execute(args...)

	return &m.Message{
		MessageType: m.Message_ExecutionStatusResponse.Enum(),
		MessageId:   msg.MessageId,
		ExecutionStatusResponse: &m.ExecutionStatusResponse{
			ExecutionResult: exeRes,
		},
	}
}

func getArgs(r *m.ExecuteStepRequest) []interface{} {
	var args []interface{}
	for _, param := range r.GetParameters() {
		if *param.ParameterType.Enum() == *m.Parameter_Table.Enum() || *param.ParameterType.Enum() == *m.Parameter_Special_Table.Enum() {
			args = append(args, models.CreateTableFromProtoTable(param.Table))
		} else {
			args = append(args, *param.Value)
		}
	}
	return args
}
