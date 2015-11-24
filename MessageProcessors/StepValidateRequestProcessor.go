package messageprocessors

import (
	m "github.com/manuviswam/gauge-go/gauge_messages"
	t "github.com/manuviswam/gauge-go/testsuit"
)

type StepValidateRequestProcessor struct{}

func (s *StepValidateRequestProcessor) Process(msg *m.Message, steps []t.Step)*m.Message {
	stepDesc := msg.StepValidateRequest.StepText
	valid := isValid(steps, stepDesc)
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
