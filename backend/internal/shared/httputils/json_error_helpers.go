package httputils

import (
	"encoding/json"

	"github.com/lokeam/qko-beta/internal/interfaces"
)

/*
	LogJSONError provides detailed logging for JSON parsing errors,
	including context around syntax errors to make debugging JSON less hellish.
*/

func LogJSONError(
	logger interfaces.Logger,
	requestID string,
	err error,
	rawJSON []byte,
) {
	// Log basic error info
	errorContext := map[string]any{
		"error":       err.Error(),
		"errorType":   err,
		"requestID":   requestID,
	}

	// For syntax errors, add context to display where error occured
	if syntaxErr, ok := err.(*json.SyntaxError); ok {
		errorContext["errorType"] = "JSON syntax error"
		errorContext["offset"] = syntaxErr.Offset

		// Pull out a snippet around the error location
		start := max(0, int64(syntaxErr.Offset)-20)
		end := min(int64(len(rawJSON)), int64(syntaxErr.Offset)+20)
		errorContext["errorContext"] = string(rawJSON[start:end])

		// For unmarshaling errors, add type info
		if _, ok := err.(*json.UnmarshalTypeError); ok {
			errorContext["errorType"] = "JSON type mismatch"
		}

		logger.Error("JSON parsing error", errorContext)
	}
}


// Helper functions
func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func min(a, b int64) int64 {
	if a > b {
		return b
	}
	return a
}