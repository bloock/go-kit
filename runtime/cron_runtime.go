package runtime

import (
	"context"
	"time"

	"github.com/bloock/go-kit/client"
	"github.com/rs/zerolog"
)

type CronRuntime struct {
	client *client.CronClient

	shutdownTime time.Duration
	logger       zerolog.Logger
}

func NewCronRuntime(c *client.CronClient, shutdownTime time.Duration, l zerolog.Logger) (*CronRuntime, error) {
	e := CronRuntime{
		client:       c,
		shutdownTime: shutdownTime,
		logger:       l,
	}

	return &e, nil
}

func (e *CronRuntime) AddHandler(name string, spec time.Duration, h client.CronHandler) {
	e.client.AddJob(name, spec, h)
}

func (e *CronRuntime) Run(ctx context.Context) {
out:
	for {
		err := e.client.Start(ctx)
		if err != nil {
			e.logger.Info().Msgf("error while starting cron worker: %s", err.Error())
			break out
		}

		e.logger.Info().Msg("cron runtime started successfully")

		select {
		case <-ctx.Done():
			break out
		}
	}

	if err := e.client.Close(e.shutdownTime); err != nil {
		e.logger.Info().Msgf("error while closing cron runtime: %s", err.Error())
	} else {
		e.logger.Info().Msg("cron runtime closed successfully")
	}
}
