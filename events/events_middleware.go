package events

import (
	"bytes"
	"github.com/bloock/go-kit/event"
	event_entity "github.com/bloock/go-kit/event/entity"
	"github.com/bloock/go-kit/publisher"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"io"
)

type MiddlewareEvent struct {
	logger    zerolog.Logger
	publisher publisher.Publisher
}

func NewMiddlewareEvent(l zerolog.Logger, publisher publisher.Publisher) MiddlewareEvent {
	return MiddlewareEvent{
		logger:    l,
		publisher: publisher,
	}
}

func (m MiddlewareEvent) MiddlewareEvents(typ string) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestBody, err := io.ReadAll(c.Request.Body)
		if err != nil {
			m.logger.Error().Msg(err.Error())
			return
		}
		c.Request.Body.Close()
		c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))

		writer := c.Writer
		rw := wrappedWriter{ResponseWriter: c.Writer}
		c.Writer = &rw
		c.Next()
		c.Writer = writer

		body, err := NewResponseContext(c, rw, typ, string(requestBody))

		eventBody := event_entity.NewEventsActivityCreateEntity(
			body.Type,
			body.Status,
			body.Path,
			body.RequestID,
			body.RequestBody,
			body.ResponseBody,
			body.IpAddress,
			body.UserID,
			body.Method)

		ev, err := event.NewEntityEvent(event.EventsActivityCreated, eventBody)
		if err != nil {
			m.logger.Error().Msgf("Events: Could not register event. %s", err.Error())
			return
		}
		m.publisher.Publish(ev, nil)
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
