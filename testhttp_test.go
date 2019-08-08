package gobdd

import (
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

	testhttp.Build(s, router)

	s.Run()
}
