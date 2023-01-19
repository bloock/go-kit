package main

import (
	"context"
	"sync"
	"time"

	"github.com/bloock/go-kit/client"
	"github.com/bloock/go-kit/observability"
	"github.com/bloock/go-kit/runtime"
	"github.com/gin-gonic/gin"
)

func main() {
	l := observability.InitLogger("local", "test_service", "1.0.0", true)
	tracer := observability.InitTracer("local", "test_service", "1.0.0", true)
	defer tracer.Stop()

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
	ginRuntime, err := runtime.NewGinRuntime("service-gin", ginClient, 5*time.Second, l)
	if err != nil {
		panic(err.Error())
	}
	ginRuntime.SetHandlers(func(e *gin.Engine) {
		e.GET("/test", GinHandler())
	})

	// Run API server

	go func() {
		defer wg.Done()
		ginRuntime.Run(context.Background())
	}()

	wg.Wait()

}

func CronHandler() client.CronHandler {
	return func(c context.Context) error {
		s, ctx := observability.NewSpan(c, "an_span")
		defer s.Finish()

		l := observability.InitLogger("local", "test_service", "1.0.0", true)
		l.Error(ctx).Str("t", "user").Msg("a cron message")
		return nil
	}
}

func GinHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		s, ctx := observability.NewSpan(c, "service.repository.action")
		defer s.Finish()

		l := observability.InitLogger("local", "test_service", "1.0.0", true)
		l.Debug(ctx).Msg("a gin message")
	}
}
