package request

import "net/http"

type RestClient struct{}

func (r RestClient) Get(url string) (*http.Response, error) {
	return http.Get(url)
}
