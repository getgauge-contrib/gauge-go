package testsuit

import (
	"reflect"

	m "github.com/getgauge-contrib/gauge-go/gauge_messages"
)

type HookType int
type Operator int

const (
	BEFORESUITE    HookType = 1
	BEFORESPEC     HookType = 2
	BEFORESCENARIO HookType = 3
	BEFORESTEP     HookType = 4

	AFTERSUITE    HookType = 5
	AFTERSPEC     HookType = 6
	AFTERSCENARIO HookType = 7
	AFTERSTEP     HookType = 8
)

const (
	AND Operator = 1
	OR  Operator = 2
)

type Hook struct {
	Type     HookType
	Impl     func()
	Tags     []string
	Operator Operator
}

func (hook *Hook) Execute() *m.ProtoExecutionResult {
	fn := reflect.ValueOf(hook.Impl)
	return executeFunc(fn)
}
