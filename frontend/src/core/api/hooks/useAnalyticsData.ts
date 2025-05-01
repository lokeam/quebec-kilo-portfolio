/**
 * Hook for fetching analytics data from the server.
 */

import { useBackendQuery } from '@/core/api/hooks/useBackendQuery';
import {
  getAnalyticsData,
  type AnalyticsDomain,
  type AnalyticsResponse
} from '@/core/api/services/analytics.service';

/**
 * Hook to fetch analytics data for specified domains.
 *
 * @param domains - Array of analytics domains to fetch (e.g., ['storage', 'general'])
 * @returns Query result with analytics data, loading state, and error information
 *
 * @example
 * ```tsx
 * // Fetch storage analytics for the digital locations page
 * const { data, isLoading } = useAnalyticsData(['storage']);
 *
 * // Fetch multiple domains for the dashboard
 * const { data } = useAnalyticsData(['general', 'financial']);
 * ```
 */
export function useAnalyticsData(domains: AnalyticsDomain[] = []) {
  return useBackendQuery<AnalyticsResponse>({
    // Use domains in the query key for proper cache management
    queryKey: ['analytics', ...domains],
    queryFn: async (getToken) => {
      const token = await getToken();
      return getAnalyticsData(domains, token);
    },
    // Only refetch after 5 minutes - analytics don't change that frequently
    staleTime: 5 * 60 * 1000,
  });
}

/**
 * Hook to fetch storage analytics data specifically.
 *
 * This is a convenience wrapper around useAnalyticsData specifically for
 * digital and physical location information.
 *
 * @returns Query result with storage analytics data
 */
export function useStorageAnalytics() {
  return useAnalyticsData(['storage']);
}

/**
 * Hook to fetch general analytics data.
 *
 * This is a convenience wrapper for dashboard overview data.
 *
 * @returns Query result with general analytics data
 */
export function useGeneralAnalytics() {
  return useAnalyticsData(['general']);
}

/**
 * Hook to fetch financial analytics data.
 *
 * This is a convenience wrapper for subscription and cost information.
 *
 * @returns Query result with financial analytics data
 */
export function useFinancialAnalytics() {
  return useAnalyticsData(['financial']);
}