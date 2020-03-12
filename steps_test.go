package gobdd

import (
	"testing"

	"github.com/go-bdd/gobdd/context"
)

func TestValidateStepFunc(t *testing.T) {
	testCases := map[string]interface{}{
		"function without arguments":                 func() context.Context { return context.Context{} },
		"function with 1 argument":                   func(StepTest) context.Context { return context.Context{} },
		"function with invalid first argument":       func(int, context.Context) context.Context { return context.Context{} },
		"function without returned values":           func(StepTest, context.Context) {},
		"function with invalid first returned value": func(StepTest, context.Context) int { return 0 },
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			if err := validateStepFunc(testCase); err == nil {
				t.Errorf("the test should fail for the function")
			}
		})
	}
}

func TestValidateStepFunc_ValidFunction(t *testing.T) {
	if err := validateStepFunc(func(StepTest, context.Context) context.Context { return context.Context{} }); err != nil {
		t.Errorf("the test should NOT fail for the function: %s", err)
	}
}
