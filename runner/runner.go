package runner

import (
	"fmt"
	"net"
	"os"

	"github.com/golang/protobuf/proto"
	c "github.com/manuviswam/gauge-go/constants"
	t "github.com/manuviswam/gauge-go/testsuit"
	m "github.com/manuviswam/gauge-go/gauge_messages"
	mu "github.com/manuviswam/gauge-go/messageutil"
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
	fmt.Println("We have got ", len(steps), " step implementations") // remove
	fmt.Println("Steps\n========") // remove
	fmt.Println(getAllDescriptions()) // remove

	var gaugePort = os.Getenv(c.GaugePortVariable)

	fmt.Println("Connecting port:", gaugePort) // remove
	conn, err := net.Dial("tcp", net.JoinHostPort("127.0.0.1", gaugePort))
	defer conn.Close()
	if err != nil {
		fmt.Println("dial error:", err)
		return
	}
	for {
		msg, err := mu.ReadMessage(conn)
		if err != nil {
			fmt.Println("Error reading message : ", err)
			return
		}
		fmt.Println("Message received : ", msg) // remove
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

func getAllDescriptions() []string {
	descs := make([]string, len(steps))
	for _, step := range steps {
		descs = append(descs, step.Description)
	}
	return descs
}