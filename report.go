package gobdd

import (
	"github.com/cucumber/messages-go/v12"
)

type reporter struct {
	t StepTest
}

func (r reporter) Pickle(scenario *messages.GherkinDocument_Feature_Scenario) *messages.Pickle {
	return &messages.Pickle{
		Id: scenario.Id,
	}
}
