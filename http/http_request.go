package http

import (
	"context"
)

//go:generate mockgen -package=mocks -source=request/http_request.go -destination mocks/request/mock_http_request.go
type HttpRequest interface {
	Get(ctx context.Context, url string, response interface{}) error
	GetWithHeaders(ctx context.Context, url string, response interface{}, headers map[string][]string) error
	Post(ctx context.Context, url string, body interface{}, response interface{}, contentType string) error
	PostWithHeaders(ctx context.Context, url string, body, response interface{}, headers map[string]string, contentType string) error
	Delete(ctx context.Context, url string, body interface{}, response interface{}, headers map[string]string) error
}
