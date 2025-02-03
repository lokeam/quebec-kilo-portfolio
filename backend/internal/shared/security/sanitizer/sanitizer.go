package security

import (
	"fmt"
	"html"
	"regexp"
	"sync"

	"github.com/microcosm-cc/bluemonday"
)

type SanitizationError struct {
	// So know where the error came from
	Message string
}

func (e *SanitizationError) Error() string {
	return fmt.Sprintf("sanitization error: %s", e.Message)
}

// Types / Interfaces
type Sanitizer struct {
	// Whitelist of HTML elements + attributes
	policy         *bluemonday.Policy

	// Regex method of sanitizing HTML content to prevent XSS attacks
	safetyPattern  *regexp.Regexp

	// Mutex prevents race conditions since multiple requests can be made at once
	mu              sync.RWMutex
}

// Constructor
func NewSanitizer() *Sanitizer {
	return &Sanitizer{
		policy:         bluemonday.UGCPolicy(),
		safetyPattern:  regexp.MustCompile(`^[\w\s\-\.,?!]+$`),
	}
}

// Methods
func (s *Sanitizer) SanitizeSearchQuery(input string) (string, error) {
	// Lock the mutex to prevent race conditions
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Escape HTML characters
	escaped := html.EscapeString(input)

	// Sanitize content
	sanitized := s.policy.Sanitize(escaped)

	// Validate that the sanitized pattern is safe
	if !s.safetyPattern.MatchString(sanitized) {
		return "", &SanitizationError{
			Message: "input contains invalid characters",
		}
	}

	return sanitized, nil
}
