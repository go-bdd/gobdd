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

// Holds all the information about the suite (options, steps to execute etc)
type Suite struct {
	t       *testing.T
	steps   []stepDef
	options SuiteOptions
}

// Holds all the information about how the suite or features/steps should be configured
type SuiteOptions struct {
	featuresPaths string
	ignoreTags    []string
}

// NewSuiteOptions creates a new suite configuration with default values
func NewSuiteOptions() SuiteOptions {
	return SuiteOptions{
		featuresPaths: "features/*.feature",
		ignoreTags:    []string{},
	}
}

// Configures a pattern (regexp) where feature can be found.
// The default value is "features/*.feature"
func (options SuiteOptions) WithFeaturesPath(path string) SuiteOptions {
	options.featuresPaths = path
	return options
}

// Configures which tags should be skipped while executing a suite
// Every tag has to start with @
func (options SuiteOptions) WithIgnoredTags(tags []string) SuiteOptions {
	options.ignoreTags = tags
	return options
}

// StepFunc every step function have to be compatible with this type
type StepFunc func(ctx Context) error

type stepDef struct {
	expr *regexp.Regexp
	f    StepFunc
}

// Creates a new suites with given configuration and empty steps defined
func NewSuite(t *testing.T, options SuiteOptions) *Suite {
	return &Suite{
		t:       t,
		steps:   []stepDef{},
		options: options,
	}
}

// Adds a step to the suite
// The first parameter is the step definition which can be:
//
// * string which will be converted to a regexp
// * [] byte which will be converted to a regexp as well
// * regexp
//
// No other types are supported
//
// The second parameter is a function which will be executed when while running a scenario one of the steps will match
// the given pattern
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

	suite.steps = append(suite.steps, stepDef{
		expr: regex,
		f:    f,
	})

	return nil
}

// Cxecutes the suite with given options and defined steps
func (suite *Suite) Run() {
	files, err := filepath.Glob(suite.options.featuresPaths)
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
			if suite.skipScenario(scenario.Tags, suite.options.ignoreTags) {
				continue
			}
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

func (suite *Suite) findStepDef(text string) (stepDef, error) {
	for _, step := range suite.steps {
		if !step.expr.MatchString(text) {
			continue
		}

		return step, nil
	}

	return stepDef{}, errors.New("cannot find step definition")
}

func (suite *Suite) skipScenario(scenarioTags []*gherkin.Tag, ignoreTags []string) bool {
	for _, tag := range scenarioTags {
		if contains(ignoreTags, tag.Name) {
			return true
		}
	}

	return false
}

// contains tells whether a contains x.
func contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}
