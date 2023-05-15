package worker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"go-rest-boilerplate/config"
	"go-rest-boilerplate/data"
	"go-rest-boilerplate/internal"
	"go-rest-boilerplate/worker/runners"

	"github.com/jmoiron/sqlx"
)

const (
	PollingInterval = 5 * time.Second
)

type Worker struct {
	cfg *config.Config
	db  *sqlx.DB
}

func New(cfg *config.Config) *Worker {
	return &Worker{cfg, nil}
}

func ScheduleNow(ctx context.Context, runner runners.Interface) (int64, error) {
	return ScheduleAt(ctx, runner, time.Time{})
}

func ScheduleNowIfNotExists(ctx context.Context, runner runners.Interface) (int64, error) {
	return ScheduleJobIfNotExists(ctx, runner, time.Time{})
}

func ScheduleAt(ctx context.Context, runner runners.Interface, t time.Time) (int64, error) {
	if payload, err := json.Marshal(runner); err != nil {
		return 0, err
	} else if job, err := scheduleJob(ctx, runner.Name(), payload, t); err != nil {
		return 0, err
	} else {
		return job.ID, nil
	}
}

func ScheduleJobIfNotExists(ctx context.Context, runner runners.Interface, t time.Time) (int64, error) {
	payload, err := json.Marshal(runner)
	if err != nil {
		return 0, err
	}

	job, err := scheduleJobIfNotExists(ctx, runner.Name(), payload, t)
	if err != nil {
		return 0, err
	} else if job == nil {
		/* job has already existed */
		return 0, nil
	} else {
		return job.ID, nil
	}
}

func (w *Worker) Start() (err error) {
	defer internal.WrapErr("worker", &err)

	ctx := config.NewContext(context.Background(), w.cfg)
	if ctx, err = w.initializeDB(ctx); err != nil {
		return err
	}

	errChan := make(chan error)
	go w.work(ctx, errChan)
	return <-errChan
}

func (w *Worker) initializeDB(inctx context.Context) (context.Context, error) {
	if db, err := data.Connect(w.cfg); err != nil {
		return inctx, fmt.Errorf("initialize: %w", err)
	} else {
		w.db = db
	}

	ctx := data.NewContext(inctx, w.db)
	return ctx, createJobsTable(ctx)
}

func (w *Worker) work(ctx context.Context, errors chan error) {
	defer close(errors)
	cfg := config.FromContext(ctx)

	for {
		job, err := takeOnePendingJob(ctx)
		if err != nil {
			errors <- err
			break
		} else if job == nil {
			time.Sleep(PollingInterval)
			continue
		}

		cfg.Printf("running %s #%d\n", job.Name, job.ID)
		start := time.Now()

		// we got one "running" job to process
		if err := processJob(ctx, job); err != nil {
			cfg.Printf("failed %s #%d in %s: %s\n",
				job.Name, job.ID,
				time.Since(start).String(), err.Error(),
			)
			if err := markJobAsFailed(ctx, job.ID, err.Error()); err != nil {
				errors <- err
				break
			}
		} else {
			cfg.Printf("completed %s #%d in %s",
				job.Name, job.ID,
				time.Since(start).String(),
			)
			if err := markJobAsCompleted(ctx, job.ID); err != nil {
				errors <- err
				break
			}
		}
	}
}

func processJob(ctx context.Context, job *Job) error {
	runner := runners.New(job.Name)
	if runner == nil {
		return errors.New("unknown (or unregistered) job: " + job.Name)
	}
	if err := json.Unmarshal([]byte(job.Payload), runner); err != nil {
		return fmt.Errorf("malformed payload: %w", err)
	}
	if err := runner.Run(ctx, func(time time.Time) {
		_, _ = ScheduleAt(ctx, runner, time)
	}); err != nil {
		return fmt.Errorf("run failed: %w", err)
	}
	return nil
}
