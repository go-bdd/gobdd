package testhttp

import (
	"fmt"

	"github.com/go-bdd/gobdd/context"
	"github.com/go-bdd/gobdd/step"
)

type addStepper interface {
	AddStep(step interface{}, f step.Func) error
}

type testHTTPMethods struct {
	tHTTP TestHTTP
}

type httpResponse struct{}

func Build(addStep addStepper, h httpHandler) TestHTTP {
	thhtp := TestHTTP{
		handler: h,
	}

	testHTTP := testHTTPMethods{tHTTP: thhtp}

	_ = addStep.AddStep(`^I make a (GET|POST|PUT|DELETE|OPTIONS) request to "([^"]*)"$`, testHTTP.makeRequest)
	_ = addStep.AddStep(`^the response code equals (\d+)$`, testHTTP.statusCodeEquals)

	return thhtp
}

func (t testHTTPMethods) makeRequest(ctx context.Context) error {
	method := ctx.GetStringParam(0)
	url := ctx.GetStringParam(1)
	resp, err := t.tHTTP.Request(method, url, nil)
	if err != nil {
		return err
	}

	ctx.Set(httpResponse{}, resp)
	return nil
}

func (t testHTTPMethods) statusCodeEquals(ctx context.Context) error {
	expectedStatus := ctx.GetIntParam(0)
	resp := ctx.Get(httpResponse{}).(Response)

	if expectedStatus != resp.Code {
		return fmt.Errorf("expected status code: %d but %d given", expectedStatus, resp.Code)
	}
	return nil
}
