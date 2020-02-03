package gobdd

import (
	"context"
	"net/http"
	"testing"

	"github.com/go-bdd/assert"

	"github.com/go-bdd/gobdd/testhttp"
)

func TestHTTP(t *testing.T) {
	s := NewSuite(t, NewSuiteOptions().WithFeaturesPath("features/http.feature"))
	router := http.NewServeMux()
	router.HandleFunc("/health", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	router.HandleFunc("/json", func(w http.ResponseWriter, req *http.Request) {
		_, _ = w.Write([]byte(`{"valid": "json"}`))
		w.WriteHeader(http.StatusOK)
	})

	router.HandleFunc("/mirror", func(w http.ResponseWriter, r *http.Request) {
		for key, values := range r.Header {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}
		w.WriteHeader(http.StatusOK)
	})

	s.AddStep(`^the request has method set to (GET|POST|PUT|DELETE|OPTIONS)$`, requestHasMethodSetTo)
	s.AddStep(`^the url is set to "([^"]*)"$`, urlIsSetTo)
	s.AddStep(`^the request body is nil$`, requestBodyIsNil)
	s.AddStep(`^I set the header "([^"]*)" to "([^"]*)"$`, ISetHeaderTo)
	s.AddStep(`^the request has header "([^"]*)" set to "([^"]*)"$`, requestHasHeaderSetTo)
	testhttp.Build(s, router)

	s.Run()
}

func requestHasMethodSetTo(ctx context.Context, method string) (context.Context, error) {
	r := ctx.Value(testhttp.RequestKey{}).(*http.Request)
	return ctx, assert.Equals(method, r.Method)
}

func requestBodyIsNil(ctx context.Context) (context.Context, error) {
	r := ctx.Value(testhttp.RequestKey{}).(*http.Request)
	return ctx, assert.Nil(r.Body)
}

func urlIsSetTo(ctx context.Context, url string) (context.Context, error) {
	r := ctx.Value(testhttp.RequestKey{}).(*http.Request)
	return ctx, assert.Equals(url, r.URL.String())
}

func ISetHeaderTo(ctx context.Context, headerName, value string) (context.Context, error) {
	r := ctx.Value(testhttp.RequestKey{}).(*http.Request)
	r.Header.Set(headerName, value)
	ctx = context.WithValue(ctx, testhttp.RequestKey{}, r)
	return ctx, nil
}

func requestHasHeaderSetTo(ctx context.Context, headerName, expected string) (context.Context, error) {
	r := ctx.Value(testhttp.RequestKey{}).(*http.Request)
	given := r.Header.Get(headerName)
	return ctx, assert.Equals(expected, given)
}
