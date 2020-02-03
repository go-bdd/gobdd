package gobdd

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"testing"

	messages "github.com/cucumber/cucumber-messages-go/v6"

	"github.com/cucumber/gherkin-go/v8"
)

// Holds all the information about the suite (options, steps to execute etc)
type Suite struct {
	t       *testing.T
	steps   []stepDef
	options SuiteOptions
}

// Holds all the information about how the suite or features/steps should be configured
type SuiteOptions struct {
	featuresPaths  string
	ignoreTags     []string
	tags           []string
	beforeScenario []func()
	afterScenario  []func()
	runInParallel  bool
}

// creates a new suite configuration with default values
func NewSuiteOptions() SuiteOptions {
	return SuiteOptions{
		featuresPaths:  "features/*.feature",
		ignoreTags:     []string{},
		tags:           []string{},
		beforeScenario: []func(){},
		afterScenario:  []func(){},
	}
}

// runs tests in parallel
func (options SuiteOptions) RunInParallel() SuiteOptions {
	options.runInParallel = true
	return options
}

// Configures a pattern (regexp) where feature can be found.
// The default value is "features/*.feature"
func (options SuiteOptions) WithFeaturesPath(path string) SuiteOptions {
	options.featuresPaths = path
	return options
}

// Configures which tags should be skipped while executing a suite
// Every tag has to start with @
func (options SuiteOptions) WithTags(tags []string) SuiteOptions {
	options.tags = tags
	return options
}

// Configures functions that should be executed before every scenario
func (options SuiteOptions) WithBeforeScenario(f func()) SuiteOptions {
	options.beforeScenario = append(options.beforeScenario, f)
	return options
}

// Configures functions that should be executed after every scenario
func (options SuiteOptions) WithAfterScenario(f func()) SuiteOptions {
	options.afterScenario = append(options.afterScenario, f)
	return options
}

// Configures which tags should be skipped while executing a suite
// Every tag has to start with @ otherwise will be ignored
func (options SuiteOptions) WithIgnoredTags(tags []string) SuiteOptions {
	options.ignoreTags = tags
	return options
}

type stepDef struct {
	expr *regexp.Regexp
	f    interface{}
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
func (s *Suite) AddStep(step interface{}, f interface{}) error {
	err := validateStepFunc(f)
	if err != nil {
		return err
	}
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

	s.steps = append(s.steps, stepDef{
		expr: regex,
		f:    f,
	})

	return nil
}

// Executes the suite with given options and defined steps
func (s *Suite) Run() {
	files, err := filepath.Glob(s.options.featuresPaths)
	if err != nil {
		s.t.Fatalf("cannot find features/ directory")
	}

	if s.options.runInParallel {
		s.t.Parallel()
	}

	for _, file := range files {
		err = s.executeFeature(file)
		if err != nil {
			s.t.Fail()
		}
	}
}

func (s *Suite) executeFeature(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("cannot open file %s", file)
	}
	defer f.Close()
	fileIO := bufio.NewReader(f)
	doc, err := gherkin.ParseGherkinDocument(fileIO)
	if err != nil {
		s.t.Fail()
		printErrorf("error while loading document: %s\n", err)
		return fmt.Errorf("error while loading document: %s\n", err)
	}

	if doc.Feature == nil {
		return nil
	}

	return s.runFeature(doc.Feature)
}

func (s *Suite) runFeature(feature *messages.GherkinDocument_Feature) error {
	log.SetOutput(ioutil.Discard)
	hasErrors := false

	s.t.Run(feature.Name, func(t *testing.T) {
		var bkgSteps *messages.GherkinDocument_Feature_Background

		for _, child := range feature.Children {
			if child.GetBackground() != nil {
				bkgSteps = child.GetBackground()
			}
			scenario := child.GetScenario()
			if scenario == nil {
				continue
			}

			if s.skipScenario(scenario.GetTags()) {
				t.Log(fmt.Sprintf("Skipping scenario %s", scenario.Name))
				continue
			}
			ctx := context.Background()
			err := s.runScenario(ctx, scenario, bkgSteps, t)
			if err != nil {
				hasErrors = true
			}
		}
	})

	if hasErrors {
		return errors.New("the feature contains errors")
	}

	return nil
}

func (s *Suite) getOutlineStep(steps []*messages.GherkinDocument_Feature_Step, examples []*messages.GherkinDocument_Feature_Scenario_Examples) []*messages.GherkinDocument_Feature_Step {
	var newSteps []*messages.GherkinDocument_Feature_Step
	for _, outlineStep := range steps {
		for _, example := range examples {
			newSteps = append(newSteps, s.stepsFromExample(outlineStep, example)...)
		}
	}
	return newSteps
}

func (s *Suite) stepsFromExample(sourceStep *messages.GherkinDocument_Feature_Step, example *messages.GherkinDocument_Feature_Scenario_Examples) []*messages.GherkinDocument_Feature_Step {
	steps := []*messages.GherkinDocument_Feature_Step{}
	text := sourceStep.GetText()
	placeholders := example.GetTableHeader().GetCells()
	for i, placeholder := range placeholders {
		ph := "<" + placeholder.GetValue() + ">"
		index := strings.Index(text, ph)
		originalText := text
		if index == -1 {
			continue
		}

		for _, row := range example.GetTableBody() {
			value := row.GetCells()[i].GetValue()
			text = strings.Replace(text, "<"+placeholder.Value+">", value, -1)
			t := getRegexpForVar(value)

			def, err := s.findStepDef(originalText)
			if err != nil {
				continue
			}

			expr := strings.Replace(def.expr.String(), ph, t, -1)
			err = s.AddStep(expr, def.f)

			if err != nil {
				continue
			}
		}
	}

	// clone a step
	step := &messages.GherkinDocument_Feature_Step{
		Location: sourceStep.Location,
		Keyword:  sourceStep.Keyword,
		Text:     text,
		Argument: sourceStep.Argument,
	}

	steps = append(steps, step)

	return steps
}

func (s *Suite) callBeforeScenarios() {
	for _, f := range s.options.beforeScenario {
		f()
	}
}

func (s *Suite) callAfterScenarios() {
	for _, f := range s.options.afterScenario {
		f()
	}
}

func (s *Suite) runScenario(ctx context.Context, scenario *messages.GherkinDocument_Feature_Scenario, bkg *messages.GherkinDocument_Feature_Background, t *testing.T) error {
	s.callBeforeScenarios()

	defer s.callAfterScenarios()
	t.Run(scenario.Name, func(t *testing.T) {
		if bkg != nil {
			steps := s.getBackgroundSteps(bkg)
			ctx = s.runSteps(ctx, t, steps)
		}
		steps := scenario.Steps
		if examples := scenario.GetExamples(); len(examples) > 0 {
			steps = s.getOutlineStep(scenario.GetSteps(), examples)
		}
		s.runSteps(ctx, t, steps)
	})
	return nil
}

func (s *Suite) runSteps(ctx context.Context, t *testing.T, steps []*messages.GherkinDocument_Feature_Step) context.Context {
	for _, step := range steps {
		ctx = s.runStep(ctx, t, step)
	}

	return ctx
}

func (s *Suite) runStep(ctx context.Context, t *testing.T, step *messages.GherkinDocument_Feature_Step) context.Context {
	defer func() {
		if r := recover(); r != nil {
			t.Error(r)
		}
	}()

	def, err := s.findStepDef(step.Text)
	if err != nil {
		t.Errorf("cannot find step definition for step: %s %s", step.Keyword, step.Text)
		return ctx
	}

	params := def.expr.FindSubmatch([]byte(step.Text))[1:]
	t.Run(step.Text, func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("%+v", r)
			}
		}()

		d := reflect.ValueOf(def.f)
		in := []reflect.Value{reflect.ValueOf(ctx)}
		if len(params)+1 != d.Type().NumIn() {
			t.Errorf("the step function %s accepts %d arguments but %d received", d.String(), d.Type().NumIn(), len(params)+1)
			return
		}
		for i, v := range params {
			inType := d.Type().In(i + 1)
			paramType := s.paramType(v, inType)
			in = append(in, paramType)
		}
		retValues := d.Call(in)
		v := retValues[1]

		if !v.IsNil() {
			err = v.Interface().(error)
			t.Errorf(step.Keyword, step.Text, err)
		}

		ctx = retValues[0].Interface().(context.Context)
	})

	return ctx
}

func (s *Suite) paramType(param []byte, inType reflect.Type) reflect.Value {
	paramType := reflect.ValueOf(param)
	if inType.Kind() == reflect.String {
		paramType = reflect.ValueOf(string(paramType.Interface().([]uint8)))
	}
	if inType.Kind() == reflect.Int {
		s := paramType.Interface().([]uint8)
		p, _ := strconv.Atoi(string(s))
		paramType = reflect.ValueOf(p)
	}
	if inType.Kind() == reflect.Float32 {
		s := paramType.Interface().([]uint8)
		p, _ := strconv.ParseFloat(string(s), 32)
		paramType = reflect.ValueOf(float32(p))
	}
	if inType.Kind() == reflect.Float64 {
		s := paramType.Interface().([]uint8)
		p, _ := strconv.ParseFloat(string(s), 32)
		paramType = reflect.ValueOf(p)
	}
	return paramType
}

func (s *Suite) findStepDef(text string) (stepDef, error) {
	var sd stepDef
	found := 0

	for _, step := range s.steps {
		if !step.expr.MatchString(text) {
			continue
		}

		if l := len(step.expr.FindAll([]byte(text), -1)); l > found {
			found = l
			sd = step
		}
	}

	if reflect.DeepEqual(sd, stepDef{}) {
		return sd, errors.New("cannot find step definition")
	}

	return sd, nil
}

func (s *Suite) skipScenario(scenarioTags []*messages.GherkinDocument_Feature_Tag) bool {
	for _, tag := range scenarioTags {
		if contains(s.options.ignoreTags, tag.Name) {
			return true
		}
	}

	if len(s.options.tags) == 0 {
		return false
	}

	for _, tag := range scenarioTags {
		if contains(s.options.tags, tag.Name) {
			return false
		}
	}

	return true
}

func (s *Suite) getBackgroundSteps(bkg *messages.GherkinDocument_Feature_Background) []*messages.GherkinDocument_Feature_Step {
	return bkg.Steps
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

func getRegexpForVar(v interface{}) string {
	s := v.(string)

	if _, err := strconv.Atoi(s); err == nil {
		return "(\\d+)"
	}

	if _, err := strconv.ParseFloat(s, 32); err == nil {
		return "([+-]?([0-9]*[.])?[0-9]+)"
	}

	return "(.*)"
}
