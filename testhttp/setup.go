package testhttp

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-bdd/assert"
	"github.com/go-bdd/gobdd/context"
	"github.com/go-bdd/gobdd/step"
	"io/ioutil"
	"net/http"
)

type addStepper interface {
	AddStep(step interface{}, f step.Func) error
}

type testHTTPMethods struct {
	tHTTP TestHTTP
}

type ResponseKey struct{}
type RequestKey struct{}

func Build(addStep addStepper, h httpHandler) TestHTTP {
	thhtp := TestHTTP{
		handler: h,
	}

	testHTTP := testHTTPMethods{tHTTP: thhtp}

	_ = addStep.AddStep(`^I make a (GET|POST|PUT|DELETE|OPTIONS) request to "([^"]*)"$`, testHTTP.makeRequest)
	_ = addStep.AddStep(`^the response code equals (\d+)$`, testHTTP.statusCodeEquals)
	_ = addStep.AddStep(`^the response contains a valid JSON$`, testHTTP.validJSON)
	_ = addStep.AddStep(`^the response is "(.*)"$`, testHTTP.theResponseIs)
	_ = addStep.AddStep(`^the response header "(.*)" equals "(.*)"$`, testHTTP.responseHeaderEquals)
	_ = addStep.AddStep(`^I have a (GET|POST|PUT|DELETE|OPTIONS) request "(.*)"$`, testHTTP.iHaveARequest)
	_ = addStep.AddStep(`^I set request header "Xyz" to "ZZZ"$`, testHTTP.iSetRequestSetTo)
	_ = addStep.AddStep(`^the request has body "(.*)"$`, testHTTP.theRequestHasBody)
	_ = addStep.AddStep(`^I make the request$`, testHTTP.iMakeRequest)

	return thhtp
}

func (t testHTTPMethods) iSetRequestSetTo(ctx context.Context) error {
	req := ctx.Get(RequestKey{}).(*http.Request)
	headerName := ctx.GetStringParam(0)
	value := ctx.GetStringParam(1)
	req.Header.Add(headerName, value)

	ctx.Set(RequestKey{}, req)
	return nil
}

func (t testHTTPMethods) responseHeaderEquals(ctx context.Context) error {
	resp := ctx.Get(ResponseKey{}).(Response)
	headerName := ctx.GetStringParam(0)
	expected := ctx.GetStringParam(1)
	given := resp.Header.Get(headerName)

	return assert.Equals(expected, given)
}

func (t testHTTPMethods) theRequestHasBody(ctx context.Context) error {
	req := ctx.Get(RequestKey{}).(*http.Request)
	body := ctx.GetStringParam(0)
	req.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(body)))
	ctx.Set(RequestKey{}, req)
	return nil
}

func (t testHTTPMethods) iHaveARequest(ctx context.Context) error {
	method := ctx.GetStringParam(0)
	url := ctx.GetStringParam(1)
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return err
	}

	ctx.Set(RequestKey{}, req)
	return nil
}

func (t testHTTPMethods) theResponseIs(ctx context.Context) error {
	expectedResponse := ctx.GetStringParam(0)

	resp := ctx.Get(ResponseKey{}).(Response)
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	ctx.Set(ResponseKey{}, resp)
	if err != nil {
		return fmt.Errorf("an error while reading the body: %s", err)
	}

	if expectedResponse != string(body) {
		return errors.New("the body doesn't contain a valid JSON")
	}
	return nil
}

func (t testHTTPMethods) validJSON(ctx context.Context) error {
	resp := ctx.Get(ResponseKey{}).(Response)
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	ctx.Set(ResponseKey{}, resp)
	if err != nil {
		return fmt.Errorf("an error while reading the body: %s", err)
	}

	if !json.Valid(body) {
		return errors.New("the body doesn't contain a valid JSON")
	}
	return nil
}

func (t testHTTPMethods) iMakeRequest(ctx context.Context) error {
	req := ctx.Get(RequestKey{}).(*http.Request)
	resp, err := t.tHTTP.MakeRequest(req)
	if err != nil {
		return err
	}

	ctx.Set(ResponseKey{}, resp)
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

	ctx.Set(ResponseKey{}, resp)
	return nil
}

func (t testHTTPMethods) statusCodeEquals(ctx context.Context) error {
	expectedStatus := ctx.GetIntParam(0)
	resp := ctx.Get(ResponseKey{}).(Response)

	if expectedStatus != resp.Code {
		return fmt.Errorf("expected status code: %d but %d given", expectedStatus, resp.Code)
	}
	return nil
}
