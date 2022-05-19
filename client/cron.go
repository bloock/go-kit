package client

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog"
)

type CronHandler func(context.Context) error
type cronJob struct {
	ctx  context.Context
	name string
	spec string
	job  CronHandler
	l    zerolog.Logger
}

func newCronJob(name, spec string, job CronHandler, l zerolog.Logger) cronJob {
	return cronJob{
		name: name,
		spec: spec,
		job:  job,
		l:    l,
	}
}

func (c *cronJob) WithContext(ctx context.Context) cronJob {
	return cronJob{
		ctx:  ctx,
		name: c.name,
		spec: c.spec,
		job:  c.job,
		l:    c.l,
	}
}

func (c cronJob) Run() {
	err := c.job(c.ctx)
	if err != nil {
		c.l.Error().Str("name", c.name).Msgf("error running cron job %s: %s", c.name, err.Error())
	}
}

type CronClient struct {
	ctx  context.Context
	cron *cron.Cron

	handlers []cronJob

	l  zerolog.Logger
	wg *sync.WaitGroup
}

func NewCronClient(ctx context.Context, l zerolog.Logger) (*CronClient, error) {
	l = l.With().Str("layer", "infrastructure").Str("component", "cron").Logger()

	c := cron.New()

	client := CronClient{
		cron:     c,
		handlers: make([]cronJob, 0),
		l:        l,
		wg:       &sync.WaitGroup{},
	}

	return &client, nil

}

func (a *CronClient) AddJob(name, spec string, handler CronHandler) {
	job := newCronJob(name, spec, handler, a.l)
	a.handlers = append(a.handlers, job)
}

func (a *CronClient) Start(ctx context.Context) error {
	for _, handler := range a.handlers {
		_, err := a.cron.AddJob(handler.spec, handler.WithContext(ctx))
		if err != nil {
			return err
		}
	}
	a.cron.Start()

	return nil
}

func (a *CronClient) Close(shutdownTime time.Duration) error {
	stopCtx := a.cron.Stop()
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTime)
	defer cancel()

	select {
	case <-stopCtx.Done():
		return nil
	case <-ctx.Done():
		return fmt.Errorf("couldn't close cron client before timeout")
	}
}
