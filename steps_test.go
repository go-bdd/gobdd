package gobdd

import (
	"testing"

	"github.com/go-bdd/gobdd/context"
)

func TestValidateStepFunc(t *testing.T) {
	testCases := map[string]interface{}{
		"function without arguments":           func() {},
		"function with 1 argument":             func(StepTest) {},
		"function with invalid first argument": func(int, context.Context) {},
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
	if err := validateStepFunc(func(StepTest, context.Context) {}); err != nil {
		t.Errorf("the test should NOT fail for the function: %s", err)
	}
}

func TestValidateStepFunc_ReturnContext(t *testing.T) {
	if err := validateStepFunc(func(StepTest, context.Context) context.Context { return context.Context{} }); err != nil {
		t.Errorf("step function returning a context should NOT fail validation: %s", err)
	}
}
