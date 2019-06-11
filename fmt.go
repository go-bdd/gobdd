package gobdd

import (
	"fmt"

	"github.com/cucumber/gherkin-go"
	. "github.com/logrusorgru/aurora"
)

func printScenario(scenario string) {
	fmt.Print(Green(fmt.Sprintf("  Scenario: %s\n", scenario)))
}

func printScenarioOutline(scenario string) {
	fmt.Print(Green(fmt.Sprintf("  Scenario Outline: %s\n", scenario)))
}

func printStep(step *gherkin.Step) {
	fmt.Printf("    %s %s\n", Bold(Blue(step.Keyword)), step.Text)
}

func printFeature(featureName string) {
	fmt.Print(Green(fmt.Sprintf("Feature: %s\n", featureName)))
}

func printErrorf(format string, params ...interface{}) {
	fmt.Print(Red(fmt.Sprintf(format, params...)))
}

func printError(err error) {
	fmt.Println(Red(err))
}
