/**
 * Sublocation Mutations
 *
 * For API standards and best practices, see:
 * @see {@link ../../docs/api-standards.md}
 */

import { axiosInstance } from '@/core/api/client/axios-instance';
import { useQueryClient } from '@tanstack/react-query';
import { useBackendMutation } from '@/core/api/hooks/useBackendMutation';
import { API_ROUTES } from '@/core/api/constants/routes';
import type { Sublocation } from '@/features/dashboard/lib/types/media-storage/sublocation';
import { logger } from '@/core/utils/logger/logger';
import { mediaStorageKeys } from '../constants/query-keys/mediaStorage';

export interface SublocationResponse {
  success: boolean;
  sublocation: Sublocation;
}

interface DeleteSublocationResponse {
  success: boolean;
  id: string;
}

export function useCreateSublocationMutation() {
  const queryClient = useQueryClient();

  return useBackendMutation<SublocationResponse, Partial<Sublocation>>(
    async (sublocationData, token) => {
      return axiosInstance.post(
        API_ROUTES.SUBLOCATION.CREATE,
        sublocationData,
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );
    },
    {
      onSuccess: (data) => {
        // Invalidate all location data
        queryClient.invalidateQueries({
          queryKey: mediaStorageKeys.all,
        });

        // Invalidate specific sublocation
        if (data?.data?.sublocation?.id) {
          queryClient.invalidateQueries({
            queryKey: mediaStorageKeys.sublocations.byId(data.data.sublocation.id),
          });

          // Also invalidate parent location if present
          if (data?.data?.sublocation?.parentLocationId) {
            queryClient.invalidateQueries({
              queryKey: mediaStorageKeys.locations.byId(data.data.sublocation.parentLocationId),
            });
          }
        }
      },
      onError: (error) => {
        logger.error('Failed to create sublocation:', error);
      }
    }
  );
}

export function useUpdateSublocationMutation() {
  const queryClient = useQueryClient();

  return useBackendMutation<SublocationResponse, Partial<Sublocation>>(
    async (sublocationData, token) => {
      return axiosInstance.put(
        API_ROUTES.SUBLOCATION.BY_ID(sublocationData.id as string),
        sublocationData,
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );
    },
    {
      onSuccess: (data) => {
        // Invalidate all location data
        queryClient.invalidateQueries({
          queryKey: mediaStorageKeys.all,
        });

        // Invalidate specific sublocation
        if (data?.data?.sublocation?.id) {
          queryClient.invalidateQueries({
            queryKey: mediaStorageKeys.sublocations.byId(data.data.sublocation.id),
          });

          // Also invalidate parent location if present
          if (data?.data?.sublocation?.parentLocationId) {
            queryClient.invalidateQueries({
              queryKey: mediaStorageKeys.locations.byId(data.data.sublocation.parentLocationId),
            });
          }
        }
      },
      onError: (error) => {
        logger.error('Failed to update sublocation:', error);
      }
    }
  );
}

export function useDeleteSublocationMutation() {
  const queryClient = useQueryClient();

  return useBackendMutation<DeleteSublocationResponse, string>(
    async (sublocationID, token) => {
      return axiosInstance.delete(
        API_ROUTES.SUBLOCATION.BY_ID(sublocationID),
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
        logger.error('Failed to delete sublocation:', error);
      }
    }
  );
}