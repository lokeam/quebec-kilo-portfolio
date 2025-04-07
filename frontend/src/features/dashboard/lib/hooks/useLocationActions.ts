import { useCallback } from 'react';
import { toast } from 'sonner';
import { logger } from '@/core/utils/logger/logger';
import { useUpdateLocationMutation, useDeleteLocationMutation } from '@/core/api/queries/useLocationMutations';
import type { PhysicalLocation } from '@/features/dashboard/lib/types/media-storage/physical';

type LocationUpdateData = {
  id: string;
  name: string;
  locationType: string;
  mapCoordinates?: string;
};

/**
 * Hook for managing physical location update actions with logging
 */
export function useLocationUpdate(onSuccess?: () => void) {
  const updateMutation = useUpdateLocationMutation();

  const updateLocation = useCallback((data: LocationUpdateData) => {
    logger.debug('📝 Starting location update:', { data });

    updateMutation.mutate(data as Partial<PhysicalLocation>, {
      onMutate: (variables) => {
        logger.debug('⚡ Update mutation started with variables:', variables);
      },
      onSuccess: (response) => {
        logger.debug('✅ Update mutation succeeded:', response);
        toast.success('Location updated successfully');
        onSuccess?.();
      },
      onError: (error) => {
        logger.error('❌ Update mutation failed:', error);
        toast.error(`Failed to update location: ${error.message}`);
      }
    });
  }, [updateMutation, onSuccess]);

  return {
    updateLocation,
    isUpdating: updateMutation.isPending
  };
}

/**
 * Hook for managing physical location deletion actions with logging
 */
export function useLocationDelete(onSuccess?: () => void) {
  const deleteMutation = useDeleteLocationMutation();

  const deleteLocation = useCallback((locationId: string) => {
    logger.debug('🗑️ Starting location deletion:', { locationId });

    deleteMutation.mutate(locationId, {
      onMutate: (id) => {
        logger.debug('⚡ Delete mutation started for ID:', id);
      },
      onSuccess: (response) => {
        logger.debug('✅ Delete mutation succeeded:', response);
        toast.success('Location deleted successfully');
        onSuccess?.();
      },
      onError: (error) => {
        logger.error('❌ Delete mutation failed:', error);
        toast.error(`Failed to delete location: ${error.message}`);
      }
    });
  }, [deleteMutation, onSuccess]);

  return {
    deleteLocation,
    isDeleting: deleteMutation.isPending
  };
}