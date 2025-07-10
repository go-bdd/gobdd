package gobdd

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"testing"

	gherkin "github.com/cucumber/gherkin/go/v33"
	msgs "github.com/cucumber/messages/go/v28"
)

const contextArgumentsNumber = 2

// Suite holds all the information about the suite (options, steps to execute etc)
type Suite struct {
	t              TestingT
	steps          []stepDef
	options        SuiteOptions
	hasStepErrors  bool
	parameterTypes map[string][]string
}

// SuiteOptions holds all the information about how the suite or features/steps should be configured
type SuiteOptions struct {
	featureSource  featureSource
	ignoreTags     []string
	tags           []string
	beforeScenario []func(ctx Context)
	afterScenario  []func(ctx Context)
	beforeStep     []func(ctx Context)
	afterStep      []func(ctx Context)
	runInParallel  bool
}

type featureSource interface {
	loadFeatures() ([]feature, error)
}

type feature interface {
	Open() (io.Reader, error)
}

type pathFeatureSource string

func (s pathFeatureSource) loadFeatures() ([]feature, error) {
	files, err := filepath.Glob(string(s))
	if err != nil {
		return nil, errors.New("cannot find features/ directory")
	}

	features := make([]feature, 0, len(files))

	for _, f := range files {
		features = append(features, fileFeature(f))
	}

	return features, nil
}

type fileFeature string

func (f fileFeature) Open() (io.Reader, error) {
	file, err := os.Open(string(f))
	if err != nil {
		return nil, fmt.Errorf("cannot open file %s", f)
	}

	return file, nil
}

// NewSuiteOptions creates a new suite configuration with default values
func NewSuiteOptions() SuiteOptions {
	return SuiteOptions{
		featureSource:  pathFeatureSource("features/*.feature"),
		ignoreTags:     []string{},
		tags:           []string{},
		beforeScenario: []func(ctx Context){},
		afterScenario:  []func(ctx Context){},
		beforeStep:     []func(ctx Context){},
		afterStep:      []func(ctx Context){},
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
		options.featureSource = pathFeatureSource(path)
	}
}

// WithTags configures which tags should be skipped while executing a suite
// Every tag has to start with @
func WithTags(tags ...string) func(*SuiteOptions) {
	return func(options *SuiteOptions) {
		options.tags = tags
	}
}

// WithBeforeScenario configures functions that should be executed before every scenario
func WithBeforeScenario(f func(ctx Context)) func(*SuiteOptions) {
	return func(options *SuiteOptions) {
		options.beforeScenario = append(options.beforeScenario, f)
	}
}

// WithAfterScenario configures functions that should be executed after every scenario
func WithAfterScenario(f func(ctx Context)) func(*SuiteOptions) {
	return func(options *SuiteOptions) {
		options.afterScenario = append(options.afterScenario, f)
	}
}

// WithBeforeStep configures functions that should be executed before every step
func WithBeforeStep(f func(ctx Context)) func(*SuiteOptions) {
	return func(options *SuiteOptions) {
		options.beforeStep = append(options.beforeStep, f)
	}
}

// WithAfterStep configures functions that should be executed after every step
func WithAfterStep(f func(ctx Context)) func(*SuiteOptions) {
	return func(options *SuiteOptions) {
		options.afterStep = append(options.afterStep, f)
	}
}

// WithIgnoredTags configures which tags should be skipped while executing a suite
// Every tag has to start with @ otherwise will be ignored
func WithIgnoredTags(tags ...string) func(*SuiteOptions) {
	return func(options *SuiteOptions) {
		options.ignoreTags = tags
	}
}

type stepDef struct {
	expr *regexp.Regexp
	f    interface{}
}

type StepTest interface {
	testing.TB
}

type TestingT interface {
	StepTest
	Parallel()
	Run(name string, f func(t *testing.T)) bool
}

// TestingTKey is used to store reference to current *testing.T instance
type TestingTKey struct{}

// FeatureKey is used to store reference to current *msgs.Feature instance
type FeatureKey struct{}

// RuleKey is used to store reference to current *msgs.Rule instance
type RuleKey struct{}

// ScenarioKey is used to store reference to current *msgs.Scenario instance
type ScenarioKey struct{}

// Creates a new suites with given configuration and empty steps defined
func NewSuite(t TestingT, optionClosures ...func(*SuiteOptions)) *Suite {
	options := NewSuiteOptions()

	for i := 0; i < len(optionClosures); i++ {
		optionClosures[i](&options)
	}

	s := &Suite{
		t:              t,
		steps:          []stepDef{},
		options:        options,
		parameterTypes: map[string][]string{},
	}

	// see https://github.com/cucumber/cucumber-expressions/blob/main/go/parameter_type_registry.go
	s.AddParameterTypes(`{int}`, []string{`(-?\d+)`})
	s.AddParameterTypes(`{float}`, []string{`([-+]?\d*\.?\d+)`})
	s.AddParameterTypes(`{word}`, []string{`([^\s]+)`})
	s.AddParameterTypes(`{text}`, []string{`"([^"\\]*(?:\\.[^"\\]*)*)"`, `'([^'\\]*(?:\\.[^'\\]*)*)'`})

	return s
}

// AddParameterTypes adds a list of parameter types that will be used to simplify step definitions.
//
// The first argument is the parameter type and the second parameter is a list of regular expressions
// that should replace the parameter type.
//
//	s.AddParameterTypes(`{int}`, []string{`(\d)`})
//
// The regular expression should compile, otherwise will produce an error and stop executing.
func (s *Suite) AddParameterTypes(from string, to []string) {
	for _, to := range to {
		_, err := regexp.Compile(to)
		if err != nil {
			s.t.Fatalf(`the regular expression for key %s doesn't compile: %s`, from, to)
		}

		s.parameterTypes[from] = append(s.parameterTypes[from], to)
	}
}

// AddStep registers a step in the suite.
//
// The second parameter is the step function that gets executed
// when a step definition matches the provided regular expression.
//
// A step function can have any number of parameters (even zero),
// but it MUST accept a gobdd.StepTest and gobdd.Context as the first parameters (if there is any):
//
//	func myStepFunction(t gobdd.StepTest, ctx gobdd.Context, first int, second int) {
//	}
func (s *Suite) AddStep(expr string, step interface{}) {
	err := validateStepFunc(step)
	if err != nil {
		s.t.Errorf("the step function for step `%s` is incorrect: %s", expr, err.Error())
		s.hasStepErrors = true

		return
	}

	exprs := s.applyParameterTypes(expr)

	for _, expr := range exprs {
		compiled, err := regexp.Compile(expr)
		if err != nil {
			s.t.Errorf("the step function is incorrect: %s", err.Error())
			s.hasStepErrors = true

			return
		}

		s.steps = append(s.steps, stepDef{
			expr: compiled,
			f:    step,
		})
	}
}

func (s *Suite) applyParameterTypes(expr string) []string {
	exprs := []string{expr}

	for from, to := range s.parameterTypes {
		for _, t := range to {
			if strings.Contains(expr, from) {
				exprs = append(exprs, s.applyParameterTypes(strings.ReplaceAll(expr, from, t))...)
			}
		}
	}

	return exprs
}

// AddRegexStep registers a step in the suite.
//
// The second parameter is the step function that gets executed
// when a step definition matches the provided regular expression.
//
// A step function can have any number of parameters (even zero),
// but it MUST accept a gobdd.StepTest and gobdd.Context as the first parameters (if there is any):
//
//	func myStepFunction(t gobdd.StepTest, ctx gobdd.Context, first int, second int) {
//	}
func (s *Suite) AddRegexStep(expr *regexp.Regexp, step interface{}) {
	err := validateStepFunc(step)
	if err != nil {
		s.t.Errorf("the step function is incorrect: %s", err.Error())
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

	features, err := s.options.featureSource.loadFeatures()
	if err != nil {
		s.t.Fatal(err.Error())
	}

	if s.options.runInParallel {
		s.t.Parallel()
	}

	for _, feature := range features {
		err = s.executeFeature(feature)
		if err != nil {
			s.t.Fail()
		}
	}
}

func (s *Suite) executeFeature(feature feature) error {
	f, err := feature.Open()
	if err != nil {
		return err
	}

	if closer, ok := f.(io.Closer); ok {
		defer closer.Close()
	}

	featureIO := bufio.NewReader(f)

	doc, err := gherkin.ParseGherkinDocument(featureIO, (&msgs.Incrementing{}).NewId)
	if err != nil {
		s.t.Fatalf("error while loading document: %s\n", err)
	}

	if doc.Feature == nil {
		return nil
	}

	return s.runFeature(doc.Feature)
}

func (s *Suite) runFeature(feature *msgs.Feature) error {
	if s.shouldSkipFeatureOrRule(feature.Tags) {
		s.t.Logf("the feature (%s) is ignored ", feature.Name)
		return nil
	}

	hasErrors := false

	s.t.Run(fmt.Sprintf("%s %s", strings.TrimSpace(feature.Keyword), feature.Name), func(t *testing.T) {
		backgrounds := []*msgs.Background{}

		for _, child := range feature.Children {
			if child.Background != nil {
				backgrounds = append(backgrounds, child.Background)
			}

			if rule := child.Rule; rule != nil {
				s.runRule(feature, rule, backgrounds, t)
			}
			if scenario := child.Scenario; scenario != nil {
				ctx := NewContext()
				ctx.Set(FeatureKey{}, feature)
				s.runScenario(ctx, scenario, backgrounds, t, feature.Tags)
			}
		}
	})

	if hasErrors {
		return errors.New("the feature contains errors")
	}

	return nil
}

func (s *Suite) getOutlineStep(
	steps []*msgs.Step,
	examples []*msgs.Examples) []*msgs.Step {
	stepsList := make([][]*msgs.Step, len(steps))

	for i, outlineStep := range steps {
		for _, example := range examples {
			stepsList[i] = append(stepsList[i], s.stepsFromExamples(outlineStep, example)...)
		}
	}

	var newSteps []*msgs.Step

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
	sourceStep *msgs.Step,
	example *msgs.Examples) []*msgs.Step {
	steps := []*msgs.Step{}

	placeholdersValues := []string{}

	if example.TableHeader != nil {
		placeholders := example.TableHeader.Cells
		for _, placeholder := range placeholders {
			ph := "<" + placeholder.Value + ">"
			placeholdersValues = append(placeholdersValues, ph)
		}
	}

	text := sourceStep.Text

	for _, row := range example.TableBody {
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
		step := &msgs.Step{
			Location:    sourceStep.Location,
			Keyword:     sourceStep.Keyword,
			Text:        stepText,
			KeywordType: sourceStep.KeywordType,
			DocString:   sourceStep.DocString,
			DataTable:   sourceStep.DataTable,
			Id:          sourceStep.Id,
		}

		steps = append(steps, step)
	}

	return steps
}

func (s *Suite) stepFromExample(
	stepName string,
	row *msgs.TableRow, placeholders []string) (string, string) {
	expr := stepName

	for i, ph := range placeholders {
		t := getRegexpForVar(row.Cells[i].Value)
		expr = strings.ReplaceAll(expr, ph, t)
		stepName = strings.ReplaceAll(stepName, ph, row.Cells[i].Value)
	}

	return stepName, expr
}

func (s *Suite) callBeforeScenarios(ctx Context) {
	for _, f := range s.options.beforeScenario {
		f(ctx)
	}
}

func (s *Suite) callAfterScenarios(ctx Context) {
	for _, f := range s.options.afterScenario {
		f(ctx)
	}
}

func (s *Suite) callBeforeSteps(ctx Context) {
	for _, f := range s.options.beforeStep {
		f(ctx)
	}
}

func (s *Suite) callAfterSteps(ctx Context) {
	for _, f := range s.options.afterStep {
		f(ctx)
	}
}
func (s *Suite) runRule(feature *msgs.Feature, rule *msgs.Rule,
	backgrounds []*msgs.Background, t *testing.T) {
	ruleTags := feature.Tags
	ruleTags = append(ruleTags, rule.Tags...)

	if s.shouldSkipFeatureOrRule(ruleTags) {
		s.t.Logf("the rule (%s) is ignored ", feature.Name)
		return
	}

	ruleBackgrounds := []*msgs.Background{}
	ruleBackgrounds = append(ruleBackgrounds, backgrounds...)

	t.Run(fmt.Sprintf("%s %s", strings.TrimSpace(rule.Keyword), rule.Name), func(t *testing.T) {
		for _, ruleChild := range rule.Children {
			if ruleChild.Background != nil {
				ruleBackgrounds = append(ruleBackgrounds, ruleChild.Background)
			}
			if scenario := ruleChild.Scenario; scenario != nil {
				ctx := NewContext()
				ctx.Set(FeatureKey{}, feature)
				ctx.Set(RuleKey{}, rule)
				s.runScenario(ctx, scenario, ruleBackgrounds, t, ruleTags)
			}
		}
	})
}
func (s *Suite) runScenario(ctx Context, scenario *msgs.Scenario,
	backgrounds []*msgs.Background, t *testing.T, parentTags []*msgs.Tag) {
	if s.shouldSkipScenario(append(parentTags, scenario.Tags...)) {
		t.Logf("Skipping scenario %s", scenario.Name)
		return
	}

	t.Run(fmt.Sprintf("%s %s", strings.TrimSpace(scenario.Keyword), scenario.Name), func(t *testing.T) {
		// NOTE consider passing t as argument to scenario hooks
		ctx.Set(ScenarioKey{}, scenario)
		ctx.Set(TestingTKey{}, t)
		defer ctx.Set(TestingTKey{}, nil)

		s.callBeforeScenarios(ctx)
		defer s.callAfterScenarios(ctx)

		if len(backgrounds) > 0 {
			steps := s.getBackgroundSteps(backgrounds)
			s.runSteps(ctx, t, steps)
		}
		steps := scenario.Steps
		if examples := scenario.Examples; len(examples) > 0 {
			c := ctx.Clone()
			steps = s.getOutlineStep(scenario.Steps, examples)
			s.runSteps(c, t, steps)
		} else {
			c := ctx.Clone()
			s.runSteps(c, t, steps)
		}
	})
}

func (s *Suite) runSteps(ctx Context, t *testing.T, steps []*msgs.Step) {
	for _, step := range steps {
		s.runStep(ctx, t, step)
	}
}

func (s *Suite) runStep(ctx Context, t *testing.T, step *msgs.Step) {
	defer func() {
		if r := recover(); r != nil {
			t.Error(r)
		}
	}()

	def, err := s.findStepDef(step.Text)
	if err != nil {
		t.Fatalf("cannot find step definition for step: %s%s", step.Keyword, step.Text)
	}

	matches := def.expr.FindSubmatch([]byte(step.Text))[1:]
	params := make([]interface{}, 0, len(matches)) // defining the slices capacity instead of the length to use append
	for _, m := range matches {
		params = append(params, m)
	}

	if step.DocString != nil {
		params = append(params, step.DocString.Content)
	}
	if step.DataTable != nil {
		params = append(params, *step.DataTable)
	}

	t.Run(fmt.Sprintf("%s %s", strings.TrimSpace(step.Keyword), step.Text), func(t *testing.T) {
		// NOTE consider passing t as argument to step hooks
		ctx.Set(TestingTKey{}, t)
		defer ctx.Set(TestingTKey{}, nil)

		s.callBeforeSteps(ctx)
		defer s.callAfterSteps(ctx)

		def.run(ctx, t, params)
	})
}

func (def *stepDef) run(ctx Context, t TestingT, params []interface{}) { // nolint:interfacer
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("%+v", r)
		}
	}()

	d := reflect.ValueOf(def.f)
	if len(params)+contextArgumentsNumber != d.Type().NumIn() {
		t.Fatalf("the step function %s accepts %d arguments but %d received",
			d.String(),
			d.Type().NumIn(),
			len(params)+contextArgumentsNumber)

		return
	}

	in := []reflect.Value{reflect.ValueOf(t), reflect.ValueOf(ctx)}

	for i, v := range params {
		if len(params) < i+1 {
			break
		}

		inType := d.Type().In(i + contextArgumentsNumber)

		paramType, err := paramType(v, inType)
		if err != nil {
			t.Fatal(err)
		}

		in = append(in, paramType)
	}

	d.Call(in)
}

func paramType(param interface{}, inType reflect.Type) (reflect.Value, error) {
	switch inType.Kind() { // nolint:exhaustive // the linter does not recognize 'default:' to satisfy exhaustiveness
	case reflect.String:
		s, err := shouldBeString(param)
		return reflect.ValueOf(s), err
	case reflect.Int:
		v, err := shouldBeInt(param)
		return reflect.ValueOf(v), err
	case reflect.Float32:
		v, err := shouldBeFloat(param, 32) // nolint:mnd
		return reflect.ValueOf(float32(v)), err
	case reflect.Float64:
		v, err := shouldBeFloat(param, 64) // nolint:mnd
		return reflect.ValueOf(v), err
	case reflect.Slice:
		// only []byte is supported
		if inType != reflect.TypeOf([]byte(nil)) {
			return reflect.Value{}, fmt.Errorf("the slice argument type %s is not supported", inType.Kind())
		}

		v, err := shouldBeByteSlice(param)

		return reflect.ValueOf(v), err
	case reflect.Struct:
		// the only struct supported is the one introduced by cucumber
		if inType != reflect.TypeOf(msgs.DataTable{}) {
			return reflect.Value{}, fmt.Errorf("the struct argument type %s is not supported", inType.Kind())
		}

		v, err := shouldBeDataTable(param)

		return reflect.ValueOf(v), err
	default:
		return reflect.Value{}, fmt.Errorf("the type %s is not supported", inType.Kind())
	}
}

func shouldBeDataTable(input interface{}) (msgs.DataTable, error) {
	if v, ok := input.(msgs.DataTable); ok {
		return v, nil
	}

	return msgs.DataTable{}, fmt.Errorf("cannot convert %v of type %T to messages.DataTable", input, input)
}

func shouldBeByteSlice(input interface{}) ([]byte, error) {
	if v, ok := input.([]byte); ok {
		return v, nil
	}

	return nil, fmt.Errorf("cannot convert %v of type %T to []byte", input, input)
}

func shouldBeInt(input interface{}) (int, error) {
	s, err := shouldBeString(input)
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(s)
}

func shouldBeFloat(input interface{}, bitSize int) (float64, error) {
	s, err := shouldBeString(input)
	if err != nil {
		return 0, err
	}

	return strconv.ParseFloat(s, bitSize)
}

func shouldBeString(input interface{}) (string, error) {
	switch v := input.(type) {
	case string:
		return v, nil
	case []byte:
		return string(v), nil
	default:
		return "", fmt.Errorf("cannot convert %v of type %T to string", input, input)
	}
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

func (s *Suite) shouldSkipFeatureOrRule(featureOrRuleTags []*msgs.Tag) bool {
	for _, tag := range featureOrRuleTags {
		if contains(s.options.ignoreTags, tag.Name) {
			return true
		}
	}

	return false
}

func (s *Suite) shouldSkipScenario(scenarioTags []*msgs.Tag) bool {
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

func (s *Suite) getBackgroundSteps(backgrounds []*msgs.Background) []*msgs.Step {
	result := []*msgs.Step{}
	for _, background := range backgrounds {
		result = append(result, background.Steps...)
	}

	return result
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
