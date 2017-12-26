package gauge

import (
	"fmt"
	"net"
	"os"
	"reflect"

	"regexp"

	c "github.com/getgauge-contrib/gauge-go/constants"
	m "github.com/getgauge-contrib/gauge-go/gauge_messages"
	mp "github.com/getgauge-contrib/gauge-go/messageprocessors"
	mu "github.com/getgauge-contrib/gauge-go/messageutil"
	t "github.com/getgauge-contrib/gauge-go/testsuit"
)

var context *t.GaugeContext
var processors mp.ProcessorDictionary

func init() {
	context = &t.GaugeContext{
		Steps:                 make([]t.Step, 0),
		Hooks:                 make([]t.Hook, 0),
		SuiteStore:            make(map[string]interface{}),
		SpecStore:             make(map[string]interface{}),
		ScenarioStore:         make(map[string]interface{}),
		CustomMessageRegistry: make([]string, 0),
	}

	processors = mp.ProcessorDictionary{}
	processors[m.Message_StepNamesRequest] = &mp.StepNamesRequestProcessor{}
	processors[m.Message_StepValidateRequest] = &mp.StepValidateRequestProcessor{}
	processors[m.Message_SuiteDataStoreInit] = &mp.SuiteDataStoreInitRequestProcessor{}
	processors[m.Message_ExecutionStarting] = &mp.ExecutionStartingRequestProcessor{}
	processors[m.Message_SpecExecutionStarting] = &mp.SpecExecutionStartingRequestProcessor{}
	processors[m.Message_ScenarioExecutionStarting] = &mp.ScenarioExecutionStartingRequestProcessor{}
	processors[m.Message_StepExecutionStarting] = &mp.StepExecutionStartingRequestProcessor{}
	processors[m.Message_ExecuteStep] = &mp.ExecuteStepProcessor{}
	processors[m.Message_ExecutionEnding] = &mp.ExecutionEndingProcessor{}
	processors[m.Message_StepExecutionEnding] = &mp.StepExecutionEndingProcessor{}
	processors[m.Message_ScenarioExecutionEnding] = &mp.ScenarioExecutionEndingProcessor{}
	processors[m.Message_SpecExecutionEnding] = &mp.SpecExecutionEndingProcessor{}
	processors[m.Message_SpecDataStoreInit] = &mp.SpecDataStoreInitProcessor{}
	processors[m.Message_ScenarioDataStoreInit] = &mp.ScenarioDataStoreInitProcessor{}

	t.CustomScreenShot = &CustomScreenshotFn
}

// BeforeSuite hook is executed before suite execution begins
// This can be used for any setup before execution begins.
func BeforeSuite(fn func(), tags []string, op t.Operator) bool {
	hook := t.Hook{
		Type:     t.BEFORESUITE,
		Impl:     fn,
		Tags:     tags,
		Operator: op,
	}

	context.Hooks = append(context.Hooks, hook)
	return true
}

// AfterSuite hook is executed after hook execution is completed.
// This can be used for any cleanup after entire execution.
func AfterSuite(fn func(), tags []string, op t.Operator) bool {
	hook := t.Hook{
		Type:     t.AFTERSUITE,
		Impl:     fn,
		Tags:     tags,
		Operator: op,
	}

	context.Hooks = append(context.Hooks, hook)
	return true
}

// BeforeSpec hook is executed before every spec execution begins
// This can be used for any spec setup before execution begins. It can be executed for only a specific set of specs using Tags.
// This hook will be executed only for specs which satisfy the tags and tag operator(AND or OR) mentioned.
func BeforeSpec(fn func(), tags []string, op t.Operator) bool {
	hook := t.Hook{
		Type:     t.BEFORESPEC,
		Impl:     fn,
		Tags:     tags,
		Operator: op,
	}

	context.Hooks = append(context.Hooks, hook)
	return true
}

// AfterSpec hook is executed after every spec execution.
// This can be used for any spec cleanup after execution. It can be executed for only a specific set of specs using Tags.
// This hook will be executed only for specs which satisfy the tags and tag operator(AND or OR) mentioned.
func AfterSpec(fn func(), tags []string, op t.Operator) bool {
	hook := t.Hook{
		Type:     t.AFTERSPEC,
		Impl:     fn,
		Tags:     tags,
		Operator: op,
	}

	context.Hooks = append(context.Hooks, hook)
	return true
}

// BeforeScenario hook is executed before every scenario execution begins.
// This can be used for any scenario setup before execution begins. It can be executed for only a specific set of specs using Tags.
// This hook will be executed only for scenario which satisfy the tags and tag operator(AND or OR) mentioned.
func BeforeScenario(fn func(), tags []string, op t.Operator) bool {
	hook := t.Hook{
		Type:     t.BEFORESCENARIO,
		Impl:     fn,
		Tags:     tags,
		Operator: op,
	}

	context.Hooks = append(context.Hooks, hook)
	return true
}

// AfterScenario hook is executed after every scenario execution.
// This can be used for any spec cleanup after execution. It can be executed for only a specific set of specs using Tags.
// This hook will be executed only for scenarios which satisfy the tags and tag operator(AND or OR) mentioned.
func AfterScenario(fn func(), tags []string, op t.Operator) bool {
	hook := t.Hook{
		Type:     t.AFTERSCENARIO,
		Impl:     fn,
		Tags:     tags,
		Operator: op,
	}

	context.Hooks = append(context.Hooks, hook)
	return true
}

// BeforeStep hook is executed before every step execution begins.
// This can be used for any step setup before execution begins.
func BeforeStep(fn func(), tags []string, op t.Operator) bool {
	hook := t.Hook{
		Type:     t.BEFORESTEP,
		Impl:     fn,
		Tags:     tags,
		Operator: op,
	}

	context.Hooks = append(context.Hooks, hook)
	return true
}

// AfterStep hook is executed before after every step execution.
// This can be used for any step cleanup after execution.
func AfterStep(fn func(), tags []string, op t.Operator) bool {
	hook := t.Hook{
		Type:     t.AFTERSTEP,
		Impl:     fn,
		Tags:     tags,
		Operator: op,
	}

	context.Hooks = append(context.Hooks, hook)
	return true
}

// GetSuiteStore returns the suite store which keeps values added to it during the lifecycle of entire suite execution.
// Values are cleared after entire suite execution.
func GetSuiteStore() map[string]interface{} {
	return context.SuiteStore

}

// GetSpecStore returns the spec data store which keeps values added to it during the lifecycle of the specification execution.
// Values are cleared after every specification executes.
func GetSpecStore() map[string]interface{} {
	return context.SpecStore
}

// GetScenarioStore returns the scenario data store which keeps values added to it during the lifecycle of the scenario execution.
// Values are cleared after every scenario executes.
func GetScenarioStore() map[string]interface{} {
	return context.ScenarioStore
}

// WriteMessage adds additional information at exec time to be available on reports
func WriteMessage(message string, args ...interface{}) {
	context.CustomMessageRegistry = append(context.CustomMessageRegistry, fmt.Sprintf(message, args...))
}

// Run opens a port for listening to gauge. For internal use only.
func Run() {
	// fmt.Printf("We have got %d step implementations\n", len(context.Steps)) // move to logger

	var gaugePort = os.Getenv(c.GaugePortVariable)

	// fmt.Println("Connecting port:", gaugePort) // move to logger
	conn, err := net.Dial("tcp", net.JoinHostPort("127.0.0.1", gaugePort))
	if err != nil {
		fmt.Println("dial error:", err)
		return
	}
	defer conn.Close()
	for {
		msg, err := mu.ReadMessage(conn)
		if err != nil {
			fmt.Println("Error reading message : ", err)
			return
		}
		if msg.MessageType == m.Message_KillProcessRequest {
			return
		}

		processor := processors[msg.MessageType]

		if processor == nil {
			fmt.Println("Unable to find processor for message type : ", msg.MessageType)
			return
		}
		msgToSend := processor.Process(msg, context)

		err = mu.WriteGaugeMessage(msgToSend, conn)
		if err != nil {
			fmt.Println("Unable to write response : ", err.Error())
			return
		}
	}
}

// CustomScreenshotFn to set custom screenshot
// Returns a byte array which will be set as screenshot in case of failures
var CustomScreenshotFn func() []byte

// Step is an executable component of a specification. This function registers a step with given step description/step text and its implementation.
func Step(stepDesc string, impl interface{}) bool {
	desc, noOfArgs := parseDesc(stepDesc)
	implType := reflect.TypeOf(impl)

	if reflect.ValueOf(impl).Kind() != reflect.Func {
		//TODO decide whether to ignore or fail test
		fmt.Printf("Expected a function implementation for '%s' but got type '%s' - Ignoring test\n", stepDesc, implType.String())
		return false
	}

	//TODO validate not just the number of arguments but method signature
	if implType.NumIn() != noOfArgs {
		//TODO decide whether to ignore or fail test
		fmt.Printf("Mismatch in number of arguments in implementation of '%s' expected : %d, actual : %d - Ignoring test\n", desc, noOfArgs, implType.NumIn())
		return false
	}
	step := t.Step{
		Description: desc,
		Impl:        impl,
	}
	context.Steps = append(context.Steps, step)
	return true
}

func parseDesc(desc string) (string, int) {
	re := regexp.MustCompile("<(.*?)>")
	return re.ReplaceAllLiteralString(desc, "{}"), len(re.FindAllString(desc, -1))
}
