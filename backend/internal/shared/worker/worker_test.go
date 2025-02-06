package worker

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	jobStopContextCancelled = "Worker stopping: context cancelled"
	jobSkipped              = "Worker condition not met; skipping job"
	jobError                = "Worker job error"
)


type TestLogger struct {
	mu            sync.Mutex
	infoCalls     []string
	debugCalls    []string
	errorCalls    []string
}

func NewTestLogger() *TestLogger {
	return &TestLogger{}
}

func (tl *TestLogger) Info(msg string, fields map[string]any) {
	tl.mu.Lock()
	defer tl.mu.Unlock()

	tl.infoCalls = append(tl.infoCalls, msg)
}

func (tl *TestLogger) Debug(msg string, fields map[string]any) {
	tl.mu.Lock()
	defer tl.mu.Unlock()
	tl.debugCalls = append(tl.debugCalls, msg)
}

func (tl *TestLogger) Error(msg string, fields map[string]any) {
	tl.mu.Lock()
	defer tl.mu.Unlock()

	errStr := msg
	if fields != nil {
		if err, ok := fields["error"].(error); ok {
			errStr += ": " + err.Error()
		}
	}
	tl.errorCalls = append(tl.errorCalls, errStr)
}

func (tl *TestLogger) Warn(msg string, fields map[string]any) {
	tl.mu.Lock()
	defer tl.mu.Unlock()
	// You can choose to log warnings similarly to Debug messages.
	tl.debugCalls = append(tl.debugCalls, msg)
}

/*
	Behaviors:
		- Worker periodically runs a job
		- Worker skips job if optional condition isn't met
		- Worker stops when context is cancelled
		- Worker logs errors when job return an error

	Scenarios:
		- Worker starts, runs job a few times, stops when context is cancelled
		- Worker is set with a condition that always returns false. Job never runs.
	  - Worker's job fails and it returns an error. Worker logs error.
*/

// Test Worker scenarios to validate behaviors:
func TestWorkerSuite(t *testing.T) {
	// Set test cases w/ BDD style naming
	testCases := []struct {
		name                      string
		description               string
		interval                  time.Duration          // How often job should run
		jobFactory                func(t *testing.T) (func(ctx context.Context) error, *int) // Factory function to create job and counter
		runCondition              func() bool            // If exists and return false, worker will skip executing the job
		runDurationBeforeCancel   time.Duration          // Duration to let the worker run before cancelling context
		runDurationAfterCancel    time.Duration          // Gives worker time to wrap up and log any final messages after cancelling
		didBehaviorOccur          func(t *testing.T, fakeLogger *TestLogger, jobCount int)
	}{
		{
			name: "Worker stops when context is cancelled",
			description: `
				GIVEN a running worker,
				WHEN the context is cancelled,
				THEN it logs the cancellation message and stops
			`,
			interval: 10 * time.Millisecond,
			jobFactory: func(t *testing.T) (func(ctx context.Context) error, *int) {
				// No operation job, we're not using the counter here
				counter := new(int)
				return func(ctx context.Context) error { return nil }, counter
			},
			runCondition:              nil,
			runDurationBeforeCancel:   30 * time.Millisecond,
			runDurationAfterCancel:    20 * time.Millisecond,
			didBehaviorOccur: func(t *testing.T, testLogger *TestLogger, jobCount int) {
				testLogger.mu.Lock()
				defer testLogger.mu.Unlock()

				var foundJob bool

				for _, msg := range testLogger.infoCalls {
					if msg == jobStopContextCancelled {
						foundJob = true
						break
					}
				}
				assert.True(t, foundJob, "expected worker to log cancellation stop message")
			},
		},
		{
			name: "Worker skips job when condition is false",
			description: `
				GIVEN a worker with a condition that always returns false,
				WHEN the worker starts,
				THEN the job is not executed and a debug log is written
			`,
			interval: 10 * time.Millisecond,
			jobFactory: func(t *testing.T) (func(context.Context) error, *int) {
				counter := new(int)
				return func(ctx context.Context) error {
					(*counter)++
					return nil
				}, counter
			},
			runCondition: func() bool { return false },
			runDurationBeforeCancel: 50 * time.Millisecond,
			runDurationAfterCancel: 0,
			didBehaviorOccur: func(t *testing.T, testLogger *TestLogger, jobCount int) {
				assert.Equal(t, 0, jobCount, "job should not have been executed when condition is false")

				testLogger.mu.Lock()
				defer testLogger.mu.Unlock()

				var foundJob bool
				for _, msg := range testLogger.debugCalls {
					if msg == jobSkipped {
						foundJob = true
						break
					}
				}
				assert.True(t, foundJob, "expected debug log 'Worker condition not met; skipping job' not found")
			},
		},
		{
			name: "Worker logs job error when job returns an error",
			description: `
				GIVEN a worker with a failing job,
				WHEN then job returns an error,
				THEN the worker logs the error
			`,
			interval: 10 * time.Millisecond,
			jobFactory: func(t *testing.T) (func(context.Context) error, *int) {
				counter := new(int)
				jobError := errors.New("job failed")
				return func(ctx context.Context) error {
					(*counter)++

					if *counter == 1 {
						return jobError
					}
					return nil
				}, counter
			},
			runCondition: nil,
			runDurationBeforeCancel:  50 * time.Millisecond,
			runDurationAfterCancel:   0,
			didBehaviorOccur: func(t *testing.T, testLogger *TestLogger, jobCount int) {
				testLogger.mu.Lock()
				defer testLogger.mu.Unlock()

				var found bool

				for _, msg := range testLogger.errorCalls {
					if msg == jobError {
						found = true
						break
					}
				}
				assert.True(t, found, "expected error log for job failure was not triggered")
			},
		},
	}

	// Test runner loop
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// Set up context, test logger and job function
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			testLogger := NewTestLogger()
			job, jobCountPointer := testCase.jobFactory(t)

			// Create worker with specified params
			workerInstance := NewWorker(testCase.interval, job, testCase.runCondition, testLogger)

			// Start worker
			go workerInstance.Start(ctx)

			// Allow worker to run
			time.Sleep(testCase.runDurationBeforeCancel)

			// Stop worker
			cancel()

			// Allow worker to wrap up
			if testCase.runDurationAfterCancel > 0 {
				time.Sleep(testCase.runDurationAfterCancel)
			}

			var jobCount int
			if jobCountPointer != nil {
				jobCount = *jobCountPointer
			}
			testCase.didBehaviorOccur(t, testLogger, jobCount)
		})
	}
}
