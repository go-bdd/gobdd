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
	t *testing.T
	steps []StepDef
}

type StepFunc func(ctx Context) error

type StepDef struct {
	expr *regexp.Regexp
	f StepFunc
}

func NewSuite(t *testing.T) *Suite {
	return &Suite{
		t: t,
	}
}

func (suite *Suite) AddStep(step interface{}, f StepFunc)  {
	var regex *regexp.Regexp

	switch t := step.(type) {
	case *regexp.Regexp:
		regex = t
	case string:
		regex = regexp.MustCompile(t)
	case []byte:
		regex = regexp.MustCompile(string(t))
	default:
		panic(fmt.Sprintf("expecting expr to be a *regexp.Regexp or a string, got type: %T", step))
	}

	suite.steps = append(suite.steps, StepDef{
		expr:regex,
		f: f,
	})
}

func (suite *Suite) Run() {
	files, err := filepath.Glob("features/*.feature")
	if err != nil {
		suite.t.Fatalf("cannot find features/ directory")
	}

	for _, file := range files  {
		f, err := os.Open(file)
		if err != nil {
			suite.t.Fatalf("cannot find features/ directory")
		}
		fileIO := bufio.NewReader(f)
		doc, err := gherkin.ParseGherkinDocument(fileIO)
		if err != nil {
			suite.t.Fatalf("error while loading document: %s", err)
		}

		_ = f.Close()

		if doc.Feature == nil {
			continue
		}

		suite.runFeature(doc.Feature)
	}
}

func (suite *Suite) runFeature(feature *gherkin.Feature) {
	fmt.Printf("Feature: %s\n", feature.Name)

	for _, s := range feature.Children {
		scenario := s.(*gherkin.Scenario)
		fmt.Printf("  Scenario: %s\n", scenario.Name)
		ctx := newContext()

		for _, step := range scenario.Steps  {
			fmt.Printf("    %s: %s\n", step.Keyword, step.Text)
			def, err := suite.findStepDef(step.Text)
			if err != nil {
				suite.t.Errorf("cannot find step definition for step '%s'", step.Text)
				continue
			}

			b := def.expr.FindSubmatch([]byte(step.Text))
			ctx.setParams(b[1:])

			err = def.f(ctx)
			if err != nil {
				suite.t.Error(err)
			}
		}
	}
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
