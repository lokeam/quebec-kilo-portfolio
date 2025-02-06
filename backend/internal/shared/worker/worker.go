package worker

import (
	"context"
	"time"

	"github.com/lokeam/qko-beta/internal/shared/logger"
)

// Worker is a generic background worker that executes a job function at a regular interval.
type Worker struct {
	interval   time.Duration
	job        func(context.Context) error
	condition  func() bool
	logger     logger.Logger
}

// NewWorker creates a new Worker with the specified interval and job function.
func NewWorker(
	interval time.Duration,
	job func(context.Context) error,
	condition  func() bool,
	logger logger.Logger,
) *Worker {
	return &Worker{
		interval:  interval,
		job:       job,
		condition: condition,
		logger:    logger,
	}
}

// Start launches the worker in a new goroutine.
// It runs the job at the specified interval until the context is cancelled.
func (w *Worker) Start(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()
	w.logger.Info("Worker started", nil)
	for {
		select {
		case <-ctx.Done():
			w.logger.Info("Worker stopped", nil)
			return
		case <-ticker.C:
			// Check if condition is met before executing the job.
			if w.condition == nil || w.condition() {
				if err := w.job(ctx); err != nil {
					w.logger.Error("Worker job failed", map[string]any{"error": err.Error()})
				}
			} else {
				w.logger.Debug("Worker condition not met; skipping job", nil)
			}
		}
	}
}
