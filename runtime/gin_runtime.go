package runtime

import (
	"context"
	"net/http"
	"time"

	pinned "github.com/bloock/go-kit/http/versioning"

	"github.com/bloock/go-kit/client"
	bloockContext "github.com/bloock/go-kit/context"
	"github.com/bloock/go-kit/http/middleware"
	"github.com/bloock/go-kit/observability"
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	gintrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/gin-gonic/gin"
)

var api_version_header = "api_version"

type GinRuntime struct {
	name           string
	client         *client.GinEngine
	shutdownTime   time.Duration
	versionManager *pinned.VersionManager
	logger         observability.Logger
}

func NewVersionedGinRuntime(name string, c *client.GinEngine, shutdownTime time.Duration, vm *pinned.VersionManager, l observability.Logger) (*GinRuntime, error) {

	runtime, err := NewGinRuntime(name, c, shutdownTime, l)
	if err != nil {
		return nil, err
	}
	runtime.versionManager = &pinned.VersionManager{
		Layout: vm.Layout,
		Header: api_version_header,
	}

	return runtime, nil
}
func NewGinRuntime(name string, c *client.GinEngine, shutdownTime time.Duration, l observability.Logger) (*GinRuntime, error) {
	e := GinRuntime{
		name:         name,
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
		logger.WithLogger(func(c *gin.Context, l zerolog.Logger) zerolog.Logger {
			return l.With().Str("user-id", bloockContext.GetUserID(c)).
				Str("request-id", bloockContext.GetRequestID(c)).
				Logger()
		}),
	)
	e.client.Engine().Use(l)
	e.client.Engine().Use(middleware.ErrorMiddleware())
	e.client.Engine().Use(middleware.ContextMiddleware())
	e.client.Engine().Use(gintrace.Middleware(
		e.name,
		gintrace.WithIgnoreRequest(func(c *gin.Context) bool {
			return c.Request.URL.Path == "/health"
		}),
	))
	f(e.client.Engine())
}

func (e *GinRuntime) Run(ctx context.Context) {

	srv := &http.Server{
		Addr:    e.client.Addr(),
		Handler: e.client.Engine(),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			e.logger.Info(ctx).Msgf("server running on %s", e.client.Addr())
		}
	}()

	<-ctx.Done()

	ctx, cancel := context.WithTimeout(context.Background(), e.shutdownTime)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		e.logger.Info(ctx).Msgf("error while closing gin runtime: %s", err.Error())
	} else {
		e.logger.Info(ctx).Msg("gin runtime closed successfully")
	}
}
