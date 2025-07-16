package gobdd

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"testing"

	msgs "github.com/cucumber/messages/go/v28"
	"github.com/go-bdd/assert"
	"github.com/stretchr/testify/require"
)

func TestScenarios(t *testing.T) {
	suite := NewSuite(t, WithFeaturesPath("features/example.feature"))
	compiled := regexp.MustCompile(`I add (\d+) and (\d+)`)
	suite.AddRegexStep(compiled, add)
	compiled = regexp.MustCompile(`the result should equal (\d+)`)
	suite.AddRegexStep(compiled, check)
	suite.Run()
}

func TestRule(t *testing.T) {
	suite := NewSuite(t, WithFeaturesPath("features/example_rule.feature"))
	compiled := regexp.MustCompile(`I add (\d+) and (\d+)`)
	suite.AddRegexStep(compiled, add)
	compiled = regexp.MustCompile(`the result should equal (\d+)`)
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

func TestParameterTypes(t *testing.T) {
	suite := NewSuite(t, WithFeaturesPath("features/parameter-types.feature"))
	suite.AddStep(`I add {int} and {int}`, add)
	suite.AddStep(`the result should equal {int}`, check)
	suite.AddStep(`I add floats {float} and {float}`, addf)
	suite.AddStep(`the result should equal float {float}`, checkf)
	suite.AddStep(`the result should equal text {text}`, checkt)
	suite.AddStep(`I use word {word}`, func(t StepTest, _ Context, word string) {
		if word != "pizza" {
			t.Fatal("it should be pizza")
		}
	})
	suite.AddStep(`I use text {text}`, func(_ StepTest, ctx Context, text string) {
		ctx.Set("stringRes", text)
	})
	suite.AddStep(`I concat word {word} and text {text}`, concat)
	suite.AddStep(`I format text {text} with int {int}`, func(_ StepTest, ctx Context, format string, value int) {
		ctx.Set("stringRes", fmt.Sprintf(format, value))
	})

	suite.Run()
}

func TestArguments(t *testing.T) {
	suite := NewSuite(t, WithFeaturesPath("features/argument.feature"))
	suite.AddStep(`the result should equal argument:`, checkt)
	suite.AddStep(`I concat text {text} and argument:`, concat)

	suite.Run()
}

func TestDatatable(t *testing.T) {
	suite := NewSuite(t, WithFeaturesPath("features/datatable.feature"))
	suite.AddStep(`I concat all the columns and row together using {text} to separate the columns`, concatTable)
	suite.AddStep(`the result should equal argument:`, checkt)

	suite.Run()
}

func TestScenarioOutlineExecutesAllTests(t *testing.T) {
	c := 0
	suite := NewSuite(t, WithFeaturesPath("features/outline.feature"))
	suite.AddStep(`I add (\d+) and (\d+)`, add)
	suite.AddStep(`the result should equal (\d+)`, func(t StepTest, ctx Context, sum int) {
		c++
		check(t, ctx, sum)
	})

	suite.Run()

	if err := assert.Equals(2, c); err != nil {
		t.Errorf("expected to run %d times but %d got", 2, c)
	}
}

func TestStepFromExample(t *testing.T) {
	s := NewSuite(t)
	st, expr := s.stepFromExample("I add <d1> and <d2>", &msgs.TableRow{
		Cells: []*msgs.TableCell{
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
	suite.AddStep(`I concat word {word} and text {text}`, concat)
	suite.AddStep(`the result should equal text {text}`, checkt)
	suite.AddStep(`the result should equal (\d+)`, check)

	suite.Run()
}

func TestTags(t *testing.T) {
	suite := NewSuite(t, WithFeaturesPath("features/tags.feature"), WithTags("@tag"))
	suite.AddStep(`fail the test`, fail)
	suite.AddStep(`the test should pass`, pass)

	suite.Run()
}

func TestFilterFeatureWithTags(t *testing.T) {
	suite := NewSuite(t, WithFeaturesPath("features/filter_tags_*.feature"), WithTags("@run-this"))
	c := false

	suite.AddStep(`the test should pass`, func(_ StepTest, _ Context) {
		c = true
	})
	suite.AddStep(`fail the test`, fail)

	suite.Run()

	if err := assert.Equals(true, c); err != nil {
		t.Error(err)
	}
}

func TestWithAfterScenario(t *testing.T) {
	c := false
	suite := NewSuite(t, WithFeaturesPath("features/empty.feature"), WithAfterScenario(func(_ Context) {
		c = true
	}))
	suite.Run()

	if err := assert.Equals(true, c); err != nil {
		t.Error(err)
	}
}

func TestWithBeforeScenario(t *testing.T) {
	c := false
	suite := NewSuite(t, WithFeaturesPath("features/empty.feature"), WithBeforeScenario(func(_ Context) {
		c = true
	}))
	suite.Run()

	if err := assert.Equals(true, c); err != nil {
		t.Error(err)
	}
}

func TestWithAfterStep(t *testing.T) {
	c := 0
	suite := NewSuite(t, WithFeaturesPath("features/background.feature"), WithAfterStep(func(ctx Context) {
		c++

		// feature should be *msgs.Feature
		feature, err := ctx.Get(FeatureKey{})
		require.NoError(t, err)
		if _, ok := feature.(*msgs.Feature); !ok {
			t.Errorf("expected feature but got %T", feature)
		}

		// scenario should be *msgs.Scenario
		scenario, err := ctx.Get(ScenarioKey{})
		require.NoError(t, err)
		if _, ok := scenario.(*msgs.Scenario); !ok {
			t.Errorf("expected scenario but got %T", scenario)
		}
	}))
	suite.AddStep(`I add (\d+) and (\d+)`, add)
	suite.AddStep(`the result should equal (\d+)`, check)
	suite.AddStep(`I concat word {word} and text {text}`, concat)
	suite.AddStep(`the result should equal text {text}`, checkt)

	suite.Run()

	if err := assert.Equals(6, c); err != nil {
		t.Error(err)
	}
}

func TestWithBeforeStep(t *testing.T) {
	c := 0
	suite := NewSuite(t, WithFeaturesPath("features/background.feature"), WithBeforeStep(func(_ Context) {
		c++
	}))
	suite.AddStep(`I add (\d+) and (\d+)`, add)
	suite.AddStep(`the result should equal (\d+)`, check)
	suite.AddStep(`I concat word {word} and text {text}`, concat)
	suite.AddStep(`the result should equal text {text}`, checkt)

	suite.Run()

	if err := assert.Equals(6, c); err != nil {
		t.Error(err)
	}
}

func TestIgnoredTags(t *testing.T) {
	suite := NewSuite(t, WithFeaturesPath("features/ignored_*tags.feature"), WithIgnoredTags("@ignore"))
	suite.AddStep(`the test should pass`, pass)
	suite.AddStep(`fail the test`, fail)
	suite.Run()
}

func TestInvalidFunctionSignature(t *testing.T) {
	testCases := map[string]struct {
		f interface{}
	}{
		"nil":                              {},
		"func without return value":        {f: func(_ Context) {}},
		"func with invalid return value":   {f: func(_ Context) int { return 0 }},
		"func without arguments":           {f: func() error { return errors.New("") }},
		"func with invalid first argument": {f: func(_ int) error { return errors.New("") }},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			tester := &mockTester{}
			suite := NewSuite(tester)
			suite.AddStep("", testCase.f)
			suite.Run()
			if err := assert.Equals(1, tester.fatalCalled); err != nil {
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
			def.run(NewContext(), tester, nil)
			err := assert.Equals(testCase.expectedErrors, tester.errors)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func addf(_ StepTest, ctx Context, var1, var2 float32) {
	res := var1 + var2
	ctx.Set("sumRes", res)
}

func add(_ StepTest, ctx Context, var1, var2 int) {
	res := var1 + var2
	ctx.Set("sumRes", res)
}

func concat(_ StepTest, ctx Context, var1, var2 string) {
	ctx.Set("stringRes", var1+var2)
}

func concatTable(_ StepTest, ctx Context, separator string, table msgs.DataTable) {
	rows := make([]string, 0, len(table.Rows))
	for _, row := range table.Rows {
		values := make([]string, 0, len(row.Cells))
		for _, cell := range row.Cells {
			values = append(values, cell.Value)
		}

		rows = append(rows, strings.Join(values, separator))
	}

	ctx.Set("stringRes", strings.Join(rows, "\n"))
}

func checkt(t StepTest, ctx Context, text string) {
	received, err := ctx.GetString("stringRes")
	if err != nil {
		t.Error(err.Error())

		return
	}

	if text != received {
		t.Errorf("expected %s but %s received", text, received)
	}
}

func checkf(t StepTest, ctx Context, sum float32) {
	received, err := ctx.Get("sumRes")
	if err != nil {
		t.Error(err.Error())

		return
	}

	if sum != received {
		t.Errorf("expected %f but %f received", sum, received)
	}
}

func check(t StepTest, ctx Context, sum int) {
	received, err := ctx.Get("sumRes")
	if err != nil {
		t.Error(err)
		return
	}

	if sum != received {
		t.Errorf("expected %d but %d received", sum, received)
	}
}

func fail(t StepTest, _ Context) {
	t.Error("the step should never be executed")
}

func failure(t StepTest, _ Context) {
	t.Error("the step failed")
}

func panics(_ StepTest, _ Context) {
	panic(errors.New("the step panicked"))
}

func pass(_ StepTest, _ Context) {}

type mockTester struct {
	testing.T
	fatalCalled int
	errors      []string
}

var _ TestingT = (*mockTester)(nil)

func (m *mockTester) Log(...interface{}) {}

func (m *mockTester) Logf(string, ...interface{}) {}

func (m *mockTester) Fatal(...interface{}) {
	m.fatalCalled++
}

func (m *mockTester) Fatalf(string, ...interface{}) {}

func (m *mockTester) Error(a ...interface{}) {
	m.errors = append(m.errors, fmt.Sprintf("%s", a...))
}

func (m *mockTester) Errorf(format string, a ...interface{}) {
	m.errors = append(m.errors, fmt.Sprintf(format, a...))
}

func (m *mockTester) Parallel() {}

func (m *mockTester) Fail() {}

func (m *mockTester) FailNow() {}
