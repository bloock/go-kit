package request

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	bloockContext "github.com/bloock/go-kit/context"
)

type RestClient struct{}

func (r RestClient) Get(ctx context.Context, url string, response interface{}) error {
	client := http.Client{}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req = r.setRequestID(ctx, req)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respByte, err := io.ReadAll(resp.Body)
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

func (r RestClient) Post(ctx context.Context, url string, body interface{}, response interface{}, contentType string) error {
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}

	client := http.Client{}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	req = r.setRequestID(ctx, req)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respByte, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New(resp.Status)
	}

	err = json.Unmarshal(respByte, &response)
	if err != nil {
		return err
	}

	return nil
}

func (r RestClient) PostWithHeaders(ctx context.Context, url string, body, response interface{}, headers map[string]string, contentType string) error {
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
	req.Header.Set("Content-Type", contentType)
	req = r.setRequestID(ctx, req)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respByte, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New(resp.Status)
	}

	err = json.Unmarshal(respByte, &response)
	if err != nil {
		return err
	}

	return nil
}

func (r RestClient) Delete(ctx context.Context, url string, body interface{}, response interface{}, headers map[string]string) error {
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
	req = r.setRequestID(ctx, req)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respByte, err := io.ReadAll(resp.Body)
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

func (r RestClient) setRequestID(ctx context.Context, req *http.Request) *http.Request {
	req.Header.Set("X-Request-ID", bloockContext.GetRequestID(ctx))
	return req
}
