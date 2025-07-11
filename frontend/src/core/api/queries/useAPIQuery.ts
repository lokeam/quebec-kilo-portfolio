/**
 * BACKEND PROTECTION SYSTEM - Safe Query Hook
 *
 * Custom version of React Query that hopefully mitigates hammering of the backend.
 *
 * WHAT IT DOES:
 * - Checks if an API call is blocked before firing
 * - Records when API calls succeed / fail
 * - Uses safe retry settings (doesn't retry forever)
 * - Turns off aggressive refetching (doesn't retry when user clicks back to tab)
 *
 * HOW TO USE:
 * Use THIS instead of directly using React Query's useQuery
 *
 * OLD (dangerous):
 * const { data } = useQuery(['user-profile'], getUserProfile);
 *
 * NEW (protected):
 * const { data } = useAPIQuery({
 *   queryKey: ['user-profile'],
 *   queryFn: getUserProfile,
 * });
 *
 * WHY THIS MATTERS:
 * - See same section in rateLimiter.ts
 */

import { useQuery } from '@tanstack/react-query';
import type { UseQueryOptions } from '@tanstack/react-query';
import { isRateLimited, recordQueryFailure, recordQuerySuccess } from '@/core/api/queries/rateLimiter';
import { logger } from '@/core/utils/logger/logger';

/**
 * A safer version of React Query that prevents backend hammering
 *
 * @param options - The same options we'd pass to useQuery, but with added protection
 * @returns The same result as useQuery, but with rate limiting protection
 *
 * EXAMPLE USAGE:
 * const { data, isLoading, error } = useAPIQuery({
 *   queryKey: ['user-profile'],
 *   queryFn: () => fetch('/api/user/profile').then(r => r.json()),
 *   staleTime: 5 * 60 * 1000, // 5 minutes
 * });
 */
export function useAPIQuery<TData = unknown, TError = unknown>(
  options: UseQueryOptions<TData, TError>
) {
  const { queryKey, queryFn, ...rest } = options;

  // Check if this query is currently blocked due to too many failures
  const enabled = !isRateLimited(queryKey) && (options.enabled ?? true);

  return useQuery({
    ...rest,
    queryKey,
    queryFn: async (context) => {
      // Double-check if we're blocked before making the request
      if (isRateLimited(queryKey)) {
        logger.warn('Query blocked by rate limiter', { queryKey });
        throw new Error('Rate limited - too many consecutive failures');
      }

      try {
        // Make sure we have a valid function to call
        if (typeof queryFn !== 'function') {
          throw new Error('queryFn must be a function');
        }

        // Make the actual API call
        const result = await queryFn(context);

        // Record success - this resets the failure count
        recordQuerySuccess(queryKey);

        return result;
      } catch (error) {
        // Record failure - will eventually block the query if it fails too often
        recordQueryFailure(queryKey);

        // Re-throw the error so React Query can handle it
        throw error;
      }
    },
    enabled,

    // SAFE DEFAULTS - These prevent aggressive retrying
    networkMode: 'online',
    retryOnMount: true,

    // SMART RETRY LOGIC - Only retry certain errors, do not continue forever
    retry: (failureCount: number, error: unknown) => {
      const err = error as { status?: number };

      // Don't retry auth errors - they won't fix themselves
      if (err?.status === 401 || err?.status === 403) return false;

      // Don't retry 404s - the route doesn't exist
      if (err?.status === 404) return false;

      // Only retry server errors 2 times max
      if (err?.status && err.status >= 500 && failureCount >= 2) return false;

      // Max 3 retries total for other errors
      return failureCount < 3;
    },

    // Capped exponential backoff - don't wait forever between retries
    retryDelay: attemptIndex => Math.min(1000 * 2 ** attemptIndex, 5000),

    // Turn off aggressive refetching that could hammer backend
    refetchOnWindowFocus: false,  // Don't retry when user clicks back to tab
    refetchOnReconnect: false,    // Don't retry when network reconnects
  });
}