package worker

import (
	"context"
	"database/sql"
	"time"

	"go-rest-boilerplate/config"
	"go-rest-boilerplate/data"
)

const (
	CreateJobsTableSQL = `
		CREATE TABLE IF NOT EXISTS jobs (
			id      BIGSERIAL NOT NULL PRIMARY KEY,
			name    VARCHAR NOT NULL,
			status  VARCHAR NOT NULL DEFAULT 'pending',
			payload VARCHAR NOT NULL DEFAULT '',
			error   VARCHAR NOT NULL DEFAULT '',
			
			created_at   TIMESTAMPTZ NOT NULL DEFAULT (CURRENT_TIMESTAMP),
			scheduled_at TIMESTAMPTZ NOT NULL DEFAULT (CURRENT_TIMESTAMP),
			updated_at   TIMESTAMPTZ NOT NULL DEFAULT (CURRENT_TIMESTAMP)
		);`

	ScheduleJobSQL = `
		INSERT INTO jobs (name, status, payload, scheduled_at)
		VALUES ($1, $2, $3, $4);`
	FindActiveJobSQL = `
		SELECT * FROM jobs
		WHERE status IN ('pending','running') AND name = $1;`
	FindScheduledJobSQL = `
		SELECT * FROM jobs
		WHERE status = 'pending' AND name = $1 AND scheduled_at = $2
		ORDER BY id DESC
		LIMIT 1;`

	// we could use FOR UPDATE locks but this means the "processing" status
	// update won't be visible to other workers and we lose visibility into
	// jobs that are actually under processing.
	//
	// NEWID() is used to randomize record selection to minimize two workers
	// picking up the same job when there's high load.
	FindPendingJobSQL = `
		SELECT * FROM jobs
		WHERE
			status = 'pending' AND
			(scheduled_at IS NULL OR
				scheduled_at < CURRENT_TIMESTAMP)
		ORDER BY RANDOM()
		LIMIT 1;`

	UpdateJobStatusSQL = `
		UPDATE jobs
		SET status = $1,
			error = $2,
			updated_at = $3
		WHERE id = $4 AND status = $5;`
)

type JobStatus string

const (
	// Job is awaiting to be picked up by a worker
	PendingStatus JobStatus = "pending"
	// Job has been picked up by a worker and is currently running
	RunningStatus JobStatus = "running"
	// Job has been ran by a worker and failed
	FailedStatus JobStatus = "failed"
	// Job has been ran by a worker and completed
	CompletedStatus JobStatus = "completed"
)

type Job struct {
	ID      int64     `db:"id" json:"id"`
	Name    string    `db:"name" json:"name"`
	Status  JobStatus `db:"status" json:"status"`
	Payload string    `db:"payload" json:"payload"`
	Error   string    `db:"error" json:"error"`

	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	ScheduledAt time.Time `db:"scheduled_at" json:"scheduled_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

func createJobsTable(ctx context.Context) error {
	_, err := data.Run(ctx, func(s data.Scope) (interface{}, error) {
		return nil, s.Exec(CreateJobsTableSQL)
	})
	return err
}

func scheduleJob(ctx context.Context, name string, payload []byte, t time.Time) (*Job, error) {
	if t.IsZero() {
		t = time.Now()
	}

	return data.Run(ctx, func(s data.Scope) (*Job, error) {
		config.FromContext(ctx).Println("Scheduling", name)
		config.FromContext(ctx).Println("Payload", string(payload))
		config.FromContext(ctx).Println("Time", t)

		if err := s.Exec(ScheduleJobSQL,
			name, PendingStatus, string(payload), t,
		); err != nil {
			return nil, err
		}

		job := &Job{}
		if err := s.Get(job, FindScheduledJobSQL, name, t); err != nil {
			return nil, err
		} else {
			return job, nil
		}
	})
}

func scheduleJobIfNotExists(ctx context.Context, name string, payload []byte, t time.Time) (*Job, error) {
	if t.IsZero() {
		t = time.Now()
	}

	return data.Run(ctx, func(s data.Scope) (*Job, error) {
		config.FromContext(ctx).Println("Scheduling", name)
		config.FromContext(ctx).Println("Payload", string(payload))
		config.FromContext(ctx).Println("Time", t)

		activeJobs := &[]Job{}
		if err := s.Select(activeJobs, FindActiveJobSQL, name); err != nil && err != sql.ErrNoRows {
			return nil, err
		}
		// job has already existed
		if len(*activeJobs) > 0 {
			return nil, nil
		}

		if err := s.Exec(ScheduleJobSQL,
			name, PendingStatus, string(payload), t,
		); err != nil {
			return nil, err
		}

		job := &Job{}
		if err := s.Get(job, FindScheduledJobSQL, name, t); err != nil {
			return nil, err
		} else {
			return job, nil
		}
	})
}

func takeOnePendingJob(ctx context.Context) (*Job, error) {
	return data.Run(ctx, func(s data.Scope) (*Job, error) {
		job := &Job{}
		if err := s.Get(job, FindPendingJobSQL); err != nil {
			if data.IsNoRows(err) {
				return nil, nil
			} else {
				return nil, err
			}
		}
		if err := s.Exec(UpdateJobStatusSQL,
			RunningStatus, "", time.Now(),
			job.ID, job.Status,
		); err != nil {
			return nil, err
		}
		return job, nil
	})
}

func markJobAsFailed(ctx context.Context, jobId int64, reason string) error {
	return data.Exec(ctx, UpdateJobStatusSQL,
		FailedStatus, reason, time.Now(),
		jobId, RunningStatus)
}

func markJobAsCompleted(ctx context.Context, jobId int64) error {
	return data.Exec(ctx, UpdateJobStatusSQL,
		CompletedStatus, "", time.Now(),
		jobId, RunningStatus)
}
