package messageutil

import (
	"net"

	"github.com/getgauge-contrib/gauge-go/gauge_messages"
	"github.com/golang/protobuf/proto"
)

func Write(conn net.Conn, messageBytes []byte) error {
	messageLen := proto.EncodeVarint(uint64(len(messageBytes)))
	data := append(messageLen, messageBytes...)
	_, err := conn.Write(data)
	return err
}

func WriteGaugeMessage(message *gauge_messages.Message, conn net.Conn) error {
	data, err := proto.Marshal(message)
	if err != nil {
		return err
	}
	return Write(conn, data)
}
