/**
 * Physical Location Query Hooks
 *
 * Provides React Query hooks for fetching and managing physical media storage locations and their sublocations.
 */

// React Query
import { useMutation, useQueryClient } from '@tanstack/react-query';

// Base Query Hook
import { useAPIQuery } from '@/core/api/queries/useAPIQuery';

// Service Layer methods
import {
  getPhysicalLocationsBFFResponse,
  getSinglePhysicalLocation,
  createPhysicalLocation,
  updatePhysicalLocation,
  deletePhysicalLocation,
  getAllSublocations,
  getSingleSublocation,
  createSublocation,
  updateSublocation,
  deleteSublocation,
} from '@/core/api/services/physicalLocation.service';

// Utils
import { showToast } from '@/shared/components/ui/TanstackMutationToast/showToast';

// Types
import type { PhysicalLocation, CreatePhysicalLocationRequest, LocationsBFFResponse } from '@/types/domain/physical-location';
import type { Sublocation, CreateSublocationRequest } from '@/types/domain/sublocation';
import type { DeletePhysicalLocationResponse, DeleteSublocationResponse } from '@/core/api/services/physicalLocation.service';

// Constants
import { TOAST_SUCCESS_MESSAGES } from '@/shared/constants/toast.success.messages';
import { TOAST_DURATIONS } from '@/shared/constants/toast.config';
import { TOAST_ERROR_MESSAGES } from '@/shared/constants/toast.error.messages';
import { logger } from '@/core/utils/logger/logger';
import { AxiosError } from 'axios';
import type { ApiError } from '@/types/api/response';

/**
 * Query key factory for physical location queries
 */
export const physicalLocationKeys = {
  all: ['physical-locations'] as const,
  lists: () => [...physicalLocationKeys.all, 'list'] as const,
  list: (filters: string) => [...physicalLocationKeys.lists(), { filters }] as const,
  details: () => [...physicalLocationKeys.all, 'detail'] as const,
  detail: (id: string) => [...physicalLocationKeys.details(), id] as const,
  sublocations: (parentId: string) => [...physicalLocationKeys.detail(parentId), 'sublocations'] as const,
};

/**
 * Hook to fetch all physical locations
 */
export const useGetPhysicalLocationsBFFResponse = () => {
  return useAPIQuery<LocationsBFFResponse>({
    queryKey: physicalLocationKeys.lists(),
    queryFn: async () => {
      try {
        //console.log('[DEBUG] useGetPhysicalLocationsBFFResponse: Starting query function');
        const response = await getPhysicalLocationsBFFResponse();
        console.log('[DEBUG] useGetPhysicalLocationsBFFResponse: Raw response:', response);

        if (!response) {
          //console.error('[DEBUG] useGetPhysicalLocationsBFFResponse: Response is null or undefined');
          throw new Error('No data received from server');
        }

        if (!response.physicalLocations || !response.sublocations) {
          // console.error('[DEBUG] useGetPhysicalLocationsBFFResponse: Invalid response structure:', {
          //   hasPhysicalLocations: !!response.physicalLocations,
          //   hasSublocations: !!response.sublocations,
          //   response
          // });
          throw new Error('Invalid response structure from server');
        }

        // console.log('[DEBUG] useGetPhysicalLocationsBFFResponse: Successfully parsed response:', {
        //   physicalLocationsCount: response.physicalLocations.length,
        //   sublocationsCount: response.sublocations.length
        // });

        return response;
      } catch (error) {
        console.error('[DEBUG] useGetPhysicalLocationsBFFResponse: Error fetching data:', error);
        throw error;
      }
    },
    staleTime: 30000,
    refetchOnMount: false,
  });
};

/**
 * Hook to fetch a single physical location by ID
 */
export const useGetSinglePhysicalLocation = (id: string) => {
  return useAPIQuery<PhysicalLocation>({
    queryKey: physicalLocationKeys.detail(id),
    queryFn: async () => {
      const singleLocation = await getSinglePhysicalLocation(id);
      return singleLocation;
    },
    enabled: !!id,
  });
};

/**
 * Hook to create a new physical location
 */
export const useCreatePhysicalLocation = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreatePhysicalLocationRequest) => createPhysicalLocation(data),
    onSuccess: () => {
      // Invalidate queries that need to be refreshed
      queryClient.invalidateQueries({ queryKey: physicalLocationKeys.lists() });
      queryClient.invalidateQueries({ queryKey: ['analytics'] });

      showToast({
        message: TOAST_SUCCESS_MESSAGES.PHYSICAL_LOCATION.CREATE,
        variant: 'success',
        duration: TOAST_DURATIONS.EXTENDED,
      });
    },
    onError: (error) => {
      console.error('Failed to create physical location:', error);
      showToast({
        message: TOAST_ERROR_MESSAGES.PHYSICAL_LOCATION.CREATE.DEFAULT,
        variant: 'error',
        duration: TOAST_DURATIONS.EXTENDED,
      });
    }
  });
};

/**
 * Hook to update an existing physical location
 */
export const useUpdatePhysicalLocation = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({ id, data }: { id: string; data: Partial<CreatePhysicalLocationRequest> }) => {
      logger.debug('Updating physical location', { id, data });
      try {
        const location = await updatePhysicalLocation(id, data);
        logger.debug('Physical location updated successfully', {
          id: location.id
        });
        return location;
      } catch (error) {
        logger.error('Failed to update physical location', { error, id, data });
        throw error;
      }
    },
    onSuccess: (data) => {
      // Invalidate all relevant queries
      queryClient.invalidateQueries({ queryKey: physicalLocationKeys.detail(data.id) });
      queryClient.invalidateQueries({ queryKey: physicalLocationKeys.lists() });
      queryClient.invalidateQueries({ queryKey: physicalLocationKeys.all }); // This will invalidate BFF queries
      queryClient.invalidateQueries({ queryKey: ['analytics'] });

      // Show success toast with the correct message
      showToast({
        message: TOAST_SUCCESS_MESSAGES.PHYSICAL_LOCATION.UPDATE,
        variant: 'success',
        duration: TOAST_DURATIONS.EXTENDED,
      });
    },
    onError: (error: AxiosError<ApiError>) => {
      const errorMessages = TOAST_ERROR_MESSAGES.PHYSICAL_LOCATION.UPDATE as {
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
 * Hook to delete a physical location
 */
export const useDeletePhysicalLocation = () => {
  const queryClient = useQueryClient();

  return useMutation<DeletePhysicalLocationResponse['physical'], Error, string | string[]>({
    mutationFn: (ids: string | string[]) => deletePhysicalLocation(ids),
    onSuccess: async (response) => {
      // Log the deletion results
      logger.debug('Successfully deleted locations:', {
        count: response.deleted_count,
        ids: response.location_ids
      });

      // First, invalidate all physical location queries
      await queryClient.invalidateQueries({ queryKey: physicalLocationKeys.all });

      // Then, force a refetch of the BFF data
      await queryClient.refetchQueries({
        queryKey: physicalLocationKeys.lists(),
        exact: true,
        type: 'active'
      });

      showToast({
        message: TOAST_SUCCESS_MESSAGES.PHYSICAL_LOCATION.DELETE,
        variant: 'success',
        duration: TOAST_DURATIONS.EXTENDED,
      });
    },
    onError: (error) => {
      logger.error('Failed to delete physical location:', error);
      showToast({
        message: TOAST_ERROR_MESSAGES.PHYSICAL_LOCATION.DELETE.DEFAULT,
        variant: 'error',
        duration: TOAST_DURATIONS.EXTENDED,
      });
    }
  });
};

// Sublocation Queries

/**
 * Hook to fetch all sublocations for a physical location
 */
export const useGetAllSublocations = (physicalLocationId: string) => {
  return useAPIQuery<Sublocation[]>({
    queryKey: physicalLocationKeys.sublocations(physicalLocationId),
    queryFn: async () => {
      const sublocations = await getAllSublocations(physicalLocationId);
      return sublocations;
    },
    enabled: !!physicalLocationId,
  });
};

/**
 * Hook to fetch a single sublocation by ID
 */
export const useGetSingleSublocation = (id: string) => {
  return useAPIQuery<Sublocation>({
    queryKey: ['sublocations', id],
    queryFn: async () => {
      const sublocation = await getSingleSublocation(id);
      return sublocation;
    },
    enabled: !!id,
  });
};

/**
 * Hook to create a new sublocation
 */
export const useCreateSublocation = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateSublocationRequest) => createSublocation(data),
    onSuccess: (data) => {
      // Invalidate relevant queries
      queryClient.invalidateQueries({ queryKey: physicalLocationKeys.sublocations(data.parentLocationId) });
      queryClient.invalidateQueries({ queryKey: physicalLocationKeys.detail(data.parentLocationId) });
      queryClient.invalidateQueries({ queryKey: ['analytics'] });

      showToast({
        message: TOAST_SUCCESS_MESSAGES.SUBLOCATION.CREATE,
        variant: 'success',
        duration: TOAST_DURATIONS.EXTENDED,
      });
    },
    onError: (error) => {
      console.error('Failed to create sublocation:', error);
      showToast({
        message: TOAST_ERROR_MESSAGES.SUBLOCATION.CREATE.DEFAULT,
        variant: 'error',
        duration: TOAST_DURATIONS.EXTENDED,
      });
    }
  });
};

/**
 * Hook to update an existing sublocation
 */
export const useUpdateSublocation = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({ id, data }: { id: string; data: Partial<CreateSublocationRequest> }) => {
      logger.debug('Updating sublocation', { id, data });
      try {
        const sublocation = await updateSublocation(id, data);
        logger.debug('Sublocation updated successfully', {
          id: sublocation.id
        });
        return sublocation;
      } catch (error) {
        logger.error('Failed to update sublocation', { error, id, data });
        throw error;
      }
    },
    onSuccess: (data) => {
      // Invalidate relevant queries
      queryClient.invalidateQueries({ queryKey: ['sublocations', data.id] });
      queryClient.invalidateQueries({ queryKey: physicalLocationKeys.sublocations(data.parentLocationId) });
      queryClient.invalidateQueries({ queryKey: physicalLocationKeys.detail(data.parentLocationId) });
      queryClient.invalidateQueries({ queryKey: ['analytics'] });

      showToast({
        message: TOAST_SUCCESS_MESSAGES.SUBLOCATION.UPDATE,
        variant: 'success',
        duration: TOAST_DURATIONS.EXTENDED,
      });
    },
    onError: (error: AxiosError<ApiError>) => {
      const errorMessages = TOAST_ERROR_MESSAGES.SUBLOCATION.UPDATE as {
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
 * Hook to delete a sublocation
 */
export const useDeleteSublocation = () => {
  const queryClient = useQueryClient();

  return useMutation<DeleteSublocationResponse['sublocation'], Error, string | string[]>({
    mutationFn: (ids: string | string[]) => deleteSublocation(ids),
    onSuccess: (response) => {
      // Log the deletion results
      // console.log('Successfully deleted sublocations:', {
      //   count: response.deleted_count,
      //   ids: response.sublocation_ids
      // });

      // Note: We can't invalidate specific queries here since we don't know the parent location
      // Instead, we'll invalidate all physical location queries
      queryClient.invalidateQueries({ queryKey: physicalLocationKeys.all });
      queryClient.invalidateQueries({ queryKey: ['analytics'] });

      showToast({
        message: TOAST_SUCCESS_MESSAGES.SUBLOCATION.DELETE,
        variant: 'success',
        duration: TOAST_DURATIONS.EXTENDED,
      })
    },
    onError: (error) => {
      console.error('Failed to delete sublocation:', error);
      showToast({
        message: TOAST_ERROR_MESSAGES.SUBLOCATION.DELETE.DEFAULT,
        variant: 'error',
        duration: TOAST_DURATIONS.EXTENDED,
      })
    }
  });
};