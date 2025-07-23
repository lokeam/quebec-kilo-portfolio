package middleware

import (
	"log"
	"net/http"

	"github.com/lokeam/qko-beta/internal/services"
	"github.com/lokeam/qko-beta/internal/shared/httputils"
)

// RequireUserExists ensures the user exists in the database
//
// WHAT THIS DOES:
// 1. Gets userID from the request (set by Auth0 middleware)
// 2. Checks if user exists in database
// 3. If user doesn't exist, creates a minimal user record
// 4. Always continues to the next handler (endpoint)
//
// WHY THIS EXISTS:
// - Auth0 handles user authentication
// - This middleware handles user MANAGEMENT (ensuring user exists in database)
// - Endpoints handle BUSINESS logic (spend tracking, locations, etc.)
//
// WHEN THIS RUNS:
// - Runs BEFORE all endpoints
// - Runs for ALL protected routes (/api/v1/*)
// - Runs once per request
//
// WHAT HAPPENS IF IT FAILS:
// - Logs the error but doesn't crash
// - Continues to endpoint anyway
// - Endpoint might fail with foreign key errors (but won't crash app)
func RequireUserExists(userService services.UserService) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// SAFETY CHECK 1: If userService is nil (failed to create), just continue
			// This prevents the app from crashing if userService creation failed
			if userService == nil {
				log.Printf("Warning: userService is nil, skipping user creation")
				next.ServeHTTP(w, r)
				return
			}

			// STEP 1: Get userID from request context (set by Auth0 middleware)
			userID := httputils.GetUserID(r)

			// SAFETY CHECK 2: If userID is empty, just continue
			// This prevents errors if Auth0 didn't set the userID properly
			if userID == "" {
				log.Printf("Warning: userID is empty, skipping user creation")
				next.ServeHTTP(w, r)
				return
			}

			// STEP 2: Check if user exists in our database
			exists, err := userService.UserExists(r.Context(), userID)
			if err != nil {
				// SAFETY CHECK 3: If database check fails, log and continue
				// This prevents the app from crashing if database is down
				log.Printf("Warning: failed to check user existence: %v", err)
				next.ServeHTTP(w, r)
				return
			}

			// STEP 3: If user doesn't exist, create them
			if !exists {
				err = userService.CreateUserFromID(r.Context(), userID)
				if err != nil {
					// SAFETY CHECK 4: If user creation fails, log but continue
					// This prevents the app from crashing if user creation fails
					log.Printf("Warning: failed to create user: %v", err)
					// Continue anyway - don't fail the request
				}
			}

			// STEP 4: Always continue to the next handler (your endpoint)
			// Whether user creation succeeded or failed, let your endpoint try
			next.ServeHTTP(w, r)
		})
	}
}