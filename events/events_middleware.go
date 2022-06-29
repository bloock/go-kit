package events

import (
	"fmt"
	"github.com/bloock/go-kit/request"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
)

func MiddlewareEvents() gin.HandlerFunc {
	return func(c *gin.Context) {
		b, err := io.ReadAll(c.Request.Body)
		log.Println(err)
		body := map[string]interface{}{
			"type":          "",
			"status":        c.Writer.Status(),
			"method":        c.Request.Method,
			"path":          c.Request.URL,
			"request_id":    c.Request.Header.Get("X-Request-ID"),
			"request_body":  b,
			"response_body": "",
			"ip":            c.ClientIP(),
			"user_id":       c.Request.Header.Get("x-user-id"),
		}

		url := fmt.Sprintf("%s://%s/%s", "https", "api.bloock.dev", "events/v1/activities")
		headers := make(map[string]string)

		log.Printf("Info --> %+v", body)
		err = request.RestClient{}.PostWithHeaders(url, body, nil, headers)

		if err != nil {
			c.Writer.WriteHeader(http.StatusBadRequest)
			c.Writer.Write([]byte(fmt.Sprintf("Events: Could not register event. %s", err.Error())))
			c.Abort()
			return
		}

		return
	}
}
