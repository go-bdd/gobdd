package gobdd

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/go-bdd/gobdd/context"
	"github.com/go-bdd/gobdd/testhttp"
)

type httpResponse struct{}

type testHTTPMethods struct {
	tHTTP testhttp.TestHTTP
}

func TestHTTP(t *testing.T) {
	s := NewSuite(t, NewSuiteOptions().WithFeaturesPath("features/http.feature"))
	router := http.NewServeMux()
	router.HandleFunc("/health", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	tHTTP := testHTTPMethods{
		tHTTP: testhttp.Build(router),
	}

	_ = s.AddStep(`^I make a (GET|POST|PUT|DELETE|OPTIONS) request to "([^"]*)"$`, tHTTP.makeRequest)
	_ = s.AddStep(`^the response code equals (\d+)$`, tHTTP.statusCodeEquals)

	s.Run()
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
	resp := ctx.Get(httpResponse{}).(testhttp.Response)

	if expectedStatus != resp.Code {
		return fmt.Errorf("expected status code: %d but %d given", expectedStatus, resp.Code)
	}
	return nil
}
