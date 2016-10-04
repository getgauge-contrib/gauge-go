package testsuit

import (
	"reflect"

	m "github.com/getgauge-contrib/gauge-go/gauge_messages"
)

type Step struct {
	Description string
	Impl        interface{}
}

// TODO: Set gauge messasges, screenshot, recoverableError
func (step *Step) Execute(args ...interface{}) *m.ProtoExecutionResult {
	fn := reflect.ValueOf(step.Impl)
	return executeFunc(fn, args...)
}
