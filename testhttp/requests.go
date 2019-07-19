package testhttp

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
)

type Response struct {
	Code   int
	Body   io.Reader
	Header http.Header
}

type httpHandler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type TestHTTP struct {
	handler httpHandler
}

func (thhtp TestHTTP) Trace(url string, body io.Reader) (Response, error) {
	return thhtp.Request(http.MethodTrace, url, body)
}

func (thhtp TestHTTP) Options(url string, body io.Reader) (Response, error) {
	return thhtp.Request(http.MethodOptions, url, body)
}

func (thhtp TestHTTP) Head(url string, body io.Reader) (Response, error) {
	return thhtp.Request(http.MethodHead, url, body)
}

func (thhtp TestHTTP) Connect(url string, body io.Reader) (Response, error) {
	return thhtp.Request(http.MethodConnect, url, body)
}

func (thhtp TestHTTP) Patch(url string, body io.Reader) (Response, error) {
	return thhtp.Request(http.MethodPatch, url, body)
}

func (thhtp TestHTTP) Post(url string, body io.Reader) (Response, error) {
	return thhtp.Request(http.MethodPost, url, body)
}

func (thhtp TestHTTP) Put(url string, body io.Reader) (Response, error) {
	return thhtp.Request(http.MethodPut, url, body)
}

func (thhtp TestHTTP) Delete(url string, body io.Reader) (Response, error) {
	return thhtp.Request(http.MethodDelete, url, body)
}

func (thhtp TestHTTP) Get(url string, body io.Reader) (Response, error) {
	return thhtp.Request(http.MethodGet, url, body)
}

func (thhtp TestHTTP) Request(method, url string, body io.Reader) (Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return Response{}, errors.New("invalid HTTP method type")
	}

	rr := httptest.NewRecorder()
	thhtp.handler.ServeHTTP(rr, req)

	return Response{
		Code:   rr.Code,
		Body:   rr.Body,
		Header: rr.Header(),
	}, nil
}
