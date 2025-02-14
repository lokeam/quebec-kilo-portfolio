package worker

import (
	"context"
	"time"

	"github.com/lokeam/qko-beta/internal/interfaces"
)

const (
	jobStart = "Worker started"
	jobStopContextCancelled = "Worker stopping: context cancelled"
	jobSkipped              = "Worker condition not met; skipping job"
	jobError                = "Worker job error"
)

// Worker is a generic background worker that executes a job function at regular intervals.
type Worker struct {
	interval  time.Duration
	job       func(context.Context) error   // The job function to execute
	condition func() bool                   // Optional condition: if provided, job runs only if condition() returns true.
	logger    interfaces.Logger
}

func NewWorker(
	interval time.Duration,
	job func(context.Context) error,
	condition func() bool,
	logger interfaces.Logger) *Worker {
	return &Worker{
		interval:  interval,
		job:       job,
		condition: condition,
		logger:    logger,
	}
}

// Start runs the worker until the context is cancelled.
func (w *Worker) Start(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	w.logger.Info(jobStart, nil)
	for {
		select {
		case <-ctx.Done():
			w.logger.Info(jobStopContextCancelled, nil)
			return
		case <-ticker.C:
			if w.condition != nil && !w.condition() {
				w.logger.Debug(jobSkipped, nil)
				continue
			}
			if err := w.job(ctx); err != nil {
				w.logger.Error(jobError, map[string]any{"error": err.Error()})
			}
		}
	}
}
