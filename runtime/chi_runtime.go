package runtime

import (
	"context"
	"github.com/bloock/go-kit/client"
	chi_middleware "github.com/bloock/go-kit/http/middleware/chi"
	pinned "github.com/bloock/go-kit/http/versioning"
	"github.com/bloock/go-kit/observability"
	sentryhttp "github.com/getsentry/sentry-go/http"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"net/http"
	"time"
)

type ChiRuntime struct {
	name           string
	client         *client.ChiRouter
	shutdownTime   time.Duration
	versionManager *pinned.VersionManager
	logger         observability.Logger
}

func NewVersionedChiRuntime(name string, c *client.ChiRouter, shutdownTime time.Duration, vm *pinned.VersionManager, l observability.Logger) (*ChiRuntime, error) {
	runtime, err := NewChiRuntime(name, c, shutdownTime, l)
	if err != nil {
		return nil, err
	}
	runtime.versionManager = &pinned.VersionManager{
		Layout: vm.Layout,
		Header: api_version_header,
	}

	return runtime, nil
}

func NewChiRuntime(name string, c *client.ChiRouter, shutdownTime time.Duration, l observability.Logger) (*ChiRuntime, error) {
	e := ChiRuntime{
		name:         name,
		client:       c,
		shutdownTime: shutdownTime,
		logger:       l,
	}

	return &e, nil
}

func (e *ChiRuntime) SetHandlers(ctx context.Context, f func(*chi.Mux)) {

	options := cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{
			http.MethodHead,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	}

	sentryHandler := sentryhttp.New(sentryhttp.Options{Repanic: true})

	e.client.Router().Use(chi_middleware.ContextMiddleware)
	e.client.Router().Use(chi_middleware.Logger(e.logger))
	e.client.Router().Use(cors.Handler(options))
	e.client.Router().Use(sentryHandler.Handle)

	f(e.client.Router())
}

func (e *ChiRuntime) Run(ctx context.Context) {
	srv := &http.Server{
		Addr:    e.client.Addr(),
		Handler: e.client.Router(),
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

func (e *ChiRuntime) Client() *client.ChiRouter {
	return e.client
}
