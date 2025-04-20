package constants

// ContextKey is a custom type for context keys to avoid collisions
type ContextKey string

// Common context keys
const (
	UserIDKey ContextKey = "userID"
	RequestIDKey ContextKey = "requestID"
)