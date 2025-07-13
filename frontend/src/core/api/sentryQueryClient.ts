import { QueryClient } from '@tanstack/react-query';
import * as Sentry from '@sentry/react';
import { logger } from '@/core/utils/logger/logger';

/**
 * Production-safe QueryClient with circuit breakers and retry protection
 *
 * PROTECTION LAYERS:
 * 1. No retries on auth errors (401/403)
 * 2. No retries on 404s (route doesn't exist)
 * 3. Limited retries on 5xx errors (max 2)
 * 4. Exponential backoff capped at 5 seconds
 * 5. Disabled aggressive refetching
 * 6. Request deduplication and caching
 */
export function createSentryQueryClient(): QueryClient {
  return new QueryClient({
    defaultOptions: {
      queries: {
        retry: (failureCount, error: unknown) => {
          const errorObj = error as { status?: number };

        // NEVER retry auth errors - they won't fix themselves
        if (errorObj?.status === 401 || errorObj?.status === 403) {
          logger.warn('Auth error - not retrying', { status: errorObj.status || 'unknown' });
          return false;
        }

        // NEVER retry 404s - route doesn't exist
        if (errorObj?.status === 404) {
          logger.warn('404 error - not retrying', { status: errorObj.status || 'unknown' });
          return false;
        }

        // Stop after 2 retries for server errors
        if (errorObj?.status && errorObj.status >= 500 && failureCount >= 2) {
          logger.error('Server error - stopping retries', {
            status: errorObj.status,
            failureCount
          });
          return false;
        }

        // Max 3 retries total
        return failureCount < 3;
      },

      // CAPPED EXPONENTIAL BACKOFF (prevents long delays)
      retryDelay: attemptIndex => Math.min(1000 * 2 ** attemptIndex, 5000),

      // AGGRESSIVE CACHING (reduces requests)
      staleTime: 10 * 60 * 1000, // 10 minutes
      gcTime: 30 * 60 * 1000,    // 30 minutes

      // DISABLE AGGRESSIVE REFETCHING (prevents hammering)
      refetchOnWindowFocus: false,
      refetchOnReconnect: false,
      refetchOnMount: false
      },

      mutations: {
        // NEVER retry mutations (they're usually destructive)
        retry: false,
      },
    },
  });
}

/**
 * Utility function to track API errors in Sentry
 * Call this in your API hooks when errors occur
 */
export function trackApiError(
  error: unknown,
  context: {
    queryKey?: unknown[];
    mutationKey?: unknown[];
    variables?: unknown;
    operation?: string;
  }
): void {
  try {
    // Add API context to error
    Sentry.setContext('api', {
      queryKey: context.queryKey,
      mutationKey: context.mutationKey,
      variables: context.variables,
      operation: context.operation,
      error: error instanceof Error ? error.message : String(error),
    });

    // Capture the error with context
    Sentry.captureException(error, {
      tags: {
        type: 'api_error',
        operation: context.operation || 'unknown',
      },
    });

    logger.error('API error captured in Sentry', {
      context,
      error: error instanceof Error ? error.message : String(error),
    });
  } catch (sentryError) {
    // Don't let Sentry errors break the app
    logger.error('Failed to track error in Sentry', { sentryError });
  }
}

/**
 * Utility function to track slow API calls
 * Call this in your API hooks for performance monitoring
 */
export function trackSlowApiCall(
  duration: number,
  context: {
    queryKey?: unknown[];
    mutationKey?: unknown[];
    operation?: string;
  }
): void {
  if (duration > 1000) { // >1 second threshold
    Sentry.addBreadcrumb({
      category: 'performance',
      message: `Slow API call: ${context.operation || 'unknown'} took ${duration}ms`,
      level: 'warning',
      data: {
        duration,
        threshold: 1000,
        ...context,
      },
    });

    logger.warn('Slow API call detected', {
      duration,
      context,
    });
  }
}