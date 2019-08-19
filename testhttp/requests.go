package testhttp

import (
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
	req, _ := http.NewRequest(http.MethodTrace, url, body)
	return thhtp.MakeRequest(req)
}

func (thhtp TestHTTP) Options(url string, body io.Reader) (Response, error) {
	req, _ := http.NewRequest(http.MethodOptions, url, body)
	return thhtp.MakeRequest(req)
}

func (thhtp TestHTTP) Head(url string, body io.Reader) (Response, error) {
	req, _ := http.NewRequest(http.MethodHead, url, body)
	return thhtp.MakeRequest(req)
}

func (thhtp TestHTTP) Connect(url string, body io.Reader) (Response, error) {
	req, _ := http.NewRequest(http.MethodConnect, url, body)
	return thhtp.MakeRequest(req)
}

func (thhtp TestHTTP) Patch(url string, body io.Reader) (Response, error) {
	req, _ := http.NewRequest(http.MethodPatch, url, body)
	return thhtp.MakeRequest(req)
}

func (thhtp TestHTTP) Post(url string, body io.Reader) (Response, error) {
	req, _ := http.NewRequest(http.MethodPost, url, body)
	return thhtp.MakeRequest(req)
}

func (thhtp TestHTTP) Put(url string, body io.Reader) (Response, error) {
	req, _ := http.NewRequest(http.MethodPut, url, body)
	return thhtp.MakeRequest(req)
}

func (thhtp TestHTTP) Delete(url string, body io.Reader) (Response, error) {
	req, _ := http.NewRequest(http.MethodDelete, url, body)
	return thhtp.MakeRequest(req)
}

func (thhtp TestHTTP) Get(url string, body io.Reader) (Response, error) {
	req, _ := http.NewRequest(http.MethodGet, url, body)
	return thhtp.MakeRequest(req)
}

func (thhtp TestHTTP) MakeRequest(req *http.Request) (Response, error) {
	rr := httptest.NewRecorder()
	thhtp.handler.ServeHTTP(rr, req)

	return Response{
		Code:   rr.Code,
		Body:   rr.Body,
		Header: rr.Header(),
	}, nil
}
