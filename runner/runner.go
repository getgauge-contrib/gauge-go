package runner

import (
	"fmt"
	"os"
	"net"

	"github.com/manuviswam/gauge-go/constants"
	"github.com/golang/protobuf/proto"
	"github.com/manuviswam/gauge-go/gauge_messages"
)

var steps map[string]func()

func init() {
	steps = make(map[string]func())
}

func Describe(stepDesc string, impl func()) bool {
	steps[stepDesc] = impl
	return true
}

func Run() {
	fmt.Println("We have got ", len(steps), " step implementations")
	fmt.Println("Steps\n========")
	for step, _ := range steps {
		fmt.Println(step)
	}

	var gaugePort = os.Getenv(constants.GaugePortVariable)

	fmt.Println("Connecting port:", gaugePort)
	conn, err := net.Dial("tcp", net.JoinHostPort("127.0.0.1", gaugePort))
	if err != nil {
		fmt.Println("dial error:", err)
		return
	}
	defer conn.Close()
	b := make([]byte, constants.MaxMessageSize)
	for {
		conn.Read(b)
		processMessage(b)
		fmt.Println("total size:",len(b))
	}
}

func processMessage(data []byte) error {
	message := gauge_messages.Message{}
	proto.Unmarshal(data, &message)
	fmt.Println("Message recieved : ", message)
	return nil
}
