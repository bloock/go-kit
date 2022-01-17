package request

//go:generate mockgen -package=mocks -source=request/http_request.go -destination mocks/mock_http_request.go
type HttpRequest interface {
	Get(url string, response interface{}) error
	Post(url string, body interface{}, response interface{}) error
	PostWithHeaders(url string, body, response interface{}, headers map[string]string) error
}
