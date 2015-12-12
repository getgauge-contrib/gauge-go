package messageprocessors

import (
	m "github.com/manuviswam/gauge-go/gauge_messages"
	t "github.com/manuviswam/gauge-go/testsuit"
)

type MessageProcessor interface {
	Process(*m.Message, []t.Step) *m.Message
}

type ProcessorDictionary map[m.Message_MessageType]MessageProcessor
