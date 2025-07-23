package middleware

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/lokeam/qko-beta/internal/shared/constants"
)

const (
	Auth0Domain = "AUTH0_DOMAIN"
	Auth0Audience = "AUTH0_AUDIENCE"

	ContentType = "Content-Type"
	ApplicationJSON = "application/json"
)

// CustomClaims contains custom data we want from the token.
type CustomClaims struct {
	Scope string `json:"scope"`
}

// NOTE: Validate does nothing for this example, but we need
// it to satisfy validator.CustomClaims interface.
func (c CustomClaims) Validate(ctx context.Context) error {
	return nil
}


// Helper function to create and configure the JWT validator
func createJWTValidator() (*jwtmiddleware.JWTMiddleware, error) {
	issuerURL, err := url.Parse("https://" + os.Getenv(Auth0Domain) + "/")
	if err != nil {
		return nil, err
	}

	provider := jwks.NewCachingProvider(issuerURL, 5 * time.Minute)

	jwtValidator, err := validator.New(
		provider.KeyFunc,
		validator.RS256,
		issuerURL.String(),
		[]string{os.Getenv(Auth0Audience)},
		validator.WithCustomClaims(
			func() validator.CustomClaims {
				return &CustomClaims{}
			},
		),
		validator.WithAllowedClockSkew(time.Minute),
	)
	if err != nil {
		return nil, err
	}

	errorHandler := func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("JWT validation error: %v", err)
		log.Printf("Request path: %s", r.URL.Path)

		w.Header().Set(ContentType, ApplicationJSON)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"message":"Failed to validate JWT."}`))
	}

	middleware := jwtmiddleware.New(
		jwtValidator.ValidateToken,
		jwtmiddleware.WithErrorHandler(errorHandler),
	)

	return middleware, nil
}


// Helper function to extract UserID from validated JWT token and add it to request context
func extractUserIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("extractUserIDMiddleware: Processing request for %s", r.URL.Path)

		// Get JWT token from the request
		if token := r.Context().Value(jwtmiddleware.ContextKey{}); token != nil {
			if claims, ok := token.(*validator.ValidatedClaims); ok {
				// Extract user ID from JWT token
				userID := claims.RegisteredClaims.Subject
				log.Printf("Extracted userID from JWT: %s", userID)

				// Add userID to request context
				ctx := context.WithValue(r.Context(), constants.UserIDKey, userID)
				r = r.WithContext(ctx)
				log.Printf("Successfully added userID to context: %s", userID)
			} else {
				log.Printf("Failed to cast token to ValidatedClaims")
			}
		} else {
			log.Printf("No token found in context")
		}

		// Pass the updated request with context to the next handler
		next.ServeHTTP(w, r)
	})
}

// EnsureValidToken is a middleware that will check the validity of our JWT.
func EnsureValidToken() func(next http.Handler) http.Handler {
	jwtMiddleware, err := createJWTValidator()
	if err != nil {
		log.Fatalf("Failed to create JWT validator: %v", err)
	}

	return func(next http.Handler) http.Handler {
		log.Printf("Auth0 middleware: Creating middleware chain")
		// Compose middleware: JWT validation -> extract userID -> Next handler
		return jwtMiddleware.CheckJWT(extractUserIDMiddleware(next))
	}
}