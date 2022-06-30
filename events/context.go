package events

import (
	"github.com/gin-gonic/gin"
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

func NewResponseContext(c *gin.Context, w wrappedWriter, typ, requestBody string) (*ResponseContext, error) {

	return &ResponseContext{
		Type:         typ,
		Status:       w.Status(),
		Path:         c.Request.URL.Path,
		RequestID:    c.Request.Header.Get("X-Request-ID"),
		RequestBody:  requestBody,
		ResponseBody: w.body.String(),
		IP:           c.ClientIP(),
		UserID:       c.Request.Header.Get("x-user-id"),
	}, nil

}
