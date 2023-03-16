package app

import "grest.dev/grest"

func HttpClient(method, url string) *httpClientImpl {
	hc := &httpClientImpl{}
	hc.Method = method
	hc.Url = url
	return hc
}

type HttpClientInterface interface {
	grest.HttpClientInterface
}

// httpClientImpl implement HttpClientInterface embed from grest.httpClientImpl for simplicity
type httpClientImpl struct {
	grest.HttpClient
}
