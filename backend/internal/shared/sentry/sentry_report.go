package sentry

import (
	"net/http"

	sentrygo "github.com/getsentry/sentry-go"
	"github.com/lokeam/qko-beta/internal/interfaces"
)

// Report Error logs an error through the logger and reports it to Sentry
func ReportError(
	logger interfaces.Logger,
	msg string,
	err error,
	request *http.Request,
	fields map[string]any,
) {
	// Log through normal channels
	if fields == nil {
		fields = make(map[string]any)
	}
	fields["error"] = err.Error()
	logger.Error(msg, fields)

	// Also send to Sentry with request context if available
	if request != nil {
		if hub := sentrygo.GetHubFromContext(request.Context()); hub != nil {
			hub.WithScope(func(scope *sentrygo.Scope) {
				for k, v := range fields {
					scope.SetExtra(k, v)
				}
				hub.CaptureException(err)
			})
		} else {
			sentrygo.CaptureException(err)
		}
	} else {
		sentrygo.CaptureException(err)
	}
}

// ReportNonHTTContextError
func ReportNonHTTPContextError(
	logger interfaces.Logger,
	msg string,
	err error,
	fields map[string]any,
) {
	ReportError(logger, msg, err, nil, fields)
}
