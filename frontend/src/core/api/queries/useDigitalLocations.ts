/**
 * Digital Locations Query Hooks
 *
 * Provides React Query hooks for fetching and managing digital locations.
 */

import { useBackendQuery } from '@/core/api/hooks/useBackendQuery';
import { useBackendMutation } from '@/core/api/hooks/useBackendMutation';
import { digitalLocationService } from '@/core/api/services/digitalLocation.service';
import { logger } from '@/core/utils/logger/logger';
import type { DigitalLocation, CreateDigitalLocationInput, UpdateDigitalLocationInput } from '@/types/digital-location';
import { useQueryClient } from '@tanstack/react-query';

// Query key for digital locations
export const digitalLocationsQueryKey = ['digitalLocations'] as const;

/**
 * Hook for fetching all digital locations
 */
export function useGetAllDigitalLocations() {
  return useBackendQuery<DigitalLocation[]>({
    queryKey: digitalLocationsQueryKey,
    queryFn: async () => {
      logger.debug('Fetching digital locations');
      const locations = await digitalLocationService.getAllLocations();
      logger.debug('Digital locations fetched successfully', {
        count: locations.length
      });
      return locations;
    }
  });
}

/**
 * Hook for creating a new digital location
 */
export function useCreateDigitalLocation() {
  const queryClient = useQueryClient();

  return useBackendMutation<DigitalLocation, CreateDigitalLocationInput>({
    mutationFn: async (input) => {
      logger.debug('Creating digital location', { input });
      const location = await digitalLocationService.createLocation(input);
      logger.debug('Digital location created successfully', {
        id: location.id
      });
      return location;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: digitalLocationsQueryKey });
    }
  });
}

/**
 * Hook for updating a digital location
 */
export function useUpdateDigitalLocation() {
  const queryClient = useQueryClient();

  return useBackendMutation<DigitalLocation, { id: string; input: UpdateDigitalLocationInput }>({
    mutationFn: async ({ id, input }) => {
      logger.debug('Updating digital location', { id, input });
      const location = await digitalLocationService.updateLocation(id, input);
      logger.debug('Digital location updated successfully', {
        id: location.id
      });
      return location;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: digitalLocationsQueryKey });
    }
  });
}

/**
 * Hook for deleting a digital location
 */
export function useDeleteDigitalLocation() {
  const queryClient = useQueryClient();

  return useBackendMutation<void, string>({
    mutationFn: async (id) => {
      logger.debug('Deleting digital location', { id });
      await digitalLocationService.deleteLocation(id);
      logger.debug('Digital location deleted successfully', { id });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: digitalLocationsQueryKey });
    }
  });
}