package gobdd

import (
	"github.com/go-bdd/assert"
	"github.com/go-bdd/gobdd/context"
	"net/http"
	"testing"

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

	_ = s.AddStep(`^the request has method set to (GET|POST|PUT|DELETE|OPTIONS)$`, requestHasMethodSetTo)
	_ = s.AddStep(`^the url is set to "([^"]*)"$`, urlIsSetTo)
	_ = s.AddStep(`^the request body is nil$`, requestBodyIsNil)
	_ = s.AddStep(`^I set the header "([^"]*)" to "([^"]*)"$`, ISetHeaderTo)
	_ = s.AddStep(`^the request has header "([^"]*)" set to "([^"]*)"$`, requestHasHeaderSetTo)
	testhttp.Build(s, router)

	s.Run()
}

func requestHasMethodSetTo(ctx context.Context, method string) error {
	r := ctx.Get(testhttp.RequestKey{}).(*http.Request)
	return assert.Equals(method, r.Method)
}

func requestBodyIsNil(ctx context.Context) error {
	r := ctx.Get(testhttp.RequestKey{}).(*http.Request)
	return assert.Nil(r.Body)
}

func urlIsSetTo(ctx context.Context, url string) error {
	r := ctx.Get(testhttp.RequestKey{}).(*http.Request)
	return assert.Equals(url, r.URL.String())
}

func ISetHeaderTo(ctx context.Context, headerName, value string) error {
	r := ctx.Get(testhttp.RequestKey{}).(*http.Request)
	r.Header.Set(headerName, value)
	ctx.Set(testhttp.RequestKey{}, r)
	return nil
}

func requestHasHeaderSetTo(ctx context.Context, headerName,expected  string) error {
	r := ctx.Get(testhttp.RequestKey{}).(*http.Request)
	given := r.Header.Get(headerName)
	return assert.Equals(expected, given)
}
