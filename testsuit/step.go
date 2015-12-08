package testsuit

type Step struct {
	Description string
	Impl        func(...interface{})
}

func (step *Step) Execute(args ...interface{}) {
	step.Impl(args...)
}
