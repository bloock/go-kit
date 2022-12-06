package events

import (
	http2 "github.com/bloock/go-kit/errors/http"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ResponseContext struct {
	Type         string
	Status       int
	Path         string
	RequestID    string
	RequestBody  string
	ResponseBody string
	IpAddress    string
	UserID       string
	Method       string
}

func NewResponseContext(c *gin.Context, w wrappedWriter, typ, requestBody string) (*ResponseContext, error) {
	status := w.Status()
	responseBody := w.body.String()
	if len(c.Errors) > 0 {
		detectedErrors := c.Errors.ByType(gin.ErrorTypeAny)
		err := detectedErrors[0].Err
		var parsedError *http2.HttpAppError

		switch err.(type) {
		case *http2.HttpAppError:
			parsedError = err.(*http2.HttpAppError)
			status = parsedError.Code
			responseBody = parsedError.Message
		default:
			status = http.StatusInternalServerError
			responseBody = "Internal Server Error"
		}
	}

	return &ResponseContext{
		Type:         typ,
		Status:       status,
		Path:         c.Request.URL.Path,
		RequestID:    c.Request.Header.Get("X-Request-ID"),
		RequestBody:  requestBody,
		ResponseBody: responseBody,
		IpAddress:    c.ClientIP(),
		UserID:       c.Request.Header.Get("X-User-ID"),
		Method:       c.Request.Method,
	}, nil
}
