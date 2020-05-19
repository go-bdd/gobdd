package gobdd

import (
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/cucumber/messages-go/v12"
	"github.com/go-bdd/assert"
	"github.com/go-bdd/gobdd/context"
)

func TestScenarios(t *testing.T) {
	suite := NewSuite(t, WithFeaturesPath("features/example.feature"))
	compiled, _ := regexp.Compile(`I add (\d+) and (\d+)`)
	suite.AddRegexStep(compiled, add)
	compiled, _ = regexp.Compile(`the result should equal (\d+)`)
	suite.AddRegexStep(compiled, check)

	suite.Run()
}

func TestAddStepWithRegexp(t *testing.T) {
	suite := NewSuite(t, WithFeaturesPath("features/example.feature"))
	suite.AddStep(`I add (\d+) and (\d+)`, add)
	suite.AddStep(`the result should equal (\d+)`, check)

	suite.Run()
}

func TestDifferentFuncTypes(t *testing.T) {
	suite := NewSuite(t, WithFeaturesPath("features/func_types.feature"))
	suite.AddStep(`I add ([+-]?[0-9]*[.]?[0-9]+) and ([+-]?[0-9]*[.]?[0-9]+)`, addf)
	suite.AddStep(`the result should equal ([+-]?[0-9]*[.]?[0-9]+)`, checkf)

	suite.Run()
}

func TestScenarioOutline(t *testing.T) {
	suite := NewSuite(t, WithFeaturesPath("features/outline.feature"))
	suite.AddStep(`I add (\d+) and (\d+)`, add)
	suite.AddStep(`the result should equal (\d+)`, check)

	suite.Run()
}

func TestScenarioOutlineExecutesAllTests(t *testing.T) {
	c := 0
	suite := NewSuite(t, WithFeaturesPath("features/outline.feature"))
	suite.AddStep(`I add (\d+) and (\d+)`, add)
	suite.AddStep(`the result should equal (\d+)`, func(t StepTest, ctx context.Context, sum int) context.Context {
		c++
		return check(t, ctx, sum)
	})

	suite.Run()
	if err := assert.Equals(2, c); err != nil {
		t.Errorf("expected to run %d times but %d got", 2, c)
	}
}

func TestStepFromExample(t *testing.T) {
	s := NewSuite(t)
	st, expr := s.stepFromExample("I add <d1> and <d2>", &messages.GherkinDocument_Feature_TableRow{
		Cells: []*messages.GherkinDocument_Feature_TableRow_TableCell{
			{Value: "1"},
			{Value: "2"},
		},
	}, []string{"<d1>", "<d2>"})
	if err := assert.NotNil(st); err != nil {
		t.Error(err)
	}

	if err := assert.Equals("I add 1 and 2", st); err != nil {
		t.Error(err)
	}

	if err := assert.Equals(`I add (\d+) and (\d+)`, expr); err != nil {
		t.Error(err)
	}
}

func TestBackground(t *testing.T) {
	suite := NewSuite(t, WithFeaturesPath("features/background.feature"))
	suite.AddStep(`I add (\d+) and (\d+)`, add)
	suite.AddStep(`the result should equal (\d+)`, check)

	suite.Run()
}

func TestTags(t *testing.T) {
	suite := NewSuite(t, WithFeaturesPath("features/tags.feature"), WithTags([]string{"@tag"}))
	suite.AddStep(`fail the test`, fail)
	suite.AddStep(`the test should pass`, pass)

	suite.Run()
}

func TestIgnoredTags(t *testing.T) {
	suite := NewSuite(t, WithFeaturesPath("features/ignored_tags.feature"), WithIgnoredTags([]string{"@ignore"}))
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
			suite := NewSuite(tester)
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
		{name: "returns error", f: failure, expectedErrors: []string{"the step failed"}},
		{name: "step panics", f: panics, expectedErrors: []string{"the step panicked"}},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			def := stepDef{f: testCase.f}

			tester := &mockTester{}
			def.run(context.New(), tester, nil)
			err := assert.Equals(testCase.expectedErrors, tester.errors)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func addf(t StepTest, ctx context.Context, var1, var2 float32) context.Context {
	res := var1 + var2
	ctx.Set("sumRes", res)
	return ctx
}

func add(t StepTest, ctx context.Context, var1, var2 int) context.Context {
	res := var1 + var2
	ctx.Set("sumRes", res)
	return ctx
}

func checkf(t StepTest, ctx context.Context, sum float32) context.Context {
	received, err := ctx.Get("sumRes")
	if err != nil {
		t.Error(err.Error())
		return ctx
	}

	if sum != received {
		t.Error("the sum doesn't match")
	}

	return ctx
}

func check(t StepTest, ctx context.Context, sum int) context.Context {
	received, err := ctx.Get("sumRes")
	if err != nil {
		t.Error(err)
		return ctx
	}

	if sum != received {
		t.Errorf("expected %d but %d received", sum, received)
		return ctx
	}

	return ctx
}

func fail(t StepTest, ctx context.Context) context.Context {
	t.Error("the step should never be executed")
	return ctx
}

func failure(t StepTest, ctx context.Context) context.Context {
	t.Error("the step failed")
	return ctx
}

func panics(t StepTest, _ context.Context) context.Context {
	panic(errors.New("the step panicked"))
}

func pass(t StepTest, ctx context.Context) context.Context {
	return ctx
}

type mockTester struct {
	fatalCalled int
	errors      []string
}

func (m *mockTester) Log(...interface{}) {
}

func (m *mockTester) Logf(string, ...interface{}) {
}
func (m *mockTester) Fatal(...interface{}) {
	m.fatalCalled++
}

func (m *mockTester) Fatalf(string, ...interface{}) {
}

func (m *mockTester) Error(a ...interface{}) {
	m.errors = append(m.errors, fmt.Sprintf("%s", a...))
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
