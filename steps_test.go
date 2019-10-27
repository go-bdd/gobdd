package gobdd

import (
	"context"
	"testing"
)

func TestValidateStepFunc(t *testing.T) {
	testCases:= map[string]interface{}{
		"function without arguments": func()(context.Context, error){ return nil, nil},
		"function with invalid first arguments": func(int)(context.Context, error){ return nil, nil},
		"function without returned values": func(context.Context){},
		"function with invalid first returned value": func(context.Context)(int, error){ return 0, nil},
		"function with invalid second returned value": func(context.Context)(context.Context, int){ return nil, 0},
	}

	for name, testCase := range testCases{
		t.Run(name, func(t *testing.T) {
			if err := validateStepFunc(testCase); err == nil {
				t.Errorf("the test should fail for the function")
			}
		})
	}
}

func TestValidateStepFunc_ValidFunction(t *testing.T) {
	if err := validateStepFunc(func(context.Context) (context.Context, error) { return nil, nil}); err != nil {
		t.Errorf("the test should NOT fail for the function: %s",err)
	}
}