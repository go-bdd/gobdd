package gobdd

import (
	"errors"
	"testing"

	"github.com/go-bdd/gobdd/context"
)

func TestScenarios(t *testing.T) {
	suite := NewSuite(t, NewSuiteOptions().WithFeaturesPath("features/example.feature"))
	err := suite.AddStep(`I add (\d+) and (\d+)`, add)
	if err != nil {
		t.Fatal(err)
	}

	err = suite.AddStep(`the result should equal (\d+)`, check)
	if err != nil {
		t.Fatal(err)
	}

	suite.Run()
}

func TestScenarioOutline(t *testing.T) {
	suite := NewSuite(t, NewSuiteOptions().WithFeaturesPath("features/outline.feature"))
	err := suite.AddStep(`I add <digit1> and <digit2>`, add)
	if err != nil {
		t.Fatal(err)
	}

	err = suite.AddStep(`the result should equal <result>`, check)
	if err != nil {
		t.Fatal(err)
	}

	suite.Run()
}

func TestBackground(t *testing.T) {
	options := NewSuiteOptions().WithFeaturesPath("features/background.feature")
	suite := NewSuite(t, options)
	err := suite.AddStep(`I add (\d+) and (\d+)`, add)
	if err != nil {
		t.Fatal(err)
	}

	err = suite.AddStep(`the result should equal (\d+)`, check)
	if err != nil {
		t.Fatal(err)
	}

	suite.Run()
}

func TestTags(t *testing.T) {
	options := NewSuiteOptions().WithFeaturesPath("features/tags.feature").WithTags([]string{"@tag"})
	suite := NewSuite(t, options)
	err := suite.AddStep(`fail the test`, fail)
	if err != nil {
		t.Fatal(err)
	}

	err = suite.AddStep(`the test should pass`, pass)
	if err != nil {
		t.Fatal(err)
	}

	suite.Run()
}

func TestIgnoredTags(t *testing.T) {
	options := NewSuiteOptions().WithFeaturesPath("features/ignored_tags.feature")
	options = options.WithIgnoredTags([]string{"@ignore"})
	suite := NewSuite(t, options)
	err := suite.AddStep(`fail the test`, fail)
	if err != nil {
		t.Fatal(err)
	}

	suite.Run()
}


func add(ctx context.Context, var1, var2 int) error {
	res := var1 + var2
	ctx.Set("sumRes", res)
	return nil
}

func check(ctx context.Context, sum int) error {
	received := ctx.GetInt("sumRes")

	if sum != received {
		return errors.New("the math does not work for you")
	}

	return nil
}

func fail(context.Context) error {
	return errors.New("the step should never be executed")
}

func pass(context.Context) error {
	return nil
}