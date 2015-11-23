package messageprocessors

import (
	"net"

	m "github.com/manuviswam/gauge-go/gauge_messages"
	t "github.com/manuviswam/gauge-go/testsuit"
)

type MessageProcessor interface {
	Process(net.Conn, *m.Message, []t.Step)
}

type ProcessorDictionary map[m.Message_MessageType]MessageProcessor

