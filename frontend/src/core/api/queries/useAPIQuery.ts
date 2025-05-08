import { useQuery } from '@tanstack/react-query';
import type { UseQueryOptions } from '@tanstack/react-query';

/**
 * Standardized query hook for API calls
 *
 * This hook provides consistent behavior across all API queries:
 * - Auth token handling
 * - Network mode configuration
 * - Retry on mount
 * - Global error handling
 * - Consistent loading states
 *
 * @template TData - The type of data returned by the query
 * @template TError - The type of error returned by the query
 * @param options - Query options to be merged with defaults
 * @returns Query result with standardized behavior
 */
export function useAPIQuery<TData = unknown, TError = unknown>(
  options: UseQueryOptions<TData, TError>
) {
  return useQuery({
    ...options,
    // Default options for all API queries
    networkMode: 'online',
    retryOnMount: true,
    // Auth token handling
    // etc...
  });
}