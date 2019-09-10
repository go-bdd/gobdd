package testhttp

import (
	"bytes"
	"net/http"
	"reflect"
	"testing"

	"github.com/go-bdd/gobdd"
)

func TestValidMethods(t *testing.T) {
	s := gobdd.NewSuite(t, gobdd.NewSuiteOptions())
	handler := Build(s, testHandler{})
	methods := []string{"Get", "Post", "Trace", "Options", "Head", "Connect", "Patch", "Put", "Delete"}

	for _, method := range methods {
		h := reflect.ValueOf(&handler)
		m := h.MethodByName(method)
		body := bytes.NewReader(nil)
		values := m.Call([]reflect.Value{reflect.ValueOf(""), reflect.ValueOf(body)})
		if !values[1].IsNil() {
			t.Errorf("the function should not return an error")
		}
	}
}

func TestGettingResponseHeaders(t *testing.T) {
	s := gobdd.NewSuite(t, gobdd.NewSuiteOptions())
	handler := Build(s, testHandler{})
	//handler.
}

type testHandler struct {
	body []byte
}

func (h testHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for key, values := range r.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	_, _ = w.Write(h.body)
}
