package reporter

import (
	"fmt"

	"github.com/cucumber/gherkin-go"
	. "github.com/logrusorgru/aurora"
)

type fmtReporter struct {
}

func (r fmtReporter) Step(step *gherkin.Step) {
	fmt.Printf("   %s%s\n", Yellow(step.Keyword), Green(step.Text))
}

func (r fmtReporter) Scenario(scenario *gherkin.Scenario) {
	if len(scenario.Tags) != 0 {
		for _, tag := range scenario.Tags {
			fmt.Printf(" %s\n", Yellow(tag.Name))
		}
	}
	fmt.Printf(" %s:%s\n", Yellow(scenario.Keyword), Green(scenario.Name))
}

func (r fmtReporter) ScenarioOutline(scenario *gherkin.ScenarioOutline) {
	if len(scenario.Tags) != 0 {
		for _, tag := range scenario.Tags {
			fmt.Printf(" %s\n", Yellow(tag.Name))
		}
	}
	fmt.Printf(" %s:%s\n", Yellow(scenario.Keyword), Green(scenario.Name))

	for _, example := range scenario.Examples {
		fmt.Printf(" %s:\n", Yellow(example.Keyword))
	}

	for _, step := range scenario.Steps {
		fmt.Printf("     %s%s\n", Yellow(step.Keyword), Green(step.Text))
	}
}

func (r fmtReporter) Undefined(s *gherkin.Step) {
	fmt.Printf("   %s%s\n", Red("Undefined step"), Yellow(s.Text))
}

func (r fmtReporter) Skip(s *gherkin.Step) {
	fmt.Printf(" %s%s\n", s.Keyword, Gray(10, s.Text))
}

func (r fmtReporter) Failed(step *gherkin.Step, e error) {
	fmt.Printf("   %s%s\n", Yellow(step.Keyword), Red(step.Text))
}

func (r fmtReporter) Background(background *gherkin.Background) {
	fmt.Printf("  %s%s\n", Yellow(background.Keyword), Green(background.Name))
	for _, s := range background.Steps {
		r.Step(s)
	}
}

func (r fmtReporter) Feature(feature *gherkin.Feature) {
	fmt.Printf("%s%s\n", Yellow(feature.Keyword), Green(feature.Name))
}

func (r fmtReporter) SkipScenario(s *gherkin.Scenario) {
	if len(s.Tags) != 0 {
		for _, tag := range s.Tags {
			fmt.Printf(" %s\n", Yellow(tag.Name))
		}
	}
	fmt.Printf(" %s%s\n", Gray(10, s.Keyword), Gray(10, s.Name))
	for _, step := range s.Steps {
		r.Skip(step)
	}
}

func (r fmtReporter) SkipScenarioOutline(s *gherkin.ScenarioOutline) {
	if len(s.Tags) != 0 {
		for _, tag := range s.Tags {
			fmt.Printf(" %s\n", Yellow(tag.Name))
		}
	}
	for _, example := range s.Examples {
		fmt.Printf(" %s:\n", Gray(10, example.Keyword))
	}

	fmt.Printf(" %s:%s\n", Gray(10, s.Keyword), Gray(10, s.Name))
	for _, step := range s.Steps {
		r.Skip(step)
	}
}

func (r fmtReporter) Report(report executionReport) {
	fmt.Printf("%d scenarios\n", report.countScenarios)
	fmt.Printf("%d steps\n", report.countSteps)
	if c := report.countFailedSteps; c > 0 {
		fmt.Printf("%d failed steps\n", report.countFailedSteps)
	}

	if c := report.countUndefined; c > 0 {
		fmt.Printf("%d undefined steps\n", report.countUndefined)
	}

	fmt.Printf("took %s\n", report.timeTook.String())
}
