package security

import (
	"fmt"
	"html"
	"regexp"
	"sync"

	"github.com/lokeam/qko-beta/internal/interfaces"
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
func NewSanitizer() (*Sanitizer, error) {
	return &Sanitizer{
		policy:         bluemonday.UGCPolicy(),
		safetyPattern:  regexp.MustCompile(`^[\w\s\-\.,?!]+$`),
	}, nil
}

// Make sure that Sanitizer implements interfaces.Sanitizer.
var _ interfaces.Sanitizer = (*Sanitizer)(nil)

// Methods
func (s *Sanitizer) SanitizeSearchQuery(input string) (string, error) {
	sanitizedContent, err := s.performBaseSanitization(input)
	if err != nil {
		return "", err
	}

	// Perform additional sanitization for search queries
	if !s.safetyPattern.MatchString(sanitizedContent) {
		return "", &SanitizationError{
			Message: "input contains invalid characters",
		}
	}

	return sanitizedContent, nil
}

func (s *Sanitizer) SanitizeString(input string) (string, error) {
	return s.performBaseSanitization(input)
}

func (s *Sanitizer) performBaseSanitization(input string) (string, error) {
	// Lock mutex to prevent race conditions
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Escape HTML characters
	escaped := html.EscapeString(input)

	// Sanitize string content
	sanitizedContent := s.policy.Sanitize(escaped)

	return sanitizedContent, nil
}
