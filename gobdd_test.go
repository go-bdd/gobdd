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

func TestDifferentFuncTypes(t *testing.T) {
	suite := NewSuite(t, NewSuiteOptions().WithFeaturesPath("features/func_types.feature"))
	err := suite.AddStep(`I add ([+-]?[0-9]*[.]?[0-9]+) and ([+-]?[0-9]*[.]?[0-9]+)`, addf)
	if err != nil {
		t.Fatal(err)
	}

	err = suite.AddStep(`the result should equal ([+-]?[0-9]*[.]?[0-9]+)`, checkf)
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

func TestInvalidFunctionSignature(t *testing.T) {
	testCases := map[string]struct {
		f interface{}
	}{
		"nil":                              {},
		"func without return value":        {f: func(ctx context.Context) {}},
		"func with invalid return value":   {f: func(ctx context.Context) int { return 0 }},
		"func without arguments":           {f: func() error { return errors.New("") }},
		"func with invalid first argument": {f: func(i int) error { return errors.New("") }},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			suite := NewSuite(t, NewSuiteOptions())
			err := suite.AddStep("", testCase.f)
			if err == nil {
				t.Errorf("the function has invalid signature; the test should fail")
			}
		})
	}
}

func addf(ctx context.Context, var1, var2 float32) error {
	res := var1 + var2
	ctx.Set("sumRes", res)
	return nil
}

func add(ctx context.Context, var1, var2 int) error {
	res := var1 + var2
	ctx.Set("sumRes", res)
	return nil
}

func checkf(ctx context.Context, sum float32) error {
	received := ctx.GetFloat32("sumRes")

	if sum != received {
		return errors.New("the math does not work for you")
	}

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
