package messageutil

import (
	"net"

	"github.com/golang/protobuf/proto"
	"github.com/getgauge/common"
	"github.com/manuviswam/gauge-go/gauge_messages"
)

func Write(conn net.Conn, messageBytes []byte) error {
	messageLen := proto.EncodeVarint(uint64(len(messageBytes)))
	data := append(messageLen, messageBytes...)
	_, err := conn.Write(data)
	return err
}

func WriteGaugeMessage(message *gauge_messages.Message, conn net.Conn) error {
	messageId := common.GetUniqueID()
	message.MessageId = &messageId

	data, err := proto.Marshal(message)
	if err != nil {
		return err
	}
	return Write(conn, data)
}
