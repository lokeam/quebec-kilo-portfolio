/**
 * Digital Services Catalog Query Hooks
 *
 * Provides React Query hooks for fetching digital services catalog data.
 */

import { useAPIQuery } from '@/core/api/queries/useAPIQuery';
import { getServicesCatalog, FALLBACK_SERVICES } from '@/core/api/services/digitalServicesCatalog.service';
import { logger } from '@/core/utils/logger/logger';
import type { DigitalServiceItem } from '@/types/domain/digital-location';

/**
 * Query key factory for digital services catalog queries
 */
export const digitalServicesCatalogKeys = {
  all: ['digitalServicesCatalog'] as const,
  lists: () => [...digitalServicesCatalogKeys.all, 'list'] as const,
  list: () => [...digitalServicesCatalogKeys.lists()] as const,
};

/**
 * Hook for fetching digital services catalog
 *
 * @returns Query result with digital services data
 */
export function useGetDigitalServicesCatalog() {
  return useAPIQuery<DigitalServiceItem[]>({
    queryKey: digitalServicesCatalogKeys.list(),
    queryFn: async () => {
      try {
        console.log('[DEBUG] useGetDigitalServicesCatalog: Starting API call');
        const services = await getServicesCatalog();
        console.log('[DEBUG] useGetDigitalServicesCatalog: API call successful, got', services?.length || 0, 'services');
        console.log('[DEBUG] useGetDigitalServicesCatalog: Services data:', services);

        // Ensure we always return an array
        return services || [];
      } catch (error) {
        console.error('[DEBUG] useGetDigitalServicesCatalog: API call failed:', error);
        logger.error('Failed to fetch digital services catalog', { error });
        // Return fallback data in case of error
        console.log('[DEBUG] useGetDigitalServicesCatalog: Returning fallback services');
        return FALLBACK_SERVICES;
      }
    },
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
}
