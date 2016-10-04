package messageutil

import (
	"bytes"
	"fmt"
	"net"

	c "github.com/getgauge-contrib/gauge-go/constants"
	m "github.com/getgauge-contrib/gauge-go/gauge_messages"
	"github.com/golang/protobuf/proto"
)

func ReadMessage(conn net.Conn) (*m.Message, error) {
	b, err := readMessageBytes(conn)
	if err != nil {
		return nil, err
	}

	msg, err := decodeMessage(b)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func readMessageBytes(conn net.Conn) ([]byte, error) {
	buffer := new(bytes.Buffer)
	data := make([]byte, c.MaxMessageSize)
	for {
		n, err := conn.Read(data)
		if err != nil {
			conn.Close()
			return nil, fmt.Errorf("Connection closed [%s] cause: %s", conn.RemoteAddr(), err.Error())
		}

		buffer.Write(data[0:n])

		messageLength, bytesRead := proto.DecodeVarint(buffer.Bytes())
		if messageLength > 0 && messageLength < uint64(buffer.Len()) {
			return buffer.Bytes()[bytesRead : messageLength+uint64(bytesRead)], nil
		}
	}
}

func decodeMessage(data []byte) (*m.Message, error) {
	message := new(m.Message)
	err := proto.Unmarshal(data, message)
	return message, err
}
