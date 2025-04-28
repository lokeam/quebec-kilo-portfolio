/**
 * Hook for fetching physical locations from the server.
 */

import { useBackendQuery } from '@/core/api/hooks/useBackendQuery';
import { getUserPhysicalLocations, getPhysicalLocationById } from '@/core/api/services/physicalLocations.service';
import { mediaStorageKeys } from '@/core/api/constants/query-keys/mediaStorage';
import type { PhysicalLocation } from '@/features/dashboard/lib/types/media-storage/physical';

/**
 * Hook to fetch all physical locations for the current user.
 *
 * @returns Query result with physical locations data, loading state, and error information
 */
export function usePhysicalLocations() {
  return useBackendQuery<PhysicalLocation[]>({
    queryKey: mediaStorageKeys.locations.all,
    queryFn: async (getToken) => {
      const token = await getToken();
      return getUserPhysicalLocations(token);
    },
  });
}

/**
 * Hook to fetch a specific physical location by ID.
 *
 * @param locationId - The ID of the physical location to fetch
 * @returns Query result with a single physical location, loading state, and error information
 */
export function usePhysicalLocationById(locationId: string) {
  return useBackendQuery<PhysicalLocation>({
    queryKey: mediaStorageKeys.locations.byId(locationId),
    queryFn: async (getToken) => {
      const token = await getToken();
      return getPhysicalLocationById(locationId, token);
    },
    enabled: !!locationId, // Only run the query if locationId exists
  });
}