package messageprocessors

import (
	m "github.com/getgauge-contrib/gauge-go/gauge_messages"
	t "github.com/getgauge-contrib/gauge-go/testsuit"
)

type StepNamesRequestProcessor struct{}

func (s *StepNamesRequestProcessor) Process(msg *m.Message, context *t.GaugeContext) *m.Message {
	return &m.Message{
		MessageType: m.Message_StepNamesResponse,
		MessageId:   msg.MessageId,
		StepNamesResponse: &m.StepNamesResponse{
			Steps: getAllDescriptions(context.Steps),
		},
	}
}

func getAllDescriptions(steps []t.Step) []string {
	descs := make([]string, len(steps))
	for _, step := range steps {
		descs = append(descs, step.Description)
	}
	return descs
}
