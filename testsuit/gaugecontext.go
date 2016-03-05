package testsuit

type GaugeContext struct {
	Steps []Step
	Hooks []Hook
}

func (c *GaugeContext) GetStepByDesc(desc string) *Step {
	for _, step := range c.Steps {
		if step.Description == desc {
			return &step
		}
	}
	return nil
}

func (c *GaugeContext) GetHooks(hookType HookType, tags []string, op Operator) []Hook {
	return nil
}
