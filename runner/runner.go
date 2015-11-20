package runner

import "fmt"

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
}
