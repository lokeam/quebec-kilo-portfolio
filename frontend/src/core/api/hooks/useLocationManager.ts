import { LocationService } from '../services/locationService';
import type { BaseLocation, PhysicalLocation, Sublocation } from '../types/location';
import { useBackendMutation } from './useBackendMutation';
import { mediaStorageKeys } from '../constants/query-keys/mediaStorage';
import { useQueryClient } from '@tanstack/react-query';
import type { ApiResponse } from '../types/api.types';
import { useQuery } from '@tanstack/react-query';

interface UseLocationManagerOptions {
  type: 'physical' | 'sublocation';
  onSuccess?: (data: PhysicalLocation | Sublocation | { id: string }) => void;
  onError?: (error: Error) => void;
}

export function useLocationManager({ type, onSuccess, onError }: UseLocationManagerOptions) {
  const locationService = LocationService.getInstance();
  const queryClient = useQueryClient();

  const { data: locations } = useQuery({
    queryKey: type === 'physical' ? mediaStorageKeys.locations.all : mediaStorageKeys.sublocations.all,
    queryFn: async () => {
      if (type === 'physical') {
        const result = await locationService.getPhysicalLocations();
        return { success: true, data: result };
      } else {
        const result = await locationService.getSublocations();
        return { success: true, data: result };
      }
    }
  });

  const createMutation = useBackendMutation<PhysicalLocation | Sublocation, BaseLocation>(
    async (location: BaseLocation) => {
      if (type === 'physical') {
        const result = await locationService.createPhysicalLocation(location as PhysicalLocation);
        return { success: true, data: result };
      } else {
        const result = await locationService.createSublocation(location as Sublocation);
        return { success: true, data: result };
      }
    },
    {
      onSuccess: (data: ApiResponse<PhysicalLocation | Sublocation>) => {
        // Invalidate both physical locations and sublocations queries
        queryClient.invalidateQueries({ queryKey: mediaStorageKeys.locations.all });
        queryClient.invalidateQueries({ queryKey: mediaStorageKeys.sublocations.all });
        onSuccess?.(data.data);
      },
      onError: (error: Error) => {
        onError?.(error);
      }
    }
  );

  const updateMutation = useBackendMutation<PhysicalLocation | Sublocation, BaseLocation>(
    async (location: BaseLocation) => {
      if (type === 'physical') {
        const result = await locationService.updatePhysicalLocation(location as PhysicalLocation);
        return { success: true, data: result };
      } else {
        const result = await locationService.updateSublocation(location as Sublocation);
        return { success: true, data: result };
      }
    },
    {
      onSuccess: (data: ApiResponse<PhysicalLocation | Sublocation>) => {
        // Invalidate both physical locations and sublocations queries
        queryClient.invalidateQueries({ queryKey: mediaStorageKeys.locations.all });
        queryClient.invalidateQueries({ queryKey: mediaStorageKeys.sublocations.all });
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
        // Invalidate both physical locations and sublocations queries
        queryClient.invalidateQueries({ queryKey: mediaStorageKeys.locations.all });
        queryClient.invalidateQueries({ queryKey: mediaStorageKeys.sublocations.all });
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
    locations: locations?.data || [],
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