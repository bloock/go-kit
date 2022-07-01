package events

import (
	httperror "github.com/bloock/go-kit/http_error"
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
		var parsedError *httperror.AppError

		switch err.(type) {
		case *httperror.AppError:
			parsedError = err.(*httperror.AppError)
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
		UserID:       c.Request.Header.Get("x-user-id"),
		Method:       c.Request.Method,
	}, nil
}
