package request

import (
	"bytes"
	"net/http"
)

type RestClient struct{}

func (r RestClient) Get(url string) (*http.Response, error) {
	return http.Get(url)
}

func (r RestClient) Post(url, contentType string, body []byte) (*http.Response, error) {
	return http.Post(url, contentType, bytes.NewBuffer(body))
}
