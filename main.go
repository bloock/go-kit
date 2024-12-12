package main

import (
	"context"
	gin_middleware "github.com/bloock/go-kit/http/middleware/gin"
	"github.com/bloock/go-kit/http/presenters"
	pinned "github.com/bloock/go-kit/http/versioning"
	"github.com/bloock/go-kit/test_utils"
	"github.com/gin-gonic/gin"
	"github.com/go-chi/chi/v5"
	"net/http"
	"sync"
	"time"

	"github.com/bloock/go-kit/client"
	"github.com/bloock/go-kit/observability"
	"github.com/bloock/go-kit/runtime"
)

func main() {
	ctx := context.Background()

	l := observability.InitLogger("local", "test_service", "1.0.0", true)
	err := observability.InitTracer(ctx, "connection", "dev", "1.0.0", observability.Logger{})

	wg := sync.WaitGroup{}
	wg.Add(2)

	cronClient, err := client.NewCronClient(context.Background(), l)
	if err != nil {
		panic(err.Error())
	}
	cronRuntime, err := runtime.NewCronRuntime(cronClient, 5*time.Second, l)
	if err != nil {
		panic(err.Error())
	}
	cronRuntime.AddHandler("test", 5*time.Second, CronHandler())

	go func() {
		defer wg.Done()
		cronRuntime.Run(context.Background())
	}()

	ginClient := client.NewGinEngine("0.0.0.0", 8080, true, l)
	chiClient := client.NewChiRouter("0.0.0.0", 8001, true, l)
	var vm = &pinned.VersionManager{
		Layout: "2006-01-02",
	}

	ginRuntime, err := runtime.NewVersionedGinRuntime("service-gin", ginClient, 5*time.Second, vm, l)
	if err != nil {
		panic(err.Error())
	}

	ginRuntime.SetHandlers(func(e *gin.Engine) {

		e.GET("/test", gin_middleware.HandlerVersioning(vm, test_utils.TestHandlerInstance.Versions()), test_utils.TestHandlerInstance.Handler())
	})

	chiRuntime, err := runtime.NewVersionedChiRuntime("service-chi", chiClient, 5*time.Second, vm, l)
	if err != nil {
		panic(err.Error())
	}

	chiRuntime.SetHandlers(ctx, func(r *chi.Mux) {
		r.Get("/test", func(writer http.ResponseWriter, request *http.Request) {
			presenters.RenderJSON(writer, request, http.StatusOK, "Hello World")
			return
		})

		r.Post("/test", func(writer http.ResponseWriter, request *http.Request) {
			presenters.RenderJSON(writer, request, http.StatusBadRequest, "post method error")
			return
		})

		r.Patch("/test", func(writer http.ResponseWriter, request *http.Request) {
			presenters.RenderJSON(writer, request, http.StatusInternalServerError, "patch method error")
			return
		})
	})

	// Run API server
	go func() {
		defer wg.Done()
		chiRuntime.Run(ctx)
	}()

	wg.Wait()

}

func CronHandler() client.CronHandler {
	return func(c context.Context) error {
		ctx := context.Background()

		l := observability.InitLogger("local", "test_service", "1.0.0", true)
		l.Error(ctx).Str("t", "user").Msg("a cron message")
		return nil
	}
}
