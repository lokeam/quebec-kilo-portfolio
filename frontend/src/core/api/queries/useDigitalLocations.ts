/**
 * Digital Locations Query Hooks
 *
 * Provides React Query hooks for fetching and managing digital locations.
 */

import { useBackendQuery } from '@/core/api/hooks/useBackendQuery';
import { useBackendMutation } from '@/core/api/hooks/useBackendMutation';
import { digitalLocationService } from '@/core/api/services/digitalLocation.service';
import { logger } from '@/core/utils/logger/logger';
import type { DigitalLocation } from '@/features/dashboard/lib/types/media-storage/digital-location.types';
import type { CreateDigitalLocationRequest } from '@/features/dashboard/lib/types/media-storage/digital-location.types';
import { useQueryClient } from '@tanstack/react-query';
import { adaptDigitalLocationToService } from '@/core/api/adapters/digitalLocation.adapter';
import type { OnlineService } from '@/features/dashboard/lib/types/online-services/services';
import { toast } from 'sonner';
import { AxiosError } from 'axios';
import type { ApiError } from '@/core/api/types/api.types';

// Query key for digital locations
export const digitalLocationsQueryKey = ['digitalLocations'] as const;

/**
 * Hook for fetching all digital locations and transforming them to OnlineService format
 */
export function useGetAllDigitalLocations() {
  return useBackendQuery<OnlineService[]>({
    queryKey: digitalLocationsQueryKey,
    queryFn: async () => {
      logger.debug('Fetching digital locations');
      try {
        const locations = await digitalLocationService.getAllLocations();
        logger.debug('Digital locations fetched successfully', {
          count: locations.length
        });
        return locations.map(adaptDigitalLocationToService);
      } catch (error) {
        logger.error('Failed to fetch digital locations', { error });
        throw error;
      }
    },
    retry: 1 // Add a single retry for network glitches
  });
}

/**
 * Hook for creating a new digital location
 */
export function useCreateDigitalLocation() {
  const queryClient = useQueryClient();

  return useBackendMutation<DigitalLocation, CreateDigitalLocationRequest>(
    async (input: CreateDigitalLocationRequest) => {
      logger.debug('Creating digital location', { input });
      try {
        const location = await digitalLocationService.createLocation(input);
        logger.debug('Digital location created successfully', {
          id: location.id
        });
        return { data: location };
      } catch (error) {
        logger.error('Failed to create digital location', { error, input });
        throw error;
      }
    },
    {
      onSuccess: () => {
        queryClient.invalidateQueries({ queryKey: digitalLocationsQueryKey });
        toast.success('Digital service created', {
          description: 'Your digital service has been successfully created.',
          duration: 3000
        });
      },
      onError: (error: AxiosError<ApiError>) => {
        let errorMessage = 'Failed to create service';
        let errorDescription = "We couldn't create your digital service. Please try again.";

        if (error.response) {
          const status = error.response.status;
          if (status === 401 || status === 403) {
            errorMessage = 'Permission denied';
            errorDescription = "You don't have permission to create digital services.";
          } else if (status === 409) {
            errorMessage = 'Service already exists';
            errorDescription = "A service with these details already exists.";
          } else if (status >= 500) {
            errorMessage = 'Server error';
            errorDescription = "The server encountered an error. Please try again later.";
          }
        } else if (error.request) {
          errorMessage = 'Network error';
          errorDescription = "Couldn't connect to the server. Please check your internet connection.";
        }

        toast.error(errorMessage, {
          description: errorDescription,
          duration: 5000
        });
      }
    }
  );
}

/**
 * Hook for updating a digital location
 */
export function useUpdateDigitalLocation() {
  const queryClient = useQueryClient();

  return useBackendMutation<DigitalLocation, { id: string; input: Partial<CreateDigitalLocationRequest> }>(
    async ({ id, input }) => {
      logger.debug('Updating digital location', { id, input });
      try {
        const location = await digitalLocationService.updateLocation(id, input);
        logger.debug('Digital location updated successfully', {
          id: location.id
        });
        return { data: location };
      } catch (error) {
        logger.error('Failed to update digital location', { error, id, input });
        throw error;
      }
    },
    {
      onSuccess: () => {
        queryClient.invalidateQueries({ queryKey: digitalLocationsQueryKey });
        toast.success('Digital service updated', {
          description: 'Your digital service has been successfully updated.',
          duration: 3000
        });
      },
      onError: (error: AxiosError<ApiError>) => {
        let errorMessage = 'Failed to update service';
        let errorDescription = "We couldn't update your digital service. Please try again.";

        if (error.response) {
          const status = error.response.status;
          if (status === 401 || status === 403) {
            errorMessage = 'Permission denied';
            errorDescription = "You don't have permission to update this service.";
          } else if (status === 404) {
            errorMessage = 'Service not found';
            errorDescription = "The service you're trying to update doesn't exist.";
          } else if (status >= 500) {
            errorMessage = 'Server error';
            errorDescription = "The server encountered an error. Please try again later.";
          }
        } else if (error.request) {
          errorMessage = 'Network error';
          errorDescription = "Couldn't connect to the server. Please check your internet connection.";
        }

        toast.error(errorMessage, {
          description: errorDescription,
          duration: 5000
        });
      }
    }
  );
}

/**
 * Hook for deleting a digital location
 */
export function useDeleteDigitalLocation() {
  const queryClient = useQueryClient();

  return useBackendMutation<void, string>(
    async (id: string) => {
      logger.debug('Deleting digital location', { id });
      try {
        await digitalLocationService.deleteLocation(id);
        logger.debug('Digital location deleted successfully', { id });
        return { data: undefined };
      } catch (error) {
        logger.error('Failed to delete digital location', { error, id });
        throw error;
      }
    },
    {
      onSuccess: () => {
        queryClient.invalidateQueries({ queryKey: digitalLocationsQueryKey });
        toast.success('Digital service deleted', {
          description: 'Your digital service has been successfully deleted.',
          duration: 3000
        });
      },
      onError: (error: AxiosError<ApiError>) => {
        let errorMessage = 'Failed to delete service';
        let errorDescription = "We couldn't delete your digital service. Please try again.";

        if (error.response) {
          const status = error.response.status;
          if (status === 401 || status === 403) {
            errorMessage = 'Permission denied';
            errorDescription = "You don't have permission to delete this service.";
          } else if (status === 404) {
            errorMessage = 'Service not found';
            errorDescription = "The service you're trying to delete doesn't exist.";
          } else if (status >= 500) {
            errorMessage = 'Server error';
            errorDescription = "The server encountered an error. Please try again later.";
          }
        } else if (error.request) {
          errorMessage = 'Network error';
          errorDescription = "Couldn't connect to the server. Please check your internet connection.";
        }

        toast.error(errorMessage, {
          description: errorDescription,
          duration: 5000
        });
      }
    }
  );
}