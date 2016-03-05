package testsuit

import (
	"reflect"
)

type Step struct {
	Description string
	Impl        interface{}
}

func (step *Step) Execute(args ...interface{}) {
	fn := reflect.ValueOf(step.Impl)
	rargs := make([]reflect.Value, len(args))
	for i, a := range args {
		rargs[i] = reflect.ValueOf(a)
	}
	fn.Call(rargs)
}
