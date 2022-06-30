package events

import (
	"bytes"
	"fmt"
	"github.com/bloock/go-kit/request"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"io"
	"log"
)

type MiddlewareEvent struct {
	logger zerolog.Logger
	endpointConfig EndpointConfig
}

func NewMiddlewareEvent(l zerolog.Logger, protocol, host string, port int, path string) MiddlewareEvent {
	return MiddlewareEvent{
		logger: l,
		endpointConfig: EndpointConfig{protocol, host, port, path},
	}
}

func (m MiddlewareEvent) MiddlewareEvents(typ string) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestBody, err := io.ReadAll(c.Request.Body)
		if err != nil {
			m.logger.Error().Msg(err.Error())
		}
		c.Request.Body.Close()
		c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))

		writer := c.Writer
		rw := wrappedWriter{ResponseWriter: c.Writer}
		c.Writer = &rw
		c.Next()
		c.Writer = writer

		body, err := NewResponseContext(c, rw, typ, string(requestBody))

		url := fmt.Sprintf("%s://%s:%d%s", m.endpointConfig.Protocol, m.endpointConfig.Host,
			m.endpointConfig.Port, m.endpointConfig.Path)

		log.Printf("Info --> %+v", body)
		err = request.RestClient{}.Post(url, body, nil)

		if err != nil {
			m.logger.Error().Msgf("Events: Could not register event. %s", err.Error())
			return
		}
	}
}

type EndpointConfig struct {
	Protocol string
	Host     string
	Port     int
	Path     string
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
