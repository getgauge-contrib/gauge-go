package gauge

import (
	"fmt"
	"reflect"

	t "github.com/manuviswam/gauge-go/testsuit"
)

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
