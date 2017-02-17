package messageprocessors

import (
	gm "github.com/getgauge-contrib/gauge-go/gauge_messages"
	t "github.com/getgauge-contrib/gauge-go/testsuit"
)

type StepValidateRequestProcessor struct{}

func (s *StepValidateRequestProcessor) Process(msg *gm.Message, context *t.GaugeContext) *gm.Message {
	stepDesc := msg.StepValidateRequest.StepText
	//TODO validate method signature
	valid := isValid(context.Steps, stepDesc)
	errorMsg := ""

	res := &gm.StepValidateResponse{}
	if !valid {
		errorMsg = "No implementation found for : " + stepDesc
		res.ErrorMessage = errorMsg
		res.ErrorType = gm.StepValidateResponse_STEP_IMPLEMENTATION_NOT_FOUND
	}
	res.IsValid = valid

	return &gm.Message{
		MessageType:          gm.Message_StepValidateResponse,
		MessageId:            msg.MessageId,
		StepValidateResponse: res,
	}
}

func isValid(steps []t.Step, desc string) bool {
	for _, step := range steps {
		if step.Description == desc {
			return true
		}
	}
	return false
}
