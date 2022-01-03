package request

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

type RestClient struct{}

func (r RestClient) Get(url string) (*http.Response, error) {
	return http.Get(url)
}

func (r RestClient) Post(url string, body interface{}) ([]byte, error) {
	b, err := json.Marshal(body)
	if err != nil {
		return []byte{}, err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(b))
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()

	respByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return []byte{}, errors.New("invalid status code")
	}

	return respByte, nil
}
