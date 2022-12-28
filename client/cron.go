package client

import (
	"context"
	"fmt"
	"sync"
	"time"

	bloockContext "github.com/bloock/go-kit/context"
	"github.com/bloock/go-kit/observability"
	"github.com/go-co-op/gocron"
	"github.com/google/uuid"
)

type CronHandler func(context.Context) error
type cronJob struct {
	ctx     context.Context
	name    string
	spec    time.Duration
	fixTime string
	job     CronHandler
	l       observability.Logger
}

func newCronJob(name string, spec time.Duration, fixTime string, job CronHandler, l observability.Logger) cronJob {
	return cronJob{
		name:    name,
		spec:    spec,
		fixTime: fixTime,
		job:     job,
		l:       l,
	}
}

func (c *cronJob) WithContext(ctx context.Context) cronJob {
	return cronJob{
		ctx:     ctx,
		name:    c.name,
		fixTime: c.fixTime,
		spec:    c.spec,
		job:     c.job,
		l:       c.l,
	}
}

func (c cronJob) Run() {

	ctx := context.WithValue(c.ctx, bloockContext.UserIDKey, "")
	ctx = context.WithValue(ctx, bloockContext.RequestIDKey, uuid.New().String())

	s, ctx := observability.NewSpan(ctx, c.name)
	defer s.Finish()

	c.l.Info(ctx).Str("job-name", c.name).Msg("starting job")

	err := c.job(ctx)
	if err != nil {
		c.l.Error(ctx).Str("job-name", c.name).Msgf("error running cron: %s", err.Error())
		return
	}
	c.l.Info(ctx).Str("job-name", c.name).Msg("job runned successfully")
}

type CronClient struct {
	ctx       context.Context
	scheduler *gocron.Scheduler

	handlers []cronJob

	l  observability.Logger
	wg *sync.WaitGroup
}

func NewCronClient(ctx context.Context, l observability.Logger) (*CronClient, error) {
	l.UpdateLogger(l.With().Str("layer", "infrastructure").Str("component", "cron").Logger())

	c := gocron.NewScheduler(time.UTC)

	client := CronClient{
		ctx:       ctx,
		scheduler: c,
		handlers:  make([]cronJob, 0),
		l:         l,
		wg:        &sync.WaitGroup{},
	}

	return &client, nil

}

func (a *CronClient) AddJob(name string, spec time.Duration, fixTime string, handler CronHandler) {
	job := newCronJob(name, spec, fixTime, handler, a.l)
	a.handlers = append(a.handlers, job)
}

func (a *CronClient) Start(ctx context.Context) error {
	for _, handler := range a.handlers {
		if handler.fixTime != "" {
			_, err := a.scheduler.Cron(handler.fixTime).Do(handler.WithContext(ctx).Run)
			if err != nil {
				return err
			}
			continue
		}
		_, err := a.scheduler.Every(handler.spec).Do(handler.WithContext(ctx).Run)
		if err != nil {
			return err
		}
	}

	a.scheduler.StartAsync()

	return nil
}

func (a *CronClient) Close(shutdownTime time.Duration) error {
	stop := make(chan bool)
	go func() {
		a.scheduler.Stop()
		stop <- true
	}()
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTime)
	defer cancel()

	select {
	case <-stop:
		return nil
	case <-ctx.Done():
		return fmt.Errorf("couldn't close cron client before timeout")
	}
}
