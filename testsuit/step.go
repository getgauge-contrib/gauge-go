package testsuit

import (
	"fmt"
	"reflect"
	"runtime/debug"
	"strings"
	"time"

	"github.com/golang/protobuf/proto"
	m "github.com/manuviswam/gauge-go/gauge_messages"
)

type Step struct {
	Description string
	Impl        interface{}
}

// TODO: Set gauge messasges, screenshot, recoverableError
func (step *Step) Execute(args ...interface{}) (res *m.ProtoExecutionResult) {
	fn := reflect.ValueOf(step.Impl)
	rargs := make([]reflect.Value, len(args))
	for i, a := range args {
		rargs[i] = reflect.ValueOf(a)
	}
	res = &m.ProtoExecutionResult{}
	start := time.Now()
	defer func() {
		if r := recover(); r != nil {
			res.Failed = proto.Bool(true)
			res.ExecutionTime = proto.Int64(time.Since(start).Nanoseconds())
			res.StackTrace = proto.String(strings.SplitN(string(debug.Stack()), "\n", 9)[8])
			res.ErrorMessage = proto.String(fmt.Sprintf("%s", r))
		}
	}()
	fn.Call(rargs)
	res.Failed = proto.Bool(false)
	res.ExecutionTime = proto.Int64(time.Since(start).Nanoseconds())
	return res
}
