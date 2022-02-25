package shutdown

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
)

func NewGracefullShutdown(l zerolog.Logger) context.Context {
	ctx, cancelFunc := context.WithCancel(context.Background())

	l = l.With().Str("layer", "infrastructure").Str("component", "shutdown").Logger()
	sigint := make(chan os.Signal, 1)

	// interrupt signal sent from terminal
	signal.Notify(sigint, os.Interrupt)
	// sigterm signal sent from kubernetes
	signal.Notify(sigint, syscall.SIGTERM)

	go func() {
		sig := <-sigint

		l.Info().Timestamp().Str("signal", sig.String()).Msg("received shutdown signal")

		cancelFunc()
	}()

	return ctx
}
