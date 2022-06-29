package events

import (
	"bytes"
	"fmt"
	"github.com/bloock/go-kit/request"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"log"
)

type MiddlewareEvent struct {
	logger zerolog.Logger
}

func NewMiddlewareEvent(l zerolog.Logger) MiddlewareEvent {
	return MiddlewareEvent{
		logger: l,
	}
}

func (m MiddlewareEvent) MiddlewareEvents(typ string) gin.HandlerFunc {
	return func(c *gin.Context) {
		writer := c.Writer
		rw := wrappedWriter{ResponseWriter: c.Writer}
		c.Writer = &rw
		c.Next()
		c.Writer = writer

		body, err := NewResponseContext(c, rw, typ)

		url := fmt.Sprintf("%s://%s/%s", "https", "api.bloock.dev", "events/v1/activities")

		log.Printf("Info --> %+v", body)
		err = request.RestClient{}.Post(url, body, nil)

		if err != nil {
			m.logger.Error().Msgf("Events: Could not register event. %s", err.Error())
			return
		}
	}
}

type wrappedWriter struct {
	gin.ResponseWriter
	body bytes.Buffer
}

func (rw *wrappedWriter) Write(body []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(body)
	if err == nil {
		rw.body.Write(body)
	}
	return n, err
}
