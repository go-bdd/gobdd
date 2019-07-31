package testhttp

import "testing"

func TestInvalidMethod(t *testing.T) {
	method := "NOT EXISTS"
	handler := TestHTTP{}
	_, err := handler.Request(method, "", nil)
	if err == nil {
		t.Error("the handler should not allow for invalid HTTP methods")
	}
}
