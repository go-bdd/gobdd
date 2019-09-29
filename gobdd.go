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

	"github.com/cucumber/gherkin-go"
	"github.com/go-bdd/gobdd/context"
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
}

// NewSuiteOptions creates a new suite configuration with default values
func NewSuiteOptions() SuiteOptions {
	return SuiteOptions{
		featuresPaths:  "features/*.feature",
		ignoreTags:     []string{},
		tags:           []string{},
		beforeScenario: []func(){},
		afterScenario:  []func(){},
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

func (s *Suite) runFeature(feature *gherkin.Feature) error {
	log.SetOutput(ioutil.Discard)
	hasErrors := false

	s.t.Run(feature.Name, func(t *testing.T) {
		var bkgSteps *gherkin.Background

		for _, child := range feature.Children {
			if scenario, ok := child.(*gherkin.Scenario); ok {
				if s.skipScenario(scenario.Tags) {
					t.Log(fmt.Sprintf("Skipping scenario %s", scenario.Name))
					continue
				}
				ctx := context.New()

				if bkgSteps != nil {
					s.runSteps(ctx, t, bkgSteps.Steps)
				}
				err := s.runScenario(ctx, scenario, bkgSteps, t)
				if err != nil {
					hasErrors = true
				}
			}

			if scenario, ok := child.(*gherkin.ScenarioOutline); ok {
				if s.skipScenario(scenario.Tags) {
					t.Log(fmt.Sprintf("Skipping scenario %s", scenario.Name))
					continue
				}

				ctx := context.New()
				if bkgSteps != nil {
					s.runSteps(ctx, t, bkgSteps.Steps)
				}

				err := s.runScenarioOutline(scenario, t)
				if err != nil {
					hasErrors = true
				}
			}

			if bkg, ok := child.(*gherkin.Background); ok {
				bkgSteps = bkg
			}
		}
	})

	if hasErrors {
		return errors.New("the feature contains errors")
	}

	return nil
}

func (s *Suite) runScenarioOutline(outline *gherkin.ScenarioOutline, t *testing.T) error {
	s.callBeforeScenarios()
	defer s.callAfterScenarios()
	t.Run(fmt.Sprintf("%s %s", outline.Keyword, outline.Name), func(t *testing.T) {
		for _, ex := range outline.Examples {
			if len(ex.TableBody) == 0 {
				continue
			}

			placeholders := ex.TableHeader.Cells
			groups := ex.TableBody

			for _, group := range groups {
				ctx := context.New()
				steps := s.runOutlineStep(outline, placeholders, group)
				s.runSteps(ctx, t, steps)
			}
		}
	})

	return nil
}

func (s *Suite) runOutlineStep(outline *gherkin.ScenarioOutline, placeholders []*gherkin.TableCell, group *gherkin.TableRow) []*gherkin.Step {
	var steps []*gherkin.Step
	for _, outlineStep := range outline.Steps {
		text := outlineStep.Text

		for i, placeholder := range placeholders {
			ph := "<" + placeholder.Value + ">"
			index := strings.Index(text, ph)
			originalText := text
			if index == -1 {
				continue
			}

			text = strings.Replace(text, "<"+placeholder.Value+">", group.Cells[i].Value, -1)
			t := getRegexpForVar(group.Cells[i].Value)

			def, err := s.findStepDef(originalText)
			if err != nil {
				continue
			}

			expr := strings.Replace(def.expr.String(), ph, t, -1)
			_ = s.AddStep(expr, def.f)
		}

		arg := s.getOutlineArguments(outlineStep, placeholders, group)

		// clone a step
		step := &gherkin.Step{
			Node:     outlineStep.Node,
			Text:     text,
			Keyword:  outlineStep.Keyword,
			Argument: arg,
		}
		steps = append(steps, step)
	}
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

func (s *Suite) getOutlineArguments(outlineStep *gherkin.Step, placeholders []*gherkin.TableCell, group *gherkin.TableRow) interface{} {
	arg := outlineStep.Argument

	switch t := outlineStep.Argument.(type) {
	case *gherkin.DataTable:
		arg = s.outlineDataTableArguments(t, placeholders, group)
	}
	return arg
}

func (s *Suite) outlineDataTableArguments(t *gherkin.DataTable, placeholders []*gherkin.TableCell, group *gherkin.TableRow) interface{} {
	tbl := &gherkin.DataTable{
		Node: t.Node,
		Rows: make([]*gherkin.TableRow, len(t.Rows)),
	}
	for i, row := range t.Rows {
		cells := make([]*gherkin.TableCell, len(row.Cells))
		for j, cell := range row.Cells {
			trans := cell.Value
			for i, placeholder := range placeholders {
				trans = strings.Replace(trans, "<"+placeholder.Value+">", group.Cells[i].Value, -1)
			}
			cells[j] = &gherkin.TableCell{
				Node:  cell.Node,
				Value: trans,
			}
		}
		tbl.Rows[i] = &gherkin.TableRow{
			Node:  row.Node,
			Cells: cells,
		}
	}
	return tbl
}

func (s *Suite) runScenario(ctx context.Context, scenario *gherkin.Scenario, bkg *gherkin.Background, t *testing.T) error {
	s.callBeforeScenarios()

	defer s.callAfterScenarios()
	t.Run(scenario.Name, func(t *testing.T) {
		if bkg != nil {
			steps := s.getBackgroundSteps(bkg)
			s.runSteps(ctx, t, steps)
		}
		s.runSteps(ctx, t, scenario.Steps)
	})
	return nil
}

func (s *Suite) runSteps(ctx context.Context, t *testing.T, steps []*gherkin.Step) {
	for _, step := range steps {
		s.runStep(ctx, t, step)
	}
}

func (s *Suite) runStep(ctx context.Context, t *testing.T, step *gherkin.Step) {
	defer func() {
		if r := recover(); r != nil {
			t.Error(r)
		}
	}()

	def, err := s.findStepDef(step.Text)
	if err != nil {
		t.Errorf("cannot find step definition for step: %s %s", step.Keyword, step.Text)
		return
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
		if len(params) + 1 != d.Type().NumIn() {
			t.Errorf("the step function %s accepts %d arguments but %d received", d.String(), d.Type().NumIn(), len(params) + 1)
			return
		}
		for i, v := range params {
			paramType := reflect.ValueOf(v)
			inType := d.Type().In(i + 1)

			if inType.Kind() == reflect.String {
				paramType = reflect.ValueOf(string(paramType.Interface().([]uint8)))
			}
			if inType.Kind() == reflect.Int {
				s := paramType.Interface().([]uint8)
				p, _ := strconv.Atoi(string(s))
				paramType = reflect.ValueOf(p)
			}
			in = append(in, paramType)
		}
		v := d.Call(in)[0]

		if !v.IsNil() {
			err = v.Interface().(error)
			t.Errorf(step.Keyword, step.Text, err)
		}
	})
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

func (s *Suite) skipScenario(scenarioTags []*gherkin.Tag) bool {
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

func (s *Suite) getBackgroundSteps(bkg *gherkin.Background) []*gherkin.Step {
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

	return "string"
}
