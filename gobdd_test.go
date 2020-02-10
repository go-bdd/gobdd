package gobdd

import (
	"errors"
	"fmt"
	"testing"

	"github.com/cucumber/messages-go/v9"
	"github.com/go-bdd/assert"
	"github.com/go-bdd/gobdd/context"
)

func TestScenarios(t *testing.T) {
	suite := NewSuite(t, NewSuiteOptions().WithFeaturesPath("features/example.feature"))
	suite.AddStep(`I add (\d+) and (\d+)`, add)
	suite.AddStep(`the result should equal (\d+)`, check)

	suite.Run()
}

func TestDifferentFuncTypes(t *testing.T) {
	suite := NewSuite(t, NewSuiteOptions().WithFeaturesPath("features/func_types.feature"))
	suite.AddStep(`I add ([+-]?[0-9]*[.]?[0-9]+) and ([+-]?[0-9]*[.]?[0-9]+)`, addf)
	suite.AddStep(`the result should equal ([+-]?[0-9]*[.]?[0-9]+)`, checkf)

	suite.Run()
}

func TestScenarioOutline(t *testing.T) {
	suite := NewSuite(t, NewSuiteOptions().WithFeaturesPath("features/outline.feature"))
	suite.AddStep(`I add (\d+) and (\d+)`, add)
	suite.AddStep(`the result should equal <result>`, check)

	suite.Run()
}

func TestBackground(t *testing.T) {
	options := NewSuiteOptions().WithFeaturesPath("features/background.feature")
	suite := NewSuite(t, options)
	suite.AddStep(`I add (\d+) and (\d+)`, add)
	suite.AddStep(`the result should equal (\d+)`, check)

	suite.Run()
}

func TestTags(t *testing.T) {
	options := NewSuiteOptions().WithFeaturesPath("features/tags.feature").WithTags([]string{"@tag"})
	suite := NewSuite(t, options)
	suite.AddStep(`fail the test`, fail)
	suite.AddStep(`the test should pass`, pass)

	suite.Run()
}

func TestIgnoredTags(t *testing.T) {
	options := NewSuiteOptions().WithFeaturesPath("features/ignored_tags.feature")
	options = options.WithIgnoredTags([]string{"@ignore"})
	suite := NewSuite(t, options)
	suite.AddStep(`fail the test`, fail)
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
			tester := &mockTester{}
			suite := NewSuite(tester, NewSuiteOptions())
			suite.AddStep("", testCase.f)
			suite.Run()
			err := assert.Equals(1, tester.fatalCalled)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestFailureOutput(t *testing.T) {
	testCases := []struct {
		name           string
		f              interface{}
		expectedErrors []string
	}{
		{name: "passes", f: pass, expectedErrors: nil},
		{name: "returns error", f: failure, expectedErrors: []string{"Test text: the step failed"}},
		{name: "step panics", f: panics, expectedErrors: []string{"the step panicked"}},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			def := stepDef{f: testCase.f}
			step := &messages.GherkinDocument_Feature_Step{Keyword: "Test", Text: "text"}

			tester := &mockTester{}
			def.run(context.New(), tester, step, nil)
			err := assert.Equals(testCase.expectedErrors, tester.errors)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func addf(ctx context.Context, var1, var2 float32) (context.Context, error) {
	res := var1 + var2
	ctx.Set("sumRes", res)
	return ctx, nil
}

func add(ctx context.Context, var1, var2 int) (context.Context, error) {
	res := var1 + var2
	ctx.Set("sumRes", res)
	return ctx, nil
}

func checkf(ctx context.Context, sum float32) (context.Context, error) {
	received := ctx.Get("sumRes")

	if sum != received {
		return ctx, errors.New("the math does not work for you")
	}

	return ctx, nil
}

func check(ctx context.Context, sum int) (context.Context, error) {
	received := ctx.Get("sumRes")

	if sum != received {
		return ctx, errors.New("the math does not work for you")
	}

	return ctx, nil
}

func fail(ctx context.Context) (context.Context, error) {
	return ctx, errors.New("the step should never be executed")
}

func failure(ctx context.Context) (context.Context, error) {
	return ctx, errors.New("the step failed")
}

func panics(_ context.Context) (context.Context, error) {
	panic(errors.New("the step panicked"))
}

func pass(ctx context.Context) (context.Context, error) {
	return ctx, nil
}

type mockTester struct {
	fatalCalled int
	errors      []string
}

func (m *mockTester) Log(...interface{}) {
}

func (m *mockTester) Fatal(...interface{}) {
	m.fatalCalled++
}

func (m *mockTester) Fatalf(string, ...interface{}) {
}

func (m *mockTester) Errorf(format string, a ...interface{}) {
	m.errors = append(m.errors, fmt.Sprintf(format, a...))
}

func (m *mockTester) Parallel() {
}

func (m *mockTester) Fail() {
}

func (m *mockTester) Run(_ string, _ func(t *testing.T)) bool {
	return true
}
