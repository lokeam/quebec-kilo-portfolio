/**
 * Location Mutations
 *
 * For API standards and best practices, see:
 * @see {@link ../../docs/api-standards.md}
 */

import { axiosInstance } from '@/core/api/client/axios-instance';
import { useQueryClient } from '@tanstack/react-query';
import { useBackendMutation } from '@/core/api/hooks/useBackendMutation';
import { API_ROUTES } from '@/core/api/constants/routes';
import type { PhysicalLocation } from '@/features/dashboard/lib/types/media-storage/physical';
import { logger } from '@/core/utils/logger/logger';
import { mediaStorageKeys } from '../constants/query-keys/mediaStorage';

export interface LocationResponse {
  success: boolean;
  location: PhysicalLocation;
}

interface DeleteLocationResponse {
  success: boolean;
  id: string;
}

// Define a type that can accept both camelCase and snake_case formats
export interface LocationPayload {
  id?: string;
  name: string;
  locationType: string;
  mapCoordinates?: string;
  bgColor?: string;
  physicalLocationId?: string;
  [key: string]: unknown;
}

export function useCreateLocationMutation(onSuccess?: () => void) {
  const queryClient = useQueryClient();

  return useBackendMutation<LocationResponse, LocationPayload>(
    async (locationData, token) => {
      // Log the request URL and payload for debugging
      console.log('[DEBUG] Creating physical location with payload:', locationData);
      console.log('[DEBUG] API endpoint:', API_ROUTES.LOCATIONS.CREATE);

      return axiosInstance.post(
        API_ROUTES.LOCATIONS.CREATE,
        locationData,
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );
    },
    {
      onSuccess: (data) => {
        // Log success
        console.log('[DEBUG] Successfully created location:', data);

        // Invalidate all location data
        queryClient.invalidateQueries({
          queryKey: mediaStorageKeys.locations.all,
        });

        // Invalidate specific location
        if (data?.data?.location?.id) {
          queryClient.invalidateQueries({
            queryKey: mediaStorageKeys.locations.byId(data.data.location.id),
          });
        }

        // Call onSuccess callback if provided
        onSuccess?.();
      },
      onError: (error) => {
        // Add detailed error handling
        console.error('[DEBUG] Failed to create location, error:', error);
        logger.error('Failed to create location:', error);
      }
    }
  );
}

export function useUpdateLocationMutation() {
  const queryClient = useQueryClient();

  return useBackendMutation<LocationResponse, Partial<PhysicalLocation>>(
    async (locationData, token) => {
      return axiosInstance.put(
        API_ROUTES.LOCATIONS.BY_ID(locationData.id as string),
        locationData,
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );
    },
    {
      onSuccess: (data) => {
        // Invalidate all location data - more targeted invalidation possible
        queryClient.invalidateQueries({
          queryKey: mediaStorageKeys.all,
        });

        // Optionally update cache directly for the specific location
        if (data?.data?.location?.id) {
          queryClient.invalidateQueries({
            queryKey: mediaStorageKeys.locations.byId(data?.data?.location?.id),
          });
        }
      },
      onError: (error) => {
        // Log error
        logger.error('Failed to update location:', error);
      }
    }
  );
}

export function useDeleteLocationMutation() {
  const queryClient = useQueryClient();

  return useBackendMutation<DeleteLocationResponse, string>(
    async (locationId, token) => {
      return axiosInstance.delete(
        API_ROUTES.LOCATIONS.BY_ID(locationId),
        {
          headers: {
            Authorization: `Bearer ${token}`,
          }
        }
      );
    },
    {
      onSuccess: () => {
        // Invalidate all locations
        queryClient.invalidateQueries({
          queryKey: mediaStorageKeys.all,
        });
      },
      onError: (error) => {
        logger.error('Failed to delete location:', error);
      }
    }
  );
}
