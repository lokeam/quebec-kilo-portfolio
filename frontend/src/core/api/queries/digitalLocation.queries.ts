/**
 * Digital Location Query Hooks
 *
 * Provides React Query hooks for fetching and managing digital media storage locations.
 */

import { AxiosError } from 'axios';

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
  getDigitalLocationsBFFResponse,
} from '@/core/api/services/digitalLocation.service';

// Adapters
import { adaptBFFResponseToDigitalLocations } from '@/core/api/adapters/digitalLocation.adapter';

// Utils
import { showToast } from '@/shared/components/ui/TanstackMutationToast/showToast';
import { logger } from '@/core/utils/logger/logger';

// Types
import type { ApiError } from '@/types/api/response';
import type {
  DigitalLocation,
  CreateDigitalLocationRequest,
} from '@/types/domain/digital-location';
import type { DeleteDigitalLocationResponse } from '@/core/api/services/digitalLocation.service';

// Constants
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


export const useGetDigitalLocationsBFFResponse = () => {
  return useAPIQuery<DigitalLocation[]>({
    queryKey: digitalLocationKeys.lists(),
    queryFn: async () => {
      try {
        const response = await getDigitalLocationsBFFResponse();
        console.log('[DEBUG] useGetDigitalLocationsBFFResponse: Raw response:', response);

        if (!response) {
          throw new Error('No data received from server');
        }

        // Transform BFF response to DigitalLocation array
        return adaptBFFResponseToDigitalLocations(response);
      } catch (error) {
        console.error('[DEBUG] useGetDigitalLocationsBFFResponse: Error fetching data: ', error);
        throw error;
      }
    },
    staleTime: 3000,
    refetchOnMount: true,
  })
}


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
    mutationFn: async ({ id, data }: { id: string; data: Partial<CreateDigitalLocationRequest> }) => {
      logger.debug('Updating digital location', { id, data });
      try {
        const location = await updateDigitalLocation(id, data);
        logger.debug('Digital location updated successfully', {
          id: location.id
        });
        return location;
      } catch (error) {
        logger.error('Failed to update digital location', { error, id, data });
        throw error;
      }
    },
    onSuccess: (data) => {
      // Invalidate relevant queries
      queryClient.invalidateQueries({ queryKey: digitalLocationKeys.detail(data.id) });
      queryClient.invalidateQueries({ queryKey: digitalLocationKeys.lists() });
      queryClient.invalidateQueries({ queryKey: ['analytics'] });

      // Show success toast
      showToast({
        message: TOAST_SUCCESS_MESSAGES.DIGITAL_LOCATION.UPDATE,
        variant: 'success',
        duration: TOAST_DURATIONS.EXTENDED,
      });
    },
    onError: (error: AxiosError<ApiError>) => {
      const errorMessages = TOAST_ERROR_MESSAGES.DIGITAL_LOCATION.UPDATE as {
        DEFAULT: string;
        PERMISSION: string;
        NOT_FOUND: string;
        SERVER: string;
      };

      let errorMessage = errorMessages.DEFAULT;

      if (error.response) {
        const status = error.response.status;
        if (status === 401 || status === 403) {
          errorMessage = errorMessages.PERMISSION;
        } else if (status === 404) {
          errorMessage = errorMessages.NOT_FOUND;
        } else if (status >= 500) {
          errorMessage = errorMessages.SERVER;
        }
      }

      showToast({
        message: errorMessage,
        variant: 'error',
        duration: TOAST_DURATIONS.EXTENDED,
      });
    }
  });
};

/**
 * Hook to delete a digital location
 */
export const useDeleteDigitalLocation = () => {
  const queryClient = useQueryClient();

  return useMutation<DeleteDigitalLocationResponse['digital'], Error, string | string[]>({
    mutationFn: (ids: string | string[]) => deleteDigitalLocation(ids),
    onSuccess: (response) => {
      // Log the deletion results
      console.log('Successfully deleted locations:', {
        count: response.deleted_count,
        ids: response.location_ids
      });

      queryClient.invalidateQueries({ queryKey: digitalLocationKeys.lists() });
      queryClient.invalidateQueries({ queryKey: ['analytics'] });

      showToast({
        message: TOAST_SUCCESS_MESSAGES.DIGITAL_LOCATION.DELETE,
        variant: 'success',
        duration: TOAST_DURATIONS.EXTENDED,
      })
    },
    onError: (error) => {
      console.error('Failed to delete digital location:', error);
      showToast({
        message: TOAST_ERROR_MESSAGES.DIGITAL_LOCATION.DELETE.DEFAULT,
        variant: 'error',
        duration: TOAST_DURATIONS.EXTENDED,
      })
    }
  });
};
