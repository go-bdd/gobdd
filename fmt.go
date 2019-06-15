package gobdd

import (
	"fmt"
	"time"

	"github.com/cucumber/gherkin-go"
	. "github.com/logrusorgru/aurora"
)

type stepResult struct {
	// typ     stepType
	// feature *feature
	// owner   interface{}
	step *gherkin.Step
	def  *stepDef
	err  error
}

type reporter struct {
	started time.Time

	features  []*gherkin.Feature
	scenarios []interface{}
	failed    []*stepResult
	passed    []*stepResult
	skipped   []*stepResult
	undefined []*stepResult
	pending   []*stepResult
}

func (r *reporter) start() {
	r.started = time.Now()
}

func (r *reporter) Feature(feature *gherkin.Feature) {
	r.features = append(r.features, feature)
	fmt.Print(Green(fmt.Sprintf("Feature: %s\n", feature.Name)))
}

func (r *reporter) Scenario(scenario *gherkin.Scenario) {
	r.scenarios = append(r.scenarios, scenario)
	fmt.Print(Green(fmt.Sprintf("  Scenario: %s\n", scenario.Name)))
}

func (r *reporter) ScenarioOutline(outline *gherkin.ScenarioOutline) {
	r.scenarios = append(r.scenarios, outline)
	fmt.Print(Green(fmt.Sprintf("  Scenario Outline: %s\n", outline.Name)))
}

func printStep(step *gherkin.Step) {
	fmt.Printf("    %s %s\n", Bold(Blue(step.Keyword)), step.Text)
}

func printErrorf(format string, params ...interface{}) {
	fmt.Print(Red(fmt.Sprintf(format, params...)))
}

func printError(err error) {
	fmt.Println(Red(err))
}
