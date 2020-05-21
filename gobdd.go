package gobdd

import (
	"bufio"
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

	"github.com/cucumber/gherkin-go/v13"
	msgs "github.com/cucumber/messages-go/v12"
	"github.com/go-bdd/gobdd/context"
)

// Holds all the information about the suite (options, steps to execute etc)
type Suite struct {
	t             TestingT
	steps         []stepDef
	options       SuiteOptions
	hasStepErrors bool
}

// Holds all the information about how the suite or features/steps should be configured
type SuiteOptions struct {
	featuresPaths  string
	ignoreTags     []string
	tags           []string
	beforeScenario []func(ctx context.Context)
	afterScenario  []func(ctx context.Context)
	runInParallel  bool
}

// creates a new suite configuration with default values
func NewSuiteOptions() SuiteOptions {
	return SuiteOptions{
		featuresPaths:  "features/*.feature",
		ignoreTags:     []string{},
		tags:           []string{},
		beforeScenario: []func(ctx context.Context){},
		afterScenario:  []func(ctx context.Context){},
	}
}

// RunInParallel runs tests in parallel
func RunInParallel() func(*SuiteOptions) {
	return func(options *SuiteOptions) {
		options.runInParallel = true
	}
}

// WithFeaturesPath configures a pattern (regexp) where feature can be found.
// The default value is "features/*.feature"
func WithFeaturesPath(path string) func(*SuiteOptions) {
	return func(options *SuiteOptions) {
		options.featuresPaths = path
	}
}

// WithTags configures which tags should be skipped while executing a suite
// Every tag has to start with @
func WithTags(tags []string) func(*SuiteOptions) {
	return func(options *SuiteOptions) {
		options.tags = tags
	}
}

// WithBeforeScenario configures functions that should be executed before every scenario
func WithBeforeScenario(f func(ctx context.Context)) func(*SuiteOptions) {
	return func(options *SuiteOptions) {
		options.beforeScenario = append(options.beforeScenario, f)
	}
}

// WithAfterScenario configures functions that should be executed after every scenario
func WithAfterScenario(f func(ctx context.Context)) func(*SuiteOptions) {
	return func(options *SuiteOptions) {
		options.afterScenario = append(options.afterScenario, f)
	}
}

// WithIgnoredTags configures which tags should be skipped while executing a suite
// Every tag has to start with @ otherwise will be ignored
func WithIgnoredTags(tags []string) func(*SuiteOptions) {
	return func(options *SuiteOptions) {
		options.ignoreTags = tags
	}
}

type stepDef struct {
	expr *regexp.Regexp
	f    interface{}
}

type StepTest interface {
	Log(...interface{})
	Logf(string, ...interface{})
	Fatal(...interface{})
	Fatalf(string, ...interface{})
	Errorf(string, ...interface{})
	Error(...interface{})
}

type TestingT interface {
	StepTest
	Parallel()
	Fail()
	Run(name string, f func(t *testing.T)) bool
}

// Creates a new suites with given configuration and empty steps defined
func NewSuite(t TestingT, optionClosures ...func(*SuiteOptions)) *Suite {
	options := NewSuiteOptions()

	for i := 0; i < len(optionClosures); i++ {
		optionClosures[i](&options)
	}

	return &Suite{
		t:       t,
		steps:   []stepDef{},
		options: options,
	}
}

// AddStep registers a step in the suite.
//
// The second parameter is the step function that gets executed
// when a step definition matches the provided regular expression.
//
// A step function can have any number of parameters (even zero),
// but it MUST accept a gobdd.StepTest and context.Context as the first parameters (if there is any):
//
// 	func myStepFunction(t gobdd.StepTest, ctx context.Context, first int, second int) {
// 	}
func (s *Suite) AddStep(expr string, step interface{}) {
	err := validateStepFunc(step)
	if err != nil {
		s.t.Errorf("the step function for step `%s` is incorrect: %w", expr, err)
		s.hasStepErrors = true

		return
	}

	compiled, err := regexp.Compile(expr)
	if err != nil {
		s.t.Errorf("the step function is incorrect: %w", err)
		s.hasStepErrors = true

		return
	}

	s.steps = append(s.steps, stepDef{
		expr: compiled,
		f:    step,
	})
}

// AddRegexStep registers a step in the suite.
//
// The second parameter is the step function that gets executed
// when a step definition matches the provided regular expression.
//
// A step function can have any number of parameters (even zero),
// but it MUST accept a gobdd.StepTest and context.Context as the first parameters (if there is any):
//
// 	func myStepFunction(t gobdd.StepTest, ctx context.Context, first int, second int) {
// 	}
func (s *Suite) AddRegexStep(expr *regexp.Regexp, step interface{}) {
	err := validateStepFunc(step)
	if err != nil {
		s.t.Errorf("the step function is incorrect: %w", err)
		s.hasStepErrors = true

		return
	}

	s.steps = append(s.steps, stepDef{
		expr: expr,
		f:    step,
	})
}

// Executes the suite with given options and defined steps
func (s *Suite) Run() {
	if s.hasStepErrors {
		s.t.Fatal("the test contains invalid step definitions")

		return
	}

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

	doc, err := gherkin.ParseGherkinDocument(fileIO, (&msgs.Incrementing{}).NewId)
	if err != nil {
		s.t.Fail()
		s.t.Errorf("error while loading document: %s\n", err)

		return fmt.Errorf("error while loading document: %s", err)
	}

	if doc.Feature == nil {
		return nil
	}

	return s.runFeature(doc.Feature)
}

func (s *Suite) runFeature(feature *msgs.GherkinDocument_Feature) error {
	log.SetOutput(ioutil.Discard)

	hasErrors := false

	s.t.Run(feature.Name, func(t *testing.T) {
		var bkgSteps *msgs.GherkinDocument_Feature_Background

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
			ctx := context.New()
			s.runScenario(ctx, scenario, bkgSteps, t)
		}
	})

	if hasErrors {
		return errors.New("the feature contains errors")
	}

	return nil
}

func (s *Suite) getOutlineStep(
	steps []*msgs.GherkinDocument_Feature_Step,
	examples []*msgs.GherkinDocument_Feature_Scenario_Examples) []*msgs.GherkinDocument_Feature_Step {
	stepsList := make([][]*msgs.GherkinDocument_Feature_Step, len(steps))

	for i, outlineStep := range steps {
		for _, example := range examples {
			stepsList[i] = append(stepsList[i], s.stepsFromExamples(outlineStep, example)...)
		}
	}

	var newSteps []*msgs.GherkinDocument_Feature_Step

	if len(stepsList) == 0 {
		return newSteps
	}

	for ei := range examples {
		for ci := range examples[ei].TableBody {
			for si := range steps {
				newSteps = append(newSteps, stepsList[si][ci])
			}
		}
	}

	return newSteps
}

func (s *Suite) stepsFromExamples(
	sourceStep *msgs.GherkinDocument_Feature_Step,
	example *msgs.GherkinDocument_Feature_Scenario_Examples) []*msgs.GherkinDocument_Feature_Step {
	steps := []*msgs.GherkinDocument_Feature_Step{}

	placeholders := example.GetTableHeader().GetCells()
	placeholdersValues := []string{}

	for _, placeholder := range placeholders {
		ph := "<" + placeholder.GetValue() + ">"
		placeholdersValues = append(placeholdersValues, ph)
	}

	text := sourceStep.GetText()

	for _, row := range example.GetTableBody() {
		// iterate over the cells and update the text
		stepText, expr := s.stepFromExample(text, row, placeholdersValues)

		// find step definition for the new step
		def, err := s.findStepDef(stepText)
		if err != nil {
			continue
		}

		// add the step to the list
		s.AddStep(expr, def.f)

		// clone a step
		step := &msgs.GherkinDocument_Feature_Step{
			Location: sourceStep.Location,
			Keyword:  sourceStep.Keyword,
			Text:     stepText,
			Argument: sourceStep.Argument,
		}

		steps = append(steps, step)
	}

	return steps
}

func (s *Suite) stepFromExample(
	stepName string,
	row *msgs.GherkinDocument_Feature_TableRow, placeholders []string) (string, string) {
	expr := stepName

	for i, ph := range placeholders {
		t := getRegexpForVar(row.Cells[i].Value)
		expr = strings.Replace(expr, ph, t, -1)
		stepName = strings.Replace(stepName, ph, row.Cells[i].Value, -1)
	}

	return stepName, expr
}

func (s *Suite) callBeforeScenarios(ctx context.Context) {
	for _, f := range s.options.beforeScenario {
		f(ctx)
	}
}

func (s *Suite) callAfterScenarios(ctx context.Context) {
	for _, f := range s.options.afterScenario {
		f(ctx)
	}
}

func (s *Suite) runScenario(ctx context.Context, scenario *msgs.GherkinDocument_Feature_Scenario,
	bkg *msgs.GherkinDocument_Feature_Background, t *testing.T) {
	s.callBeforeScenarios(ctx)

	defer s.callAfterScenarios(ctx)
	t.Run(scenario.Name, func(t *testing.T) {
		if bkg != nil {
			steps := s.getBackgroundSteps(bkg)
			s.runSteps(ctx, t, steps)
		}
		steps := scenario.Steps
		if examples := scenario.GetExamples(); len(examples) > 0 {
			c := ctx.Clone()
			steps = s.getOutlineStep(scenario.GetSteps(), examples)
			s.runSteps(c, t, steps)
		} else {
			c := ctx.Clone()
			s.runSteps(c, t, steps)
		}
	})
}

func (s *Suite) runSteps(ctx context.Context, t *testing.T, steps []*msgs.GherkinDocument_Feature_Step) {
	for _, step := range steps {
		s.runStep(ctx, t, step)
	}
}

func (s *Suite) runStep(ctx context.Context, t *testing.T, step *msgs.GherkinDocument_Feature_Step) {
	defer func() {
		if r := recover(); r != nil {
			t.Error(r)
		}
	}()

	def, err := s.findStepDef(step.Text)
	if err != nil {
		t.Errorf("cannot find step definition for step: %s%s", step.Keyword, step.Text)

		return
	}

	params := def.expr.FindSubmatch([]byte(step.Text))[1:]
	t.Run(step.Text, func(t *testing.T) {
		def.run(ctx, t, params)
	})
}

func (def *stepDef) run(ctx context.Context, t TestingT, params [][]byte) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("%+v", r)
		}
	}()

	d := reflect.ValueOf(def.f)
	if len(params)+2 != d.Type().NumIn() {
		t.Errorf("the step function %s accepts %d arguments but %d received", d.String(), d.Type().NumIn(), len(params)+2)

		return
	}

	in := []reflect.Value{reflect.ValueOf(t), reflect.ValueOf(ctx)}

	for i, v := range params {
		if len(params) < i+1 {
			break
		}

		inType := d.Type().In(i + 2)
		paramType := paramType(v, inType)
		in = append(in, paramType)
	}

	d.Call(in)
}

func paramType(param []byte, inType reflect.Type) reflect.Value {
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

func (s *Suite) skipScenario(scenarioTags []*msgs.GherkinDocument_Feature_Tag) bool {
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

func (s *Suite) getBackgroundSteps(bkg *msgs.GherkinDocument_Feature_Background) []*msgs.GherkinDocument_Feature_Step {
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
