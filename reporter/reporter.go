package reporter

import "github.com/cucumber/gherkin-go"

type Reporter interface {
	Scenario(scenario *gherkin.Scenario)
	UndefinedStep(step *gherkin.Step)
	SucceededStep(step *gherkin.Step)
	FailedStep(step *gherkin.Step, err error)
	ScenarioOutline(outline *gherkin.ScenarioOutline)
	Background(bkg *gherkin.Background)
}
