package worker

import (
	"context"
	"time"

	"github.com/lokeam/qko-beta/internal/shared/logger"
)

// Worker is a generic background worker that executes a job function at regular intervals.
type Worker struct {
	interval  time.Duration
	job       func(context.Context) error   // The job function to execute
	condition func() bool                   // Optional condition: if provided, job runs only if condition() returns true.
	logger    logger.Logger
}


func NewWorker(interval time.Duration, job func(context.Context) error, condition func() bool, logger logger.Logger) *Worker {
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

	w.logger.Info("Worker started", nil)
	for {
		select {
		case <-ctx.Done():
			w.logger.Info("Worker stopping: context cancelled", nil)
			return
		case <-ticker.C:
			if w.condition != nil && !w.condition() {
				w.logger.Debug("Worker condition not met; skipping job", nil)
				continue
			}
			if err := w.job(ctx); err != nil {
				w.logger.Error("Worker job error", map[string]any{"error": err.Error()})
			}
		}
	}
}
