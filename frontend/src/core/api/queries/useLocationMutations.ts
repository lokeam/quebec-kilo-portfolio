import { axiosInstance } from '@/core/api/client/axios-instance';
import { useQueryClient } from '@tanstack/react-query';
import { useBackendMutation } from '@/core/api/hooks/useBackendMutation';
import { API_ROUTES } from '@/core/api/constants/routes';
import type { PhysicalLocation } from '@/features/dashboard/lib/types/media-storage/physical';
import { logger } from '@/core/utils/logger/logger';
import { mediaStorageKeys } from '../constants/query-keys/mediaStorage';

interface LocationResponse {
  success: boolean;
  location: PhysicalLocation;
}

interface DeleteLocationResponse {
  success: boolean;
  id: string;
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
    async (locationID, token) => {
      return axiosInstance.delete(
        API_ROUTES.LOCATIONS.BY_ID(locationID),
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
  )
}
