package testhttp

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/go-bdd/assert"
	"github.com/go-bdd/gobdd/context"
)

type StepTest interface {
	Log(...interface{})
	Logf(string, ...interface{})
	Fatal(...interface{})
	Fatalf(string, ...interface{})
	Errorf(string, ...interface{})
	Error(...interface{})
}

type addStepper interface {
	AddStep(definition string, step interface{})
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

	addStep.AddStep(`^I make a (GET|POST|PUT|DELETE|OPTIONS) request to "([^"]*)"$`, testHTTP.makeRequest)
	addStep.AddStep(`^the response code equals (\d+)$`, testHTTP.statusCodeEquals)
	addStep.AddStep(`^the response contains a valid JSON$`, testHTTP.validJSON)
	addStep.AddStep(`^the response is "(.*)"$`, testHTTP.theResponseIs)
	addStep.AddStep(`^the response header "(.*)" equals "(.*)"$`, testHTTP.responseHeaderEquals)
	addStep.AddStep(`^I have a (GET|POST|PUT|DELETE|OPTIONS) request "(.*)"$`, testHTTP.iHaveARequest)
	addStep.AddStep(`^I set request header "(.*)" to "(.*)"$`, testHTTP.iSetRequestSetTo)
	addStep.AddStep(`^I set request body to "([^"]*)"$`, testHTTP.iSetRequestBodyTo)
	addStep.AddStep(`^the request has body "(.*)"$`, testHTTP.theRequestHasBody)
	addStep.AddStep(`^I make the request$`, testHTTP.iMakeRequest)

	return thhtp
}

func (th testHTTPMethods) iSetRequestBodyTo(t StepTest, ctx context.Context, body string) context.Context {
	r, err := GetRequest(ctx)
	if err != nil {
		t.Error(err)
		return ctx
	}
	r.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(body)))
	ctx.Set(RequestKey{}, r)
	return ctx
}

func (th testHTTPMethods) iSetRequestSetTo(t StepTest, ctx context.Context, headerName, value string) context.Context {
	req, err := GetRequest(ctx)
	if err != nil {
		t.Error(err)
		return ctx
	}
	req.Header.Add(headerName, value)

	ctx.Set(RequestKey{}, req)
	return ctx
}

func (th testHTTPMethods) responseHeaderEquals(t StepTest, ctx context.Context, headerName, expected string) context.Context {
	resp, err := GetResponse(ctx)
	if err != nil {
		t.Error(err)
		return ctx
	}
	given := resp.Header.Get(headerName)

	if err := assert.Equals(expected, given); err != nil {
		t.Error(err)
	}
	return ctx
}

func (th testHTTPMethods) theRequestHasBody(t StepTest, ctx context.Context, body string) context.Context {
	req, err := GetRequest(ctx)
	if err != nil {
		t.Error(err)
		return ctx
	}
	req.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(body)))
	ctx.Set(RequestKey{}, req)
	return ctx
}

func (th testHTTPMethods) iHaveARequest(t StepTest, ctx context.Context, method, url string) context.Context {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		t.Error(err)
		return ctx
	}

	ctx.Set(RequestKey{}, req)
	return ctx
}

func (th testHTTPMethods) theResponseIs(t StepTest, ctx context.Context, expectedResponse string) context.Context {
	resp, err := GetResponse(ctx)
	if err != nil {
		t.Error(err)
		return ctx
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	ctx.Set(ResponseKey{}, resp)
	if err != nil {
		t.Error(fmt.Errorf("an error while reading the body: %s", err))
		return ctx
	}

	if expectedResponse != string(body) {
		t.Error(errors.New("the body doesn't contain a valid JSON"))
	}
	return ctx
}

func (th testHTTPMethods) validJSON(t StepTest, ctx context.Context) context.Context {
	resp, err := GetResponse(ctx)
	if err != nil {
		t.Error(err)
		return ctx
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	ctx.Set(ResponseKey{}, resp)
	if err != nil {
		t.Error(fmt.Errorf("an error while reading the body: %s", err))
		return ctx
	}

	if !json.Valid(body) {
		t.Error(errors.New("the body doesn't contain a valid JSON"))
	}
	return ctx
}

func (th testHTTPMethods) iMakeRequest(t StepTest, ctx context.Context) context.Context {
	req, err := GetRequest(ctx)
	if err != nil {
		t.Error(err)
		return ctx
	}
	resp, err := th.tHTTP.MakeRequest(req)
	if err != nil {
		t.Error(err)
		return ctx
	}

	ctx.Set(ResponseKey{}, resp)
	return ctx
}

func (th testHTTPMethods) makeRequest(t StepTest, ctx context.Context, method, url string) context.Context {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		t.Error(err)
		return ctx
	}

	resp, err := th.tHTTP.MakeRequest(req)
	if err != nil {
		t.Error(err)
		return ctx
	}

	ctx.Set(ResponseKey{}, resp)
	return ctx
}

func (th testHTTPMethods) statusCodeEquals(t StepTest, ctx context.Context, expectedStatus int) context.Context {
	resp, err := GetResponse(ctx)
	if err != nil {
		t.Error(err)
		return ctx
	}

	if expectedStatus != resp.Code {
		t.Error(fmt.Errorf("expected status code: %d but %d given", expectedStatus, resp.Code))
		return ctx
	}
	return ctx
}

func GetResponse(ctx context.Context) (Response, error) {
	v, err := ctx.Get(ResponseKey{})
	if err != nil {
		return Response{}, err
	}

	return v.(Response), nil
}

func GetRequest(ctx context.Context) (*http.Request, error) {
	v, err := ctx.Get(RequestKey{})
	if err != nil {
		return nil, err
	}

	return v.(*http.Request), nil
}
