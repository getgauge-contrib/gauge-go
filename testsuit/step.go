package testsuit

type Step struct {
	Description string
	Impl func()
}

func (step *Step)Execute()  {
	step.Impl()
}
