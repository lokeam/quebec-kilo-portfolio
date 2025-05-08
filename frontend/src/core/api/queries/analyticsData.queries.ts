/**
 * Analytics Data Query Hooks
 *
 * Provides React Query hooks for fetching analytics data across different domains.
 */

import { useAPIQuery } from './useAPIQuery';
import { getAnalyticsData, type AnalyticsDomain, type AnalyticsResponse } from '../services/analytics.service';
import { logger } from '@/core/utils/logger/logger';
import { QUERY_STALE_TIME } from '../config';
import {
  adaptAnalyticsToStorageMetadata,
  adaptAnalyticsToPhysicalLocations,
  adaptAnalyticsToDigitalLocations
} from '../adapters/analytics.adapter';
import { useMemo } from 'react';

/**
 * Query key factory for analytics queries
 */
export const analyticsKeys = {
  all: ['analytics'] as const,
  lists: () => [...analyticsKeys.all, 'list'] as const,
  list: (domains: AnalyticsDomain[]) => [...analyticsKeys.lists(), { domains }] as const,
};

/**
 * Hook for fetching analytics data for specified domains
 *
 * @function useAnalyticsData
 * @param {AnalyticsDomain[]} domains - Array of analytics domains to fetch (e.g., ['storage', 'general'])
 * @returns {UseQueryResult} React Query result object containing the analytics data
 *
 * @remarks
 * Responsibilities:
 * - Fetches analytics data for specified domains
 * - Configurable stale time
 * - Proper error handling and logging
 *
 * @example
 * ```typescript
 * function AnalyticsComponent() {
 *   const { data, isLoading, error } = useAnalyticsData(['storage', 'general']);
 *
 *   if (isLoading) return <div>Loading...</div>;
 *   if (error) return <div>Error: {error.message}</div>;
 *
 *   return (
 *     <div>
 *       <h2>Storage Stats</h2>
 *       <p>Total Physical Locations: {data?.storage?.totalPhysicalLocations}</p>
 *       <p>Total Digital Locations: {data?.storage?.totalDigitalLocations}</p>
 *     </div>
 *   );
 * }
 * ```
 */
export function useAnalyticsData(domains: AnalyticsDomain[] = []) {
  return useAPIQuery<AnalyticsResponse>({
    queryKey: analyticsKeys.list(domains),
    queryFn: async () => {
      try {
        const response = await getAnalyticsData(domains);
        logger.debug('Analytics data fetched successfully', {
          domains,
          dataKeys: Object.keys(response || {})
        });
        return response;
      } catch (error) {
        logger.error('Failed to fetch analytics data', { domains, error });
        throw error;
      }
    },
    staleTime: QUERY_STALE_TIME,
  });
}

/**
 * Hook to fetch storage analytics data specifically
 *
 * This is a convenience wrapper around useAnalyticsData specifically for
 * digital and physical location information.
 *
 * @returns Query result with storage analytics data
 */
export function useStorageAnalytics() {
  const { data, isLoading, error } = useAnalyticsData(['storage']);

  const transformedData = useMemo(() => {
    if (!data) return null;

    return {
      metadata: adaptAnalyticsToStorageMetadata(data),
      physicalLocations: adaptAnalyticsToPhysicalLocations(data),
      digitalLocations: adaptAnalyticsToDigitalLocations(data)
    };
  }, [data]);

  return {
    data: transformedData,
    isLoading,
    error
  };
}

/**
 * Hook to fetch general analytics data
 *
 * This is a convenience wrapper for dashboard overview data.
 *
 * @returns Query result with general analytics data
 */
export function useGeneralAnalytics() {
  return useAnalyticsData(['general']);
}

/**
 * Hook to fetch financial analytics data
 *
 * This is a convenience wrapper for subscription and cost information.
 *
 * @returns Query result with financial analytics data
 */
export function useFinancialAnalytics() {
  return useAnalyticsData(['financial']);
}