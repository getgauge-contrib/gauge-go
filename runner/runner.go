package runner

import (
	"fmt"
	"net"
	"os"

	c "github.com/manuviswam/gauge-go/constants"
	t "github.com/manuviswam/gauge-go/testsuit"
	m "github.com/manuviswam/gauge-go/gauge_messages"
	mu "github.com/manuviswam/gauge-go/messageutil"
	mp "github.com/manuviswam/gauge-go/messageprocessors"
)

var steps []t.Step
var processors mp.ProcessorDictionary

func init() {
	steps = make([]t.Step, 0)
	processors = mp.ProcessorDictionary{}
	processors[*m.Message_StepNamesRequest.Enum()] = &mp.StepNamesRequestProcessor{}
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

		processor := processors[*msg.MessageType.Enum()]

		if processor == nil {
			fmt.Println("Unable to find processor for message type : ", msg.MessageType)
			return
		}
		processor.Process(conn, msg, steps)
	}
}

