package reporter

import (
	"fmt"
	"time"

	"github.com/cucumber/gherkin-go"
)

func NewReporter() *Reporter {
	return &Reporter{
		timer:          time.Now(),
		scenarios:      []interface{}{},
		undefinedSteps: []*gherkin.Step{},
		succeededSteps: []*gherkin.Step{},
		failedSteps:    []*gherkin.Step{},
		skippedSteps:   []*gherkin.Step{},
		backgrounds:    []*gherkin.Background{},
	}
}

type Reporter struct {
	timer            time.Time
	output           fmtReporter
	features         []*gherkin.Feature
	scenarios        []interface{} // *gherkin.Scenario or *gherkin.ScenarioOutline
	skippedScenarios []interface{} // *gherkin.Scenario or *gherkin.ScenarioOutline
	undefinedSteps   []*gherkin.Step
	succeededSteps   []*gherkin.Step
	failedSteps      []*gherkin.Step
	skippedSteps     []*gherkin.Step
	backgrounds      []*gherkin.Background
}

type executionReport struct {
	countScenarios   int
	countSteps       int
	countFailedSteps int
	countUndefined   int
	timeTook         time.Duration
}

func (r *Reporter) Report(reported interface{}) error {
	switch reported.(type) {
	case *gherkin.Step:
		s := reported.(*gherkin.Step)
		r.output.Step(s)

		r.succeededSteps = append(r.succeededSteps, s)
	case *gherkin.Scenario:
		s := reported.(*gherkin.Scenario)
		r.output.Scenario(s)
		r.scenarios = append(r.scenarios, s)
	case *gherkin.ScenarioOutline:
		s := reported.(*gherkin.ScenarioOutline)
		r.output.ScenarioOutline(s)
		r.scenarios = append(r.scenarios, s)

	case *gherkin.Feature:
		f := reported.(*gherkin.Feature)
		r.output.Feature(f)
		r.features = append(r.features, f)
	case *gherkin.Background:
		s := reported.(*gherkin.Background)
		r.output.Background(s)

		r.backgrounds = append(r.backgrounds, s)
	default:
		return fmt.Errorf("%T", reported)
	}

	return nil
}

func (r *Reporter) Undefined(s *gherkin.Step) {
	r.output.Undefined(s)
	r.undefinedSteps = append(r.undefinedSteps, s)
}

func (r *Reporter) Skip(reported interface{}) {
	switch reported.(type) {
	case *gherkin.Step:
		s := reported.(*gherkin.Step)
		r.output.Skip(s)

		r.skippedSteps = append(r.skippedSteps, s)
	case *gherkin.Scenario:
		s := reported.(*gherkin.Scenario)
		r.output.SkipScenario(s)
		r.skippedScenarios = append(r.skippedScenarios, s)
	case *gherkin.ScenarioOutline:
		s := reported.(*gherkin.ScenarioOutline)
		r.output.SkipScenarioOutline(s)
		r.skippedScenarios = append(r.skippedScenarios, s)
	default:
		panic("should never happen")
	}
}

func (r *Reporter) Failed(s *gherkin.Step, err error) {
	r.output.Failed(s, err)
	r.failedSteps = append(r.failedSteps, s)
}

func (r Reporter) GenerateReport() {
	er := executionReport{
		countScenarios:   r.countScenarios(),
		countSteps:       r.countSteps(),
		countFailedSteps: r.countFailed(),
		countUndefined:   r.countUndefined(),
		timeTook:         time.Since(r.timer),
	}

	r.output.Report(er)
}

func (r Reporter) countScenarios() int {
	return len(r.scenarios)
}

func (r Reporter) countSteps() int {
	return len(r.succeededSteps)
}

func (r *Reporter) countUndefined() int {
	return len(r.undefinedSteps)
}

func (r *Reporter) countFailed() int {
	return len(r.failedSteps)
}

func (r *Reporter) countSkippedSteps() int {
	return len(r.skippedSteps)
}
