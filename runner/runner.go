package runner

import (
	"fmt"
	"net"
	"os"
	"bytes"

	"github.com/golang/protobuf/proto"
	c "github.com/manuviswam/gauge-go/constants"
	t "github.com/manuviswam/gauge-go/testsuit"
	m "github.com/manuviswam/gauge-go/gauge_messages"
)

var steps []t.Step

func init() {
	steps = make([]t.Step, 0)
}

func Describe(stepDesc string, impl func()) bool {
	step := t.Step{
		Description:stepDesc,
		Impl:impl,
	}
	steps = append(steps, step)
	return true
}

func Run() {
	fmt.Println("We have got ", len(steps), " step implementations")
	fmt.Println("Steps\n========")
	for _,step := range steps {
		fmt.Println(step.Description)
	}

	var gaugePort = os.Getenv(c.GaugePortVariable)

	fmt.Println("Connecting port:", gaugePort)
	conn, err := net.Dial("tcp", net.JoinHostPort("127.0.0.1", gaugePort))
	defer conn.Close()
	if err != nil {
		fmt.Println("dial error:", err)
		return
	}
	for {
		data, err := readMessageBytes(conn)
		if err != nil {
			fmt.Println("Error reading message : ", err)
		}
		msg, err := decodeMessage(data)
		if err != nil {
			fmt.Println("Error decoding message :", err)
		}
		fmt.Println("Message received : ", msg)
		msgToSend := m.Message{
			MessageType: m.Message_StepNamesResponse.Enum(),
			MessageId:   msg.MessageId,
			StepNamesResponse: &m.StepNamesResponse{
				Steps: getAllDescriptions(),
			},
		}
		protoMsg, _ := proto.Marshal(&msgToSend)
		conn.Write(protoMsg)
	}
}

func decodeMessage(data []byte) (*m.Message, error) {
	message := new(m.Message)
	err := proto.Unmarshal(data, message)
	return message, err
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

func getAllDescriptions() []string {
	descs := make([]string, len(steps))
	for _, step := range steps {
		descs = append(descs, step.Description)
	}
	return descs
}