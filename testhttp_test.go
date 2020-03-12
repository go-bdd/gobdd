package gobdd

import (
	"net/http"
	"testing"

	"github.com/go-bdd/assert"
	"github.com/go-bdd/gobdd/context"
	"github.com/go-bdd/gobdd/testhttp"
)

func TestHTTP(t *testing.T) {
	s := NewSuite(t, WithFeaturesPath("features/http.feature"))
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

func requestHasMethodSetTo(t StepTest, ctx context.Context, method string) context.Context {
	r, err := testhttp.GetRequest(ctx)
	if err != nil {
		t.Error(err)
		return ctx
	}

	if err = assert.Equals(method, r.Method); err != nil {
		t.Error(err)
	}
	return ctx
}

func requestBodyIsNil(t StepTest, ctx context.Context) context.Context {
	r, err := testhttp.GetRequest(ctx)
	if err != nil {
		t.Error(err)
		return ctx
	}
	if err = assert.Nil(r.Body); err != nil {
		t.Error(err)
	}
	return ctx
}

func urlIsSetTo(t StepTest, ctx context.Context, url string) context.Context {
	r, err := testhttp.GetRequest(ctx)
	if err != nil {
		t.Error(err)
		return ctx
	}

	if err := assert.Equals(url, r.URL.String()); err != nil {
		t.Error(err)
	}
	return ctx
}

func ISetHeaderTo(t StepTest, ctx context.Context, headerName, value string) context.Context {
	r, err := testhttp.GetRequest(ctx)
	if err != nil {
		t.Error(err)
		return ctx
	}
	r.Header.Set(headerName, value)
	ctx.Set(testhttp.RequestKey{}, r)
	return ctx
}

func requestHasHeaderSetTo(t StepTest, ctx context.Context, headerName, expected string) context.Context {
	r, err := testhttp.GetRequest(ctx)
	if err != nil {
		t.Error(err)
		return ctx
	}
	given := r.Header.Get(headerName)
	if err := assert.Equals(expected, given); err != nil {
		t.Error(err)
	}
	return ctx
}
