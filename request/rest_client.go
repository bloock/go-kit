package request

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

type RestClient struct{}

func (r RestClient) Get(url string, response interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New("invalid status code")
	}

	err = json.Unmarshal(respByte, &response)
	if err != nil {
		return err
	}

	return nil
}

func (r RestClient) Post(url string, body interface{}, response interface{}) error {
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New("invalid status code")
	}

	err = json.Unmarshal(respByte, &response)
	if err != nil {
		return err
	}

	return nil
}

func (r RestClient) PostWithHeaders(url string, body, response interface{}, headers map[string]string) error {
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}

	client := http.Client{}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(b))
	if err != nil {
		return err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New("invalid status code")
	}

	err = json.Unmarshal(respByte, &response)
	if err != nil {
		return err
	}

	return nil
}

func (r RestClient) Delete(url string, body interface{}, response interface{}, headers map[string]string) error {
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}

	client := http.Client{}
	req, err := http.NewRequest(http.MethodDelete, url, bytes.NewBuffer(b))
	if err != nil {
		return err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New("invalid status code")
	}

	err = json.Unmarshal(respByte, &response)
	if err != nil {
		return err
	}

	return nil
}
