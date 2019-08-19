package testhttp

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

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
	_ = addStep.AddStep(`^the response contains a valid JSON$`, testHTTP.validJSON)
	_ = addStep.AddStep(`^the response is "(.*)"$`, testHTTP.theResponseIs)

	return thhtp
}

func (t testHTTPMethods) theResponseIs(ctx context.Context) error {
	expectedResponse := ctx.GetStringParam(0)

	resp := ctx.Get(httpResponse{}).(Response)
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	ctx.Set(httpResponse{}, resp)
	if err != nil {
		return fmt.Errorf("an error while reading the body: %s", err)
	}

	if expectedResponse != string(body) {
		return errors.New("the body doesn't contain a valid JSON")
	}
	return nil
}

func (t testHTTPMethods) validJSON(ctx context.Context) error {
	resp := ctx.Get(httpResponse{}).(Response)
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	ctx.Set(httpResponse{}, resp)
	if err != nil {
		return fmt.Errorf("an error while reading the body: %s", err)
	}

	if !json.Valid(body) {
		return errors.New("the body doesn't contain a valid JSON")
	}
	return nil
}

func (t testHTTPMethods) makeRequest(ctx context.Context) error {
	method := ctx.GetStringParam(0)
	url := ctx.GetStringParam(1)
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return err
	}

	resp, err := t.tHTTP.MakeRequest(req)
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
