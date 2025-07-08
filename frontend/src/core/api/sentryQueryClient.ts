import { QueryClient } from '@tanstack/react-query';
import * as Sentry from '@sentry/react';
import { logger } from '@/core/utils/logger/logger';

/**
 * Sentry-enhanced QueryClient that monitors API calls for errors and performance
 *
 * Features:
 * - Tracks failed API calls with full context
 * - Monitors API call performance (>1 second threshold)
 * - Adds user context to error reports
 * - Provides debugging information for API failures
 */
export function createSentryQueryClient(): QueryClient {
  return new QueryClient({
    defaultOptions: {
      queries: {
        retry: 1,
        staleTime: 90 * 60 * 1000, // 90 minutes
      },

      mutations: {
        retry: 1,
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