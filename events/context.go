package events

import (
	"github.com/gin-gonic/gin"
	"io"
)

type ResponseContext struct {
	Type         string
	Status       int
	Path         string
	RequestID    string
	RequestBody  string
	ResponseBody string
	IP           string
	UserID       string
}

func NewResponseContext(c *gin.Context, w wrappedWriter, typ string) (*ResponseContext, error) {

	requestBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return &ResponseContext{}, err
	}

	return &ResponseContext{
		Type:         typ,
		Status:       w.Status(),
		Path:         c.Request.URL.Path,
		RequestID:    c.Request.Header.Get("X-Request-ID"),
		RequestBody:  string(requestBody),
		ResponseBody: w.body.String(),
		IP:           c.ClientIP(),
		UserID:       c.Request.Header.Get("x-user-id"),
	}, nil

}
