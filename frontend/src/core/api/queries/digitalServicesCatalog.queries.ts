/**
 * Digital Services Catalog Query Hooks
 *
 * Provides React Query hooks for fetching digital services catalog data.
 */

import { useAPIQuery } from '@/core/api/queries/useAPIQuery';
import { getServicesCatalog, FALLBACK_SERVICES } from '@/core/api/services/digitalServicesCatalog.service';
import { adaptDigitalServicesToLocations } from '@/core/api/adapters/digitalServiceCatalog.adapter';
import { logger } from '@/core/utils/logger/logger';
import type { DigitalLocation } from '@/types/domain/online-service';

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
 * @returns Query result with transformed digital services data
 */
export function useGetDigitalServicesCatalog() {
  return useAPIQuery<DigitalLocation[]>({
    queryKey: digitalServicesCatalogKeys.list(),
    queryFn: async () => {
      try {
        const services = await getServicesCatalog();
        return adaptDigitalServicesToLocations(services);
      } catch (error) {
        logger.error('Failed to fetch digital services catalog', { error });
        // Return fallback data in case of error
        return adaptDigitalServicesToLocations(FALLBACK_SERVICES);
      }
    },
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
}
