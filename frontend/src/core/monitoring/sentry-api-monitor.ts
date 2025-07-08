import * as Sentry from '@sentry/react';

/**
 * Sentry API Monitoring Service
 *
 * WHAT:
 * When API calls fail in QKO, this code sends the error details to Sentry
 *
 * WHY:
 * With this, we can see exactly which API call failed, why it failed, and which user was affected.
 *
 * HOW:
 * The axios interceptor automatically calls these functions when API calls fail.
 * We can also call these functions manually if we need to log custom errors.
 */

/**
 * Information about an API error that we want to log
 */
export interface ApiErrorContext {
  url?: string;           // Which API endpoint failed (e.g., '/api/users')
  method?: string;        // What type of request (GET, POST, etc.)
  status?: number;        // HTTP status code (404, 500, etc.)
  errorMessage: string;   // What the error actually says
  responseData?: unknown; // Any extra data from the server
}

/**
 * Information about an authentication error (such as token refresh failure)
 */
export interface AuthErrorContext {
  originalUrl?: string;     // What the user was trying to do when auth failed
  originalMethod?: string;  // What type of request they were making
  originalStatus?: number;  // The 401 status that triggered the auth failure
  errorMessage: string;     // Why the auth failed
}

/**
 * Information about an API event we want to track
 */
export interface ApiBreadcrumbContext {
  url?: string;           // Which API endpoint
  method?: string;        // What type of request (GET, POST, etc.)
  message: string;        // What happened (human-readable)
  data?: Record<string, unknown>; // Any extra information
}

/**
 * Log an API error to Sentry so you can see what went wrong
 *
 * WHAT:
 * Takes an error that happened during an API call and sends all the details to Sentry.
 * This includes what the user was trying to do, what went wrong, and any extra context.
 *
 * WHEN:
 * Called automatically by the axios interceptor when any API call fails.
 * We can also call it manually if you want to log custom errors.
 *
 * INFO DISPLAYED IN SENTRY:
 * - The exact API endpoint that failed
 * - The HTTP method (GET, POST, etc.)
 * - The error status code (404, 500, etc.)
 * - The error message
 * - Which user was affected
 * - When it happened
 *
 * EXAMPLE:
 * User tries to load their profile, gets a 500 error.
 * In Sentry we should see: "GET /api/users failed with 500 - Database connection failed"
 */
export function logApiError(error: Error, context: ApiErrorContext): void {
  Sentry.captureException(error, {
    tags: {
      source: 'api',                    // Tells Sentry it's an API error
      endpoint: context.url,            // Which endpoint failed
      method: context.method,           // What type of request
    },
    contexts: {
      api_request: {
        url: context.url,               // The full URL that failed
        method: context.method,         // GET, POST, PUT, DELETE, etc.
        status: context.status,         // HTTP status code (404, 500, etc.)
        error_message: context.errorMessage,  // What the error actually says
        response_data: context.responseData,  // Any extra data from the server
      },
    },
  });
}

/**
 * Log an authentication error (such as when token refresh fails)
 *
 * WHAT:
 * When a user's login token expires and the app can't refresh it, this logs the problem.
 * This is different from regular API errors because it's specifically about authentication.
 *
 * WHEN:
 * Called when the axios interceptor tries to refresh an expired token and fails.
 *
 * INFO DISPLAYED IN SENTRY:
 * - What the user was trying to do when their auth failed
 * - Why the token refresh failed
 * - Which user was affected
 *
 * EXAMPLE:
 * User's session expires, app tries to refresh it, fails.
 * In Sentry we should see: "Auth failed for /api/protected-data - Token refresh failed"
 */
export function logAuthError(error: Error, context: AuthErrorContext): void {
  Sentry.captureException(error, {
    tags: {
      source: 'auth',                   // Tells Sentry it's an auth error
      error_type: 'token_refresh_failed', // What kind of auth error
    },
    contexts: {
      auth_error: {
        original_url: context.originalUrl,     // What the user was trying to do
        original_method: context.originalMethod, // What type of request
        original_status: context.originalStatus, // The 401 that triggered this
        error_message: context.errorMessage,     // Why the auth failed
      },
    },
  });
}

/**
 * Add a breadcrumb to track what happened with an API call
 *
 * WHAT:
 * Adds a note to the timeline showing what happened during an API call.
 * Breadcrumbs create a list of actions that describe the path the user took.
 * Used to understand what happened before an error occurred.
 *
 * WHEN:
 * Use this to track API events that aren't errors but are important to know about.
 * Example: when an API call starts, when it succeeds, or when something noteworthy happens.
 *
 * INFO DISPLAYED IN SENTRY:
 * - What happened (some custom message)
 * - Which API endpoint was involved
 * - What type of request it was
 * - Any extra data you want to include
 * - When it happened
 *
 * EXAMPLE:
 * User starts uploading a file, we want to track it.
 * In Sentry we should see: "File upload started for /api/upload" in the breadcrumb trail
 *
 * DIFFERENCE FROM logApiError:
 * - logApiError = "Something broke!" (creates an error report)
 * - logApiBreadcrumb = "Here's what happened" (adds a note to the timeline)
 */
export function logApiBreadcrumb(context: ApiBreadcrumbContext): void {
  Sentry.addBreadcrumb({
    category: 'api',                    // Groups it with other API events
    message: context.message,           // What happened (human-readable)
    level: 'info',                      // This is informational, not an error
    data: {
      url: context.url,                 // Which API endpoint
      method: context.method,           // GET, POST, etc.
      ...context.data,                  // Any extra information you want to include
    },
  });
}

/**
 * Log a successful API call (optional - for performance monitoring)
 *
 * WHAT:
 * Records successful API calls so we may track how fast the APIs are responding.
 * Ideally helps identify slow endpoints before users complain.
 *
 * WHEN:
 * Call this manually after successful API calls, or add it to the axios interceptor.
 *
 * INFO DISPLAYED IN SENTRY:
 * - Which API call succeeded
 * - How long it took
 * - The response status
 *
 * EXAMPLE:
 * User loads their profile successfully in 250ms.
 * In Sentry we should see: "GET /api/users succeeded in 250ms"
 */
export function logApiSuccess(context: {
  url?: string;        // Which API endpoint
  method?: string;     // What type of request
  status?: number;     // HTTP status (usually 200)
  responseTime?: number; // How long it took in milliseconds
}): void {
  Sentry.addBreadcrumb({
    category: 'api.success',            // This groups it with other API successes
    message: `API Success: ${context.method} ${context.url}`, // Human-readable message
    level: 'info',                      // This is informational, not an error
    data: {
      url: context.url,                 // Which endpoint
      method: context.method,           // GET, POST, etc.
      status: context.status,           // Usually 200
      response_time_ms: context.responseTime, // How long it took
    },
  });
}

/**
 * Log a slow API response (performance warning)
 *
 * WHAT:
 * Warns you when an API call takes longer than expected.
 * This helps you identify performance problems before they become user complaints.
 *
 * WHEN:
 * When an API call takes longer than your performance threshold (e.g., 5 seconds).
 *
 * INFO DISPLAYED IN SENTRY:
 * - Which API call was slow
 * - How long it took
 * - What your threshold was
 *
 * EXAMPLE:
 * User uploads a file, it takes 6.5 seconds (over the 5-second limit).
 * In Sentry we should see: "Warning: POST /api/upload took 6.5 seconds (threshold: 5s)"
 */
export function logSlowApiResponse(context: {
  url?: string;        // Which API endpoint was slow
  method?: string;     // What type of request
  responseTime: number; // How long it actually took
  threshold: number;   // What your performance limit is
}): void {
  Sentry.addBreadcrumb({
    category: 'performance',            // This groups it with other performance issues
    message: 'Slow API response detected', // Human-readable warning
    level: 'warning',                   // This is a warning, not an error
    data: {
      url: context.url,                 // Which endpoint was slow
      method: context.method,           // GET, POST, etc.
      response_time_ms: context.responseTime, // How long it took
      threshold_ms: context.threshold,  // What your limit was
    },
  });
}