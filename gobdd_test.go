package gobdd

import (
	"errors"
	"testing"

	"github.com/go-bdd/gobdd/context"
)

func add(ctx context.Context) error {
	res := ctx.GetIntParam(0) + ctx.GetIntParam(1)
	ctx.Set("sumRes", res)
	return nil
}

func check(ctx context.Context) error {
	expected := ctx.GetIntParam(0)
	received := ctx.GetInt("sumRes")

	if expected != received {
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
	err := suite.AddStep(`I add <digit1> and <digit2>`, func(ctx context.Context) error {
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	err = suite.AddStep(`the result should equal <result>`, func(ctx context.Context) error {

		return nil
	})
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
