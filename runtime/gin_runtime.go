package runtime

import (
	"context"
	"net/http"
	"time"

	"github.com/bloock/go-kit/client"
	"github.com/gin-gonic/gin"
	openApiMiddleware "github.com/go-openapi/runtime/middleware"
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
	f(e.client.Engine())
	e.enableSwagger()
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

func (e *GinRuntime) enableSwagger() {
	o := openApiMiddleware.SwaggerUIOpts{SpecURL: "/swagger.yaml"}
	sw := openApiMiddleware.SwaggerUI(o, nil)
	e.client.Engine().GET("/docs", gin.WrapH(sw))
	e.client.Engine().GET("/swagger.yaml", func(c *gin.Context) {
		c.File("./swagger.yaml")
	})
}
