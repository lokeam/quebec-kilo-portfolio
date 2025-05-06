/**
 * Hook for fetching digital locations from the server.
 */

import { useBackendQuery } from '@/core/api/hooks/useBackendQuery';
import { getUserDigitalLocations, getDigitalLocationById } from '@/core/api/services/mediaStorage.service';
import { adaptDigitalLocationToService } from '@/core/api/adapters/digitalLocation.adapter';
import type { OnlineService } from '@/features/dashboard/lib/types/online-services/services';
import { type DigitalLocation } from '@/features/dashboard/lib/types/media-storage/digital-location.types';

/**
 * Hook to fetch all digital locations for the current user.
 * Transforms API responses into the format expected by the UI components.
 *
 * @returns Query result with transformed digital locations data, loading state, and error information
 */
export function useDigitalLocations() {
  return useBackendQuery<OnlineService[]>({
    queryKey: ['digitalLocations'],
    queryFn: async (getToken) => {
      const token = await getToken();
      const locations = await getUserDigitalLocations(token);
      return locations.map(adaptDigitalLocationToService);
    },
  });
}

/**
 * Hook to fetch a specific digital location by ID.
 *
 * @param locationId - The ID of the digital location to fetch
 * @returns Query result with a single digital location, loading state, and error information
 */
export function useDigitalLocationById(locationId: string) {
  return useBackendQuery<DigitalLocation>({
    queryKey: ['digitalLocations', locationId],
    queryFn: async (getToken) => {
      const token = await getToken();
      return getDigitalLocationById(locationId, token);
    },
    enabled: !!locationId, // Only run the query if locationId exists
  });
}