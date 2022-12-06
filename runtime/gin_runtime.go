package runtime

import (
	"bytes"
	"context"
	httperror "github.com/bloock/go-kit/errors/http"
	"io"
	"net/http"
	"time"

	"github.com/bloock/go-kit/client"
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type GinRuntime struct {
	client       *client.GinEngine
	shutdownTime time.Duration
	logger       zerolog.Logger
}

func NewGinRuntime(c *client.GinEngine, shutdownTime time.Duration, l zerolog.Logger) (*GinRuntime, error) {
	e := GinRuntime{
		client:       c,
		shutdownTime: shutdownTime,
		logger:       l,
	}

	return &e, nil
}

func (e *GinRuntime) SetHandlers(f func(*gin.Engine)) {
	l := logger.SetLogger(
		logger.WithSkipPath([]string{"/health"}),
		logger.WithUTC(true),
		logger.WithLogger(func(c *gin.Context, _ io.Writer, latency time.Duration) zerolog.Logger {
			bodyReader := c.Request.Body
			var body []byte
			_, _ = bodyReader.Read(body)
			str := bytes.NewBuffer(body).String()
			xClientID := c.Request.Header.Get("X-Client-ID")
			xRequestID := c.Request.Header.Get("X-Request-ID")
			c.Set("X-Request-ID", xRequestID)
			return e.logger.With().
				Int("status", c.Writer.Status()).
				Str("method", c.Request.Method).
				Str("path", c.Request.URL.Path).
				Str("body", str).
				Str("ip", c.ClientIP()).
				Dur("latency", latency).
				Str("user_agent", c.Request.UserAgent()).
				Str("x-client-id", xClientID).
				Str("x-request-id", xRequestID).
				Logger()
		}),
	)
	e.client.Engine().Use(l)

	e.client.Engine().Use(httperror.ErrorMiddleware())
	f(e.client.Engine())
}

func (e *GinRuntime) Run(ctx context.Context) {
	srv := &http.Server{
		Addr:    e.client.Addr(),
		Handler: e.client.Engine(),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			e.logger.Info().Msgf("server running on %s", e.client.Addr())
		}
	}()

	<-ctx.Done()

	ctx, cancel := context.WithTimeout(context.Background(), e.shutdownTime)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		e.logger.Info().Msgf("error while closing gin runtime: %s", err.Error())
	} else {
		e.logger.Info().Msg("gin runtime closed successfully")
	}
}
