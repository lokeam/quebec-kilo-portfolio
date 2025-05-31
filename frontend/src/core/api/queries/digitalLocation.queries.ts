/**
 * Digital Location Query Hooks
 *
 * Provides React Query hooks for fetching and managing digital media storage locations.
 */

// React Query
import { useMutation, useQueryClient } from '@tanstack/react-query';

// Base Query Hook
import { useAPIQuery } from '@/core/api/queries/useAPIQuery';

// Service Layer methods
import {
  getAllDigitalLocations,
  getSingleDigitalLocation,
  createDigitalLocation,
  updateDigitalLocation,
  deleteDigitalLocation,
} from '@/core/api/services/digitalLocation.service';

// Utils
import { showToast } from '@/shared/components/ui/TanstackMutationToast/showToast';

// Types
import type { DigitalLocation, CreateDigitalLocationRequest } from '@/types/domain/digital-location';
import { TOAST_SUCCESS_MESSAGES } from '@/shared/constants/toast.success.messages';
import { TOAST_DURATIONS } from '@/shared/constants/toast.config';
import { TOAST_ERROR_MESSAGES } from '@/shared/constants/toast.error.messages';

/**
 * Query key factory for digital location queries
 */
export const digitalLocationKeys = {
  all: ['digital-locations'] as const,
  lists: () => [...digitalLocationKeys.all, 'list'] as const,
  list: (filters: string) => [...digitalLocationKeys.lists(), { filters }] as const,
  details: () => [...digitalLocationKeys.all, 'detail'] as const,
  detail: (id: string) => [...digitalLocationKeys.details(), id] as const,
};

/**
 * Hook to fetch all digital locations
 */
export const useGetAllDigitalLocations = () => {
  return useAPIQuery<DigitalLocation[]>({
    queryKey: digitalLocationKeys.lists(),
    queryFn: async () => {
      const locations = await getAllDigitalLocations();
      return locations;
    },
  });
};

/**
 * Hook to fetch a single digital location by ID
 */
export const useGetSingleDigitalLocation = (id: string) => {
  return useAPIQuery<DigitalLocation>({
    queryKey: digitalLocationKeys.detail(id),
    queryFn: async () => {
      const singleLocation = await getSingleDigitalLocation(id);
      return singleLocation;
    },
    enabled: !!id,
  });
};

/**
 * Hook to create a new digital location
 */
export const useCreateDigitalLocation = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateDigitalLocationRequest) => createDigitalLocation(data),
    onSuccess: () => {
      // Invalidate queries that need to be refreshed
      queryClient.invalidateQueries({ queryKey: digitalLocationKeys.lists() });
      queryClient.invalidateQueries({ queryKey: ['analytics'] });

      showToast({
        message: TOAST_SUCCESS_MESSAGES.DIGITAL_LOCATION.CREATE,
        variant: 'success',
        duration: TOAST_DURATIONS.EXTENDED,
      });
    },
    onError: (error) => {
      console.error('Failed to create digital location:', error);
      showToast({
        message: TOAST_ERROR_MESSAGES.DIGITAL_LOCATION.CREATE.DEFAULT,
        variant: 'error',
        duration: TOAST_DURATIONS.EXTENDED,
      });
    }
  });
};

/**
 * Hook to update an existing digital location
 */
export const useUpdateDigitalLocation = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<CreateDigitalLocationRequest> }) => {
      // NOTE: double check that we don't need to adapt data here
      return updateDigitalLocation(id, data);
    },
    onSuccess: (data) => {
      queryClient.invalidateQueries({ queryKey: digitalLocationKeys.detail(data.id) });
      queryClient.invalidateQueries({ queryKey: digitalLocationKeys.lists() });
    },
  });
};

/**
 * Hook to delete a digital location
 */
export const useDeleteDigitalLocation = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => deleteDigitalLocation(id),
    onSuccess: (_, id) => {
      queryClient.invalidateQueries({ queryKey: digitalLocationKeys.detail(id) });
      queryClient.invalidateQueries({ queryKey: digitalLocationKeys.lists() });
    },
  });
};
