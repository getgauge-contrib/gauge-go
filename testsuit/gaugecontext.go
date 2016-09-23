package testsuit

import (
	"errors"
	"fmt"
)

type GaugeContext struct {
	Steps                 []Step
	Hooks                 []Hook
	ScenarioStore         map[string]interface{}
	SpecStore             map[string]interface{}
	SuiteStore            map[string]interface{}
	CustomMessageRegistry []string
}

func (c *GaugeContext) GetStepByDesc(desc string) (*Step, error) {
	for _, step := range c.Steps {
		if step.Description == desc {
			return &step, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("No implementation found for : %s", desc))
}

func (c *GaugeContext) GetHooks(hookType HookType, tags []string) []Hook {
	filteredByType := filterByType(c.Hooks, hookType)
	h := make([]Hook, 0)
	//TODO complexity is O(n^3) optimize it
	for _, hook := range filteredByType {
		switch hook.Operator {
		case OR:
			if containsAny(tags, hook.Tags) {
				h = append(h, hook)
			}
		case AND:
			if containsAll(tags, hook.Tags) {
				h = append(h, hook)
			}
		}
	}
	return h
}

// ClearCustomMessages clears the CustomMessageRegistry
func (c *GaugeContext) ClearCustomMessages() {
	c.CustomMessageRegistry = make([]string, 0)
}

func filterByType(hooks []Hook, t HookType) []Hook {
	h := make([]Hook, 0)
	for _, hook := range hooks {
		if hook.Type == t {
			h = append(h, hook)
		}
	}
	return h
}

func containsAny(s []string, k []string) bool {
	for _, tag := range k {
		if contains(s, tag) {
			return true
		}
	}
	return false
}

func containsAll(s []string, k []string) bool {
	for _, tag := range k {
		if !contains(s, tag) {
			return false
		}
	}
	return true
}

func contains(s []string, k string) bool {
	for _, a := range s {
		if a == k {
			return true
		}
	}
	return false
}
