/**
 * Digital Media Storage Query Hooks
 *
 * Provides React Query hooks for fetching and managing digital media storage locations.
 */

import { useMutation, useQueryClient } from '@tanstack/react-query';
import { useAPIQuery } from './useAPIQuery';
import {
  getAllDigitalLocationsBFFResponse,
  getSingleDigitalLocation,
  createDigitalLocation,
  updateDigitalLocation,
  deleteDigitalLocation,
} from '@/core/api/services/digitalLocation.service';
import type { DigitalLocation, CreateDigitalLocationRequest } from '@/types/domain/digital-location';
import { adaptDigitalLocationToService } from '../adapters/digitalLocation.adapter';
import { logger } from '@/core/utils/logger/logger';
import { AxiosError } from 'axios';
import type { ApiError } from '@/types/api/';
import { mediaStorageKeys } from '../constants/query-keys/mediaStorage';
import { showToast } from '@/shared/components/ui/TanstackMutationToast/showToast';

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
 * Hook to fetch all digital locations -- NOTE: CHANGE THIS TO BFF
 */
export const useDigitalLocations = () => {
  return useAPIQuery<DigitalLocation[]>({
    queryKey: digitalLocationKeys.lists(),
    queryFn: () => getAllDigitalLocations(),
  });
};

/**
 * Hook to fetch a single digital location by ID
 */
export const useDigitalLocation = (id: string) => {
  return useAPIQuery<DigitalLocation>({
    queryKey: digitalLocationKeys.detail(id),
    queryFn: () => digitalLocationService.getLocationById(id),
    enabled: !!id,
  });
};

/**
 * Hook to create a new digital location
 */
export const useCreateDigitalLocation = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateDigitalLocationRequest) =>
      digitalLocationService.createLocation(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: digitalLocationKeys.lists() });
    },
  });
};

/**
 * Hook to update an existing digital location
 */
export const useUpdateDigitalLocation = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<CreateDigitalLocationRequest> }) =>
      digitalLocationService.updateLocation(id, data),
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
    mutationFn: (id: string) => digitalLocationService.deleteLocation(id),
    onSuccess: (_, id) => {
      queryClient.invalidateQueries({ queryKey: digitalLocationKeys.detail(id) });
      queryClient.invalidateQueries({ queryKey: digitalLocationKeys.lists() });
    },
  });
};

/**
 * Hook to fetch all digital locations as online services
 */
export const useOnlineServices = () => {
  return useAPIQuery<DigitalLocation[]>({
    queryKey: digitalLocationKeys.lists(),
    queryFn: async () => {
      const locations = await digitalLocationService.getAllLocations();
      return locations.map(adaptDigitalLocationToService);
    },
  });
};

/**
 * Hook to fetch a single online service by ID
 */
export const useOnlineService = (id: string) => {
  return useAPIQuery<DigitalLocation>({
    queryKey: digitalLocationKeys.detail(id),
    queryFn: async () => {
      const location = await digitalLocationService.getLocationById(id);
      return adaptDigitalLocationToService(location);
    },
    enabled: !!id,
  });
};

/**
 * Hook for creating a new digital location
 */
export function useCreateDigitalLocationOld() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (input: CreateDigitalLocationRequest) => {
      logger.debug('Creating digital location', { input });
      try {
        const location = await digitalLocationService.createLocation(input);
        logger.debug('Digital location created successfully', {
          id: location.id
        });
        return location;
      } catch (error) {
        logger.error('Failed to create digital location', { error, input });
        throw error;
      }
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: mediaStorageKeys.digitalLocations.all });
      showToast({
        message: 'Digital service created',
        variant: 'success',
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

      showToast({
        message: errorMessage,
        variant: 'error',
        duration: 5000
      });
    }
  });
}

/**
 * Hook for updating a digital location
 */
export function useUpdateDigitalLocationOld() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({ id, input }: { id: string; input: Partial<CreateDigitalLocationRequest> }) => {
      logger.debug('Updating digital location', { id, input });
      try {
        const location = await digitalLocationService.updateLocation(id, input);
        logger.debug('Digital location updated successfully', {
          id: location.id
        });
        return location;
      } catch (error) {
        logger.error('Failed to update digital location', { error, id, input });
        throw error;
      }
    },
    onSuccess: (data) => {
      queryClient.invalidateQueries({ queryKey: mediaStorageKeys.digitalLocations.all });
      queryClient.invalidateQueries({ queryKey: mediaStorageKeys.digitalLocations.byId(data.id) });
      showToast({
        message: 'Digital service updated',
        variant: 'success',
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

      showToast({
        message: errorMessage,
        variant: 'error',
        duration: 5000
      });
    }
  });
}

/**
 * Hook for deleting a digital location
 */
export function useDeleteDigitalLocationOld() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (id: string) => {
      logger.debug('Deleting digital location', { id });
      try {
        await digitalLocationService.deleteLocation(id);
        logger.debug('Digital location deleted successfully', { id });
      } catch (error) {
        logger.error('Failed to delete digital location', { error, id });
        throw error;
      }
    },
    onSuccess: (_, id) => {
      queryClient.invalidateQueries({ queryKey: mediaStorageKeys.digitalLocations.all });
      queryClient.invalidateQueries({ queryKey: mediaStorageKeys.digitalLocations.byId(id) });
      showToast({
        message: 'Digital service deleted',
        variant: 'success',
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

      showToast({
        message: errorMessage,
        variant: 'error',
        duration: 5000
      });
    }
  });
}