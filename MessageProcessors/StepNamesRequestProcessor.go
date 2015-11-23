package messageprocessors

import (
	"github.com/golang/protobuf/proto"
	m "github.com/manuviswam/gauge-go/gauge_messages"
	t "github.com/manuviswam/gauge-go/testsuit"
	"net"
)

type StepNamesRequestProcessor struct {}

func (s *StepNamesRequestProcessor) Process(conn net.Conn, msg *m.Message, steps []t.Step)  {
	msgToSend := m.Message{
		MessageType: m.Message_StepNamesResponse.Enum(),
		MessageId:   msg.MessageId,
		StepNamesResponse: &m.StepNamesResponse{
			Steps: getAllDescriptions(steps),
		},
	}
	protoMsg, _ := proto.Marshal(&msgToSend)
	conn.Write(protoMsg)
}

func getAllDescriptions(steps []t.Step) []string {
	descs := make([]string, len(steps))
	for _, step := range steps {
		descs = append(descs, step.Description)
	}
	return descs
}