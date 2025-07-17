package email

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/lokeam/qko-beta/internal/appcontext"
)

// EmailJob represents an email job to be processed
type EmailJob struct {
	ID        string
	Type      EmailJobType
	UserID    string
	Email     string
	Data      map[string]interface{}
	CreatedAt time.Time
	Retries   int
}

// EmailJobType defines the type of email job
type EmailJobType string

const (
	EmailJobTypeDeletionRequest      EmailJobType = "deletion_request"
	EmailJobTypeGracePeriodReminder  EmailJobType = "grace_period_reminder"
	EmailJobTypeDeletionConfirmation EmailJobType = "deletion_confirmation"
	EmailJobTypeDataExport           EmailJobType = "data_export"
	EmailJobTypeWelcomeBack          EmailJobType = "welcome_back"
)

// EmailQueue handles asynchronous email processing
type EmailQueue struct {
	appCtx       *appcontext.AppContext
	emailService EmailService
	jobs         chan EmailJob
	workers      int
	maxRetries   int
	stopChan     chan struct{}
	wg           sync.WaitGroup
	mu           sync.RWMutex
	running      bool
}

// NewEmailQueue creates a new email queue
func NewEmailQueue(appCtx *appcontext.AppContext, emailService EmailService, workers, maxRetries int) *EmailQueue {
	if workers <= 0 {
		workers = 3 // Default number of workers
	}
	if maxRetries <= 0 {
		maxRetries = 3 // Default max retries
	}

	return &EmailQueue{
		appCtx:       appCtx,
		emailService: emailService,
		jobs:         make(chan EmailJob, 100), // Buffer for 100 jobs
		workers:      workers,
		maxRetries:   maxRetries,
		stopChan:     make(chan struct{}),
	}
}

// Start starts the email queue workers
func (eq *EmailQueue) Start() error {
	eq.mu.Lock()
	defer eq.mu.Unlock()

	if eq.running {
		return fmt.Errorf("email queue is already running")
	}

	eq.running = true
	eq.appCtx.Logger.Info("Starting email queue", map[string]any{
		"workers": eq.workers,
	})

	// Start workers
	for i := 0; i < eq.workers; i++ {
		eq.wg.Add(1)
		go eq.worker(i)
	}

	return nil
}

// Stop stops the email queue
func (eq *EmailQueue) Stop() error {
	eq.mu.Lock()
	defer eq.mu.Unlock()

	if !eq.running {
		return fmt.Errorf("email queue is not running")
	}

	eq.appCtx.Logger.Info("Stopping email queue", map[string]any{})
	eq.running = false
	close(eq.stopChan)
	close(eq.jobs)

	// Wait for all workers to finish
	eq.wg.Wait()

	return nil
}

// EnqueueJob adds a job to the queue
func (eq *EmailQueue) EnqueueJob(ctx context.Context, jobType EmailJobType, userID, email string, data map[string]interface{}) error {
	eq.mu.RLock()
	defer eq.mu.RUnlock()

	if !eq.running {
		return fmt.Errorf("email queue is not running")
	}

	job := EmailJob{
		ID:        fmt.Sprintf("%s_%s_%d", jobType, userID, time.Now().Unix()),
		Type:      jobType,
		UserID:    userID,
		Email:     email,
		Data:      data,
		CreatedAt: time.Now(),
		Retries:   0,
	}

	select {
	case eq.jobs <- job:
		eq.appCtx.Logger.Info("Email job enqueued", map[string]any{
			"jobID":   job.ID,
			"type":    jobType,
			"userID":  userID,
			"email":   email,
		})
		return nil
	case <-ctx.Done():
		return fmt.Errorf("context cancelled while enqueuing job")
	default:
		return fmt.Errorf("email queue is full")
	}
}

// worker processes email jobs
func (eq *EmailQueue) worker(id int) {
	defer eq.wg.Done()

	eq.appCtx.Logger.Info("Email worker started", map[string]any{
		"workerID": id,
	})

	for {
		select {
		case job, ok := <-eq.jobs:
			if !ok {
				eq.appCtx.Logger.Info("Email worker stopping", map[string]any{
					"workerID": id,
				})
				return
			}

			eq.processJob(job)
		case <-eq.stopChan:
			eq.appCtx.Logger.Info("Email worker stopping", map[string]any{
				"workerID": id,
			})
			return
		}
	}
}

// processJob processes a single email job
func (eq *EmailQueue) processJob(job EmailJob) {
	ctx := context.Background()

	eq.appCtx.Logger.Info("Processing email job", map[string]any{
		"jobID":   job.ID,
		"type":    job.Type,
		"userID":  job.UserID,
		"email":   job.Email,
		"retries": job.Retries,
	})

	var err error

	switch job.Type {
	case EmailJobTypeDeletionRequest:
		gracePeriodEnd, ok := job.Data["gracePeriodEnd"].(time.Time)
		if !ok {
			err = fmt.Errorf("invalid gracePeriodEnd data")
			break
		}
		err = eq.emailService.SendDeletionRequestEmail(ctx, job.UserID, job.Email, job.Data["userName"].(string), gracePeriodEnd)

	case EmailJobTypeGracePeriodReminder:
		daysLeft, ok := job.Data["daysLeft"].(int)
		if !ok {
			err = fmt.Errorf("invalid daysLeft data")
			break
		}
		deletionDate, ok := job.Data["deletionDate"].(time.Time)
		if !ok {
			err = fmt.Errorf("invalid deletionDate data")
			break
		}
		userName, ok := job.Data["userName"].(string)
		if !ok {
			err = fmt.Errorf("invalid userName data")
			break
		}

		err = eq.emailService.SendGracePeriodReminderEmail(ctx, job.UserID, job.Email, userName, daysLeft, deletionDate)

	case EmailJobTypeDeletionConfirmation:
		deletionDate, ok := job.Data["deletionDate"].(time.Time)
		if !ok {
			err = fmt.Errorf("invalid deletionDate data")
			break
		}
		userName, ok := job.Data["userName"].(string)
		err = eq.emailService.SendDeletionConfirmationEmail(ctx, job.UserID, job.Email, userName, deletionDate)

	case EmailJobTypeDataExport:
		exportURL, ok := job.Data["exportURL"].(string)
		if !ok {
			err = fmt.Errorf("invalid exportURL data")
			break
		}
		userName, ok := job.Data["userName"].(string)
		if !ok {
			err = fmt.Errorf("invalid userName data")
			break
		}
		err = eq.emailService.SendDataExportEmail(ctx, job.UserID, job.Email, userName, exportURL)

	case EmailJobTypeWelcomeBack:
		userName, ok := job.Data["userName"].(string)
		if !ok {
			err = fmt.Errorf("invalid userName data")
			break
		}
		err = eq.emailService.SendWelcomeBackEmail(ctx, job.UserID, job.Email, userName)

	default:
		err = fmt.Errorf("unknown email job type: %s", job.Type)
	}

	if err != nil {
		eq.appCtx.Logger.Error("Failed to process email job", map[string]any{
			"jobID":   job.ID,
			"type":    job.Type,
			"userID":  job.UserID,
			"email":   job.Email,
			"error":   err.Error(),
			"retries": job.Retries,
		})

		// Retry if under max retries
		if job.Retries < eq.maxRetries {
			job.Retries++
			// Exponential backoff: 1s, 2s, 4s
			backoff := time.Duration(1<<uint(job.Retries-1)) * time.Second

			eq.appCtx.Logger.Info("Retrying email job", map[string]any{
				"jobID":   job.ID,
				"retries": job.Retries,
				"backoff": backoff.String(),
			})

			// Schedule retry
			go func() {
				time.Sleep(backoff)
				select {
				case eq.jobs <- job:
				default:
					eq.appCtx.Logger.Error("Failed to retry email job - queue full", map[string]any{
						"jobID": job.ID,
					})
				}
			}()
		} else {
			eq.appCtx.Logger.Error("Email job failed permanently", map[string]any{
				"jobID":   job.ID,
				"type":    job.Type,
				"userID":  job.UserID,
				"email":   job.Email,
				"retries": job.Retries,
			})
		}
	} else {
		eq.appCtx.Logger.Info("Email job completed successfully", map[string]any{
			"jobID":  job.ID,
			"type":   job.Type,
			"userID": job.UserID,
			"email":  job.Email,
		})
	}
}

// IsRunning returns true if the queue is running
func (eq *EmailQueue) IsRunning() bool {
	eq.mu.RLock()
	defer eq.mu.RUnlock()
	return eq.running
}