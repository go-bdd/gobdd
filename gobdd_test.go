package gobdd

import (
	"errors"
	"testing"
)

func add(ctx Context) error {
	res := ctx.GetIntParam(0) + ctx.GetIntParam(1)
	ctx.Set("sumRes", res)
	return nil
}

func check(ctx Context) error {
	expected := ctx.GetIntParam(0)
	received := ctx.GetInt("sumRes")

	if expected != received {
		return errors.New("the math does not work for you")
	}

	return nil
}

func TestScenarios(t *testing.T) {
	suite := NewSuite(t)
	suite.AddStep(`I add (\d+) and (\d+)`, add)
	suite.AddStep(`I the result should equal (\d+)`, check)
	suite.Run()
}
