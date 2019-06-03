package gobdd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/cucumber/gherkin-go"
)

type Suite struct {
	t       *testing.T
	steps   []StepDef
	options SuiteOptions
}

type SuiteOptions struct {
	featuresPaths string
}

func NewSuiteOptions() SuiteOptions {
	return SuiteOptions{
		featuresPaths: "features/*.feature",
	}
}

type StepFunc func(ctx Context) error

type StepDef struct {
	expr *regexp.Regexp
	f    StepFunc
}

func NewSuite(t *testing.T, options SuiteOptions) *Suite {
	return &Suite{
		t:       t,
		steps:   []StepDef{},
		options: options,
	}
}

func (suite *Suite) AddStep(step interface{}, f StepFunc) error {
	var regex *regexp.Regexp

	switch t := step.(type) {
	case *regexp.Regexp:
		regex = t
	case string:
		regex = regexp.MustCompile(t)
	case []byte:
		regex = regexp.MustCompile(string(t))
	default:
		return fmt.Errorf("expecting expr to be a *regexp.Regexp or a string, got type: %T", step)
	}

	suite.steps = append(suite.steps, StepDef{
		expr: regex,
		f:    f,
	})

	return nil
}

func (suite *Suite) Run() {
	files, err := filepath.Glob("features/*.feature")
	if err != nil {
		suite.t.Fatalf("cannot find features/ directory")
	}

	for _, file := range files {
		err = suite.executeFeature(file)
		if err != nil {
			suite.t.Fail()
			printError(err)
		}
	}
}

func (suite *Suite) executeFeature(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("cannot open file %s", file)
	}
	defer f.Close()
	fileIO := bufio.NewReader(f)
	doc, err := gherkin.ParseGherkinDocument(fileIO)
	if err != nil {
		suite.t.Fail()
		printErrorf("error while loading document: %s", err)
	}

	if doc.Feature == nil {
		return nil
	}

	return suite.runFeature(doc.Feature)
}

func (suite *Suite) runFeature(feature *gherkin.Feature) error {
	printFeature(feature.Name)
	hasErrors := false

	for _, s := range feature.Children {
		scenario, ok := s.(*gherkin.Scenario)
		if ok {
			err := suite.runScenario(scenario)
			if err != nil {
				hasErrors = true
			}
		}
	}

	if hasErrors {
		return errors.New("the feature contains errors")
	}

	return nil
}

func (suite *Suite) runScenario(scenario *gherkin.Scenario) error {
	printScenario(scenario.Name)
	ctx := newContext()

	for _, step := range scenario.Steps {
		printStep(step)
		def, err := suite.findStepDef(step.Text)
		if err != nil {
			suite.t.Errorf("cannot find step definition for step '%s'", step.Text)
			continue
		}

		b := def.expr.FindSubmatch([]byte(step.Text))
		ctx.setParams(b[1:])

		err = def.f(ctx)
		if err != nil {
			printError(err)
			suite.t.Fail()
		}
	}

	return nil
}

func (suite *Suite) findStepDef(text string) (StepDef, error) {
	for _, step := range suite.steps {
		if !step.expr.MatchString(text) {
			continue
		}

		return step, nil
	}

	return StepDef{}, errors.New("cannot find step definition")
}
