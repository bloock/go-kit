package middleware

import (
	"bytes"
	"fmt"
	"github.com/bloock/go-kit/context"
	"io"
	"net/http"

	"github.com/bloock/go-kit/domain"
	"github.com/bloock/go-kit/domain/event_entity"
	"github.com/bloock/go-kit/errors"
	"github.com/bloock/go-kit/observability"
	"github.com/bloock/go-kit/repository/publisher"
	"github.com/gin-gonic/gin"
)

type MiddlewareEvent struct {
	logger    observability.Logger
	publisher publisher.Publisher
}

func NewMiddlewareEvent(l observability.Logger, publisher publisher.Publisher) MiddlewareEvent {
	return MiddlewareEvent{
		logger:    l,
		publisher: publisher,
	}
}

func (m MiddlewareEvent) MiddlewareEvents(typ string) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestBody, err := io.ReadAll(c.Request.Body)
		if err != nil {
			m.logger.Error(c).Msg(err.Error())
			return
		}
		c.Request.Body.Close()
		c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))

		writer := c.Writer
		rw := eventsWrappedWriter{ResponseWriter: c.Writer}
		c.Writer = &rw
		c.Next()
		c.Writer = writer

		body := NewResponseContext(c, rw, typ, string(requestBody))

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

		ev, err := domain.NewEntityEvent(domain.EventsActivityCreated, eventBody)
		if err != nil {
			m.logger.Error(c).Msgf("Events: Could not register event. %s", err.Error())
			return
		}
		m.publisher.Publish(c, ev, nil)
	}
}

type eventsWrappedWriter struct {
	gin.ResponseWriter
	body bytes.Buffer
}

func (rw *eventsWrappedWriter) Write(body []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(body)
	if err == nil {
		rw.body.Write(body)
	}
	return n, err
}

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

func NewResponseContext(c *gin.Context, w eventsWrappedWriter, typ, requestBody string) *ResponseContext {
	status := w.Status()
	responseBody := w.body.String()
	if len(c.Errors) > 0 {
		detectedErrors := c.Errors.ByType(gin.ErrorTypeAny)
		err := detectedErrors[0].Err
		var parsedError *errors.HttpAppError

		switch err.(type) {
		case *errors.HttpAppError:
			parsedError = err.(*errors.HttpAppError)
			status = parsedError.Code
			responseBody = parsedError.Message
		default:
			status = http.StatusInternalServerError
			responseBody = "Internal Server Error"
		}
	}
	userID := ""
	uid, ok := c.Get(context.UserIDKey)
	if ok {
		userID = uid.(string)
	}

	return &ResponseContext{
		Type:         typ,
		Status:       status,
		Path:         c.Request.URL.Path,
		RequestID:    c.Request.Header.Get("X-Request-ID"),
		RequestBody:  requestBody,
		ResponseBody: responseBody,
		IpAddress:    c.ClientIP(),
		UserID:       userID,
		Method:       c.Request.Method,
	}
}

type EventType int32
type EventTypeTest int32

const (
	GetAnchor EventType = iota
	GetProof
	NewMessages
)

const (
	GetAnchorTest EventTypeTest = iota
	GetProofTest
	NewMessagesTest
)

const (
	CorePrefix     = "core"
	CoreTestPrefix = "core-test"
)

func (e EventType) Str() string {
	switch e {
	case GetAnchor:
		return fmt.Sprintf("%s.%s", CorePrefix, "get_anchor")
	case GetProof:
		return fmt.Sprintf("%s.%s", CorePrefix, "get_proof")
	case NewMessages:
		return fmt.Sprintf("%s.%s", CorePrefix, "new_messages")
	default:
		return ""
	}
}

func (e EventTypeTest) Str() string {
	switch e {
	case GetAnchorTest:
		return fmt.Sprintf("%s.%s", CoreTestPrefix, "get_anchor")
	case GetProofTest:
		return fmt.Sprintf("%s.%s", CoreTestPrefix, "get_proof")
	case NewMessagesTest:
		return fmt.Sprintf("%s.%s", CoreTestPrefix, "new_messages")
	default:
		return ""
	}
}
