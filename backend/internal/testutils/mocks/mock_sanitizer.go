// internal/testutils/mocks/mock_sanitizer.go
package mocks

import (
	"errors"

	"github.com/lokeam/qko-beta/internal/interfaces"
)

// SimpleMockSanitizer is a mock for interfaces.Sanitizer.
   type SimpleMockSanitizer struct {
       // Impl is the overridable function for sanitization.
       Impl func(query string) (string, error)
   }

   // SanitizeSearchQuery calls the overridden implementation if set.
   func (ms *SimpleMockSanitizer) SanitizeSearchQuery(query string) (string, error) {
       if ms.Impl != nil {
           return ms.Impl(query)
       }
       if query == "<script>alert('trigger xss sanitizer error');</script>" {
           return "", errors.New("sanitizer failure")
       }
       return query, nil
   }

   // Ensure SimpleMockSanitizer implements interfaces.Sanitizer:
   var _ interfaces.Sanitizer = (*SimpleMockSanitizer)(nil)
