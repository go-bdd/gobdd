package reporter

import (
	"fmt"
	"time"

	"github.com/cucumber/gherkin-go"
	. "github.com/logrusorgru/aurora"
)

func NewFmt() *fmtReporter {
	return &fmtReporter{
		timer:          time.Now(),
		scenarios:      []interface{}{},
		undefinedSteps: []*gherkin.Step{},
		succeededSteps: []*gherkin.Step{},
		failedSteps:    []*gherkin.Step{},
	}
}

type fmtReporter struct {
	timer                   time.Time
	scenarios               []interface{} // *gherkin.Scenario or *gherkin.ScenarioOutline
	skippedScenarios        []*gherkin.Scenario
	skippedScenarioOutlines []*gherkin.ScenarioOutline
	undefinedSteps          []*gherkin.Step
	succeededSteps          []*gherkin.Step
	failedSteps             []*gherkin.Step
}

func (r *fmtReporter) Scenario(scenario *gherkin.Scenario) {
	if len(scenario.Tags) != 0 {
		for _, tag := range scenario.Tags {
			fmt.Printf("%s\n", Yellow(tag.Name))
		}
	}

	fmt.Printf("%s: %s\n", scenario.Keyword, Gray(10, scenario.Name))
	r.scenarios = append(r.scenarios, scenario)
}

func (r *fmtReporter) SkippedScenario(scenario *gherkin.Scenario) {
	if len(scenario.Tags) != 0 {
		for _, tag := range scenario.Tags {
			fmt.Printf("%s\n", Yellow(tag.Name))
		}
	}

	fmt.Printf("%s: %s\n", scenario.Keyword, Gray(10, scenario.Name))
	r.skippedScenarios = append(r.skippedScenarios, scenario)
}

func (r *fmtReporter) SkippedScenarioOutline(scenario *gherkin.ScenarioOutline) {
	if len(scenario.Tags) != 0 {
		for _, tag := range scenario.Tags {
			fmt.Printf("%s\n", Yellow(tag.Name))
		}
	}

	fmt.Printf("%s: %s\n", scenario.Keyword, Blue(scenario.Name))
	r.skippedScenarioOutlines = append(r.skippedScenarioOutlines, scenario)
}

func (r *fmtReporter) Background(bkg *gherkin.Background) {
	fmt.Printf("%s: %s\n", bkg.Keyword, Green(bkg.Name))
}

func (r *fmtReporter) UndefinedStep(step *gherkin.Step) {
	fmt.Printf("  %s: %s\n", Red("Undefined step"), Yellow(step.Text))
	r.undefinedSteps = append(r.undefinedSteps, step)
}

func (r *fmtReporter) SucceededStep(step *gherkin.Step) {
	fmt.Printf("  %s: %s\n", Yellow(step.Keyword), Green(step.Text))
	r.succeededSteps = append(r.succeededSteps, step)
}

func (r *fmtReporter) FailedStep(step *gherkin.Step, err error) {
	fmt.Printf("  %s: %s\n", step.Keyword, Red(step.Text))
	r.failedSteps = append(r.failedSteps, step)
}

func (r *fmtReporter) ScenarioOutline(outline *gherkin.ScenarioOutline) {
	if len(outline.Tags) != 0 {
		for _, tag := range outline.Tags {
			fmt.Printf("%s\n", Yellow(tag.Name))
		}
	}
	fmt.Printf("%s: %s\n", outline.Keyword, Green(outline.Name))
	r.scenarios = append(r.scenarios, outline)
}

func (r fmtReporter) GenerateReport() {
	fmt.Printf("%d scenarios\n", len(r.scenarios))
	fmt.Printf("%d steps\n", len(r.succeededSteps))
	if c := len(r.failedSteps); c > 0 {
		fmt.Printf("%d failed steps\n", len(r.failedSteps))
	}

	if c := len(r.undefinedSteps); c > 0 {
		fmt.Printf("%d undefined steps\n", len(r.undefinedSteps))
	}

	fmt.Printf("took %s\n", time.Since(r.timer).String())
}
