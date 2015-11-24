package runner

import (
	"fmt"
	"net"
	"os"

	c "github.com/manuviswam/gauge-go/constants"
	m "github.com/manuviswam/gauge-go/gauge_messages"
	mp "github.com/manuviswam/gauge-go/messageprocessors"
	mu "github.com/manuviswam/gauge-go/messageutil"
	t "github.com/manuviswam/gauge-go/testsuit"
	"regexp"
)

var steps []t.Step
var processors mp.ProcessorDictionary

func init() {
	steps = make([]t.Step, 0)
	processors = mp.ProcessorDictionary{}
	processors[*m.Message_StepNamesRequest.Enum()] = &mp.StepNamesRequestProcessor{}
	processors[*m.Message_StepValidateRequest.Enum()] = &mp.StepValidateRequestProcessor{}
	processors[*m.Message_SuiteDataStoreInit.Enum()] = &mp.SuiteDatastoreInitRequestProcessor{}
	processors[*m.Message_ExecutionStarting.Enum()] = &mp.ExecutionStartingRequestProcessor{}
	processors[*m.Message_SpecExecutionStarting.Enum()] = &mp.SpecExecutionStartingRequestProcessor{}
	processors[*m.Message_ScenarioExecutionStarting.Enum()] = &mp.ScenarioExecutionStartingRequestProcessor{}
	processors[*m.Message_StepExecutionStarting.Enum()] = &mp.StepExecutionStartingRequestProcessor{}
	processors[*m.Message_ExecuteStep.Enum()] = &mp.ExecuteStepProcessor{}
	processors[*m.Message_ExecutionEnding.Enum()] = &mp.ExecutionEndingProcessor{}
	processors[*m.Message_StepExecutionEnding.Enum()] = &mp.StepExecutionEndingProcessor{}
	processors[*m.Message_ScenarioExecutionEnding.Enum()] = &mp.ScenarioExecutionEndingProcessor{}
	processors[*m.Message_SpecExecutionEnding.Enum()] = &mp.SpecExecutionEndingProcessor{}
}

func Describe(stepDesc string, impl func()) bool {
	desc := parseDesc(stepDesc)
	step := t.Step{
		Description: desc,
		Impl:        impl,
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
		if *msg.MessageType.Enum() == *m.Message_KillProcessRequest.Enum() {
			return
		}

		processor := processors[*msg.MessageType.Enum()]

		if processor == nil {
			fmt.Println("Unable to find processor for message type : ", msg.MessageType)
			return
		}
		msgToSend := processor.Process(msg, steps)

		mu.WriteGaugeMessage(msgToSend, conn)
	}
}

func parseDesc(desc string) string {
	re := regexp.MustCompile("<(.*?)>")
	return re.ReplaceAllLiteralString(desc, "{}")
}
