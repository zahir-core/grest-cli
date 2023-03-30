package app

import (
	"net/http"
	"time"

	"grest.dev/grest"
)

func HttpClient(method, url string) HttpClientInterface {
	hc := &httpClientUtil{}
	hc.Method = method
	hc.Url = url
	return hc
}

type HttpClientInterface interface {
	Debug()
	AddHeader(key, value string)
	AddMultipartBody(body any) error
	AddUrlEncodedBody(body any) error
	AddJsonBody(body any) error
	AddXmlBody(body any) error
	SetTimeout(timeout time.Duration)
	Send() (*http.Response, error)
	BodyResponseStr() string
	UnmarshalJson(v any) error
	UnmarshalXml(v any) error
}

// httpClientUtil implement HttpClientInterface embed from grest.httpClientUtil for simplicity
type httpClientUtil struct {
	grest.HttpClient
}
