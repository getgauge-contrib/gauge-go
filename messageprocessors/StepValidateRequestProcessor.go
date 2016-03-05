package messageprocessors

import (
	m "github.com/manuviswam/gauge-go/gauge_messages"
	t "github.com/manuviswam/gauge-go/testsuit"
)

type StepValidateRequestProcessor struct{}

func (s *StepValidateRequestProcessor) Process(msg *m.Message, context t.GaugeContext)*m.Message {
	stepDesc := msg.StepValidateRequest.StepText
	//TODO validate method signature
	valid := isValid(context.Steps, stepDesc)
	errorMsg := ""

	if !valid {
		errorMsg = "No implementation found for : " + *stepDesc
	}

	return &m.Message{
		MessageType: m.Message_StepValidateResponse.Enum(),
		MessageId:   msg.MessageId,
		StepValidateResponse: &m.StepValidateResponse{
			IsValid:      &valid,
			ErrorMessage: &errorMsg,
		},
	}
}

func isValid(steps []t.Step, desc *string) bool {
	for _, step := range steps {
		if step.Description == *desc {
			return true
		}
	}
	return false
}
