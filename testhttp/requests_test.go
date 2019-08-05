package testhttp

import (
	"bytes"
	"net/http"
	"reflect"
	"testing"
)

func TestValidMethods(t *testing.T) {
	handler := Build(testHandler{})
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

func TestInvalidMethod(t *testing.T) {
	method := "NOT EXISTS"
	handler := TestHTTP{}
	_, err := handler.Request(method, "", nil)
	if err == nil {
		t.Error("the handler should not allow for invalid HTTP methods")
	}
}

type testHandler struct {
}

func (h testHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
}
