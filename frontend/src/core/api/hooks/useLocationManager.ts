import { LocationService } from '../services/locationService';
import type { BaseLocation, PhysicalLocation, Sublocation } from '../types/location';
import { useBackendMutation } from './useBackendMutation';
import { mediaStorageKeys } from '../constants/query-keys/mediaStorage';
import { useQueryClient } from '@tanstack/react-query';
import type { ApiResponse } from '../types/api.types';

interface UseLocationManagerOptions {
  type: 'physical' | 'sublocation';
  onSuccess?: (data: PhysicalLocation | Sublocation | { id: string }) => void;
  onError?: (error: Error) => void;
}

export function useLocationManager({ type, onSuccess, onError }: UseLocationManagerOptions) {
  const locationService = LocationService.getInstance();
  const queryClient = useQueryClient();

  const createMutation = useBackendMutation<PhysicalLocation | Sublocation, BaseLocation>(
    async (location: BaseLocation) => {
      const result = await locationService.createLocation(location);
      return { success: true, data: result };
    },
    {
      onSuccess: (data: ApiResponse<PhysicalLocation | Sublocation>) => {
        // Invalidate the locations query to trigger a refetch
        queryClient.invalidateQueries({ queryKey: mediaStorageKeys.locations.all });
        onSuccess?.(data.data);
      },
      onError: (error: Error) => {
        onError?.(error);
      }
    }
  );

  const updateMutation = useBackendMutation<PhysicalLocation | Sublocation, BaseLocation>(
    async (location: BaseLocation) => {
      const result = await locationService.updateLocation(location);
      return { success: true, data: result };
    },
    {
      onSuccess: (data: ApiResponse<PhysicalLocation | Sublocation>) => {
        // Invalidate the locations query to trigger a refetch
        queryClient.invalidateQueries({ queryKey: mediaStorageKeys.locations.all });
        onSuccess?.(data.data);
        },
        onError: (error: Error) => {
        onError?.(error);
        }
    }
  );

  const deleteMutation = useBackendMutation<void, string>(
    async (id: string) => {
      await locationService.deleteLocation(id, type);
      return { success: true, data: undefined };
    },
    {
      onSuccess: (_, id: string) => {
        // Invalidate the locations query to trigger a refetch
        queryClient.invalidateQueries({ queryKey: mediaStorageKeys.locations.all });
        onSuccess?.({ id });
        },
      onError: (error: Error) => {
        onError?.(error);
        }
    }
  );

  return {
    create: createMutation.mutate,
    update: updateMutation.mutate,
    delete: deleteMutation.mutate,
    isCreating: createMutation.isPending,
    isUpdating: updateMutation.isPending,
    isDeleting: deleteMutation.isPending,
  };
}

// For backward compatibility
export const useLocationDelete = (type: 'physical' | 'sublocation', onSuccess?: (data: { id: string }) => void, onError?: (error: Error) => void) => {
  const { delete: deleteLocation, isDeleting } = useLocationManager({
    type,
    onSuccess: (data) => {
      if ('id' in data && data.id) {
        onSuccess?.({ id: data.id });
      }
    },
    onError
  });
  return { deleteLocation, isDeleting };
};