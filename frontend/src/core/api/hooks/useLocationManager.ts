import { useCallback } from 'react';

// API Hooks
import {
  useCreateLocationMutation,
  useUpdateLocationMutation,
  useDeleteLocationMutation
} from '@/core/api/queries/useLocationMutations';

import {
  useCreateSublocationMutation,
  useUpdateSublocationMutation,
  useDeleteSublocationMutation
} from '@/core/api/queries/useSublocationMutations';

// Guards
//import { isPhysicalLocation } from '@/features/dashboard/lib/types/media-storage/guards';

// Types
import type { PhysicalLocation } from '@/features/dashboard/lib/types/media-storage/physical';
import type { Sublocation } from '@/features/dashboard/lib/types/media-storage/sublocation';
import type { LocationType } from '@/features/dashboard/lib/types/media-storage/constants';

// Import response types
import type { LocationResponse } from '@/core/api/queries/useLocationMutations';
import type { SublocationResponse } from '@/core/api/queries/useSublocationMutations';

// Utils
import { toast } from 'sonner';
import { logger } from '@/core/utils/logger/logger';

// Types
type LocationBaseData = {
  id?: string;
  name: string;
  [key: string]: unknown;
};

type PhysicalLocationData = LocationBaseData & {
  locationType: string;
  mapCoordinates?: string;
};

type SublocationData = LocationBaseData & {
  type: string;
  parentLocationId: string;
};

// Specific hook exports for backward compatibility
export function useLocationUpdate(onSuccess?: () => void) {
  const updateMutation = useUpdateLocationMutation();

  const updateLocation = useCallback((data: PhysicalLocationData) => {
    logger.debug('📝 Starting location update:', { data });

    updateMutation.mutate(data as Partial<PhysicalLocation>, {
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

export function useLocationDelete(onSuccess?: () => void) {
  const deleteMutation = useDeleteLocationMutation();

  const deleteLocation = useCallback((locationId: string) => {
    logger.debug('🗑️ Starting location deletion:', { locationId });

    deleteMutation.mutate(locationId, {
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

// New unified manager hook
export interface LocationManagerConfig {
  type: LocationType;
  onSuccess?: () => void;
}

/**
 * A unified hook for managing both physical locations and sublocations
 * with consistent interface for CRUD operations
 */
export function useLocationManager({ type, onSuccess }: LocationManagerConfig) {
  // Always initialize all hooks (React rules)
  const createLocationMutation = useCreateLocationMutation();
  const updateLocationMutation = useUpdateLocationMutation();
  const deleteLocationMutation = useDeleteLocationMutation();

  const createSublocationMutation = useCreateSublocationMutation();
  const updateSublocationMutation = useUpdateSublocationMutation();
  const deleteSublocationMutation = useDeleteSublocationMutation();

  // Get entity name for messages
  const entityName = type === 'physical' ? 'location' : 'sublocation';

  // CREATE
  const create = useCallback((data: PhysicalLocationData | SublocationData) => {
    // Add direct console log to bypass logger
    console.log('DIRECT LOG: About to create location/sublocation', data);

    logger.debug(`📝 Starting ${entityName} creation:`, { data });

    if (type === 'physical') {
      createLocationMutation.mutate(data as Partial<PhysicalLocation>, {
        onSuccess: (response: { data: LocationResponse }) => {
          logger.debug(`✅ ${entityName} created successfully:`, response);
          toast.success(`${entityName} created successfully`);
          onSuccess?.();
        },
        onError: (error: Error) => {
          logger.error(`❌ Failed to create ${entityName}:`, error);
          toast.error(`Failed to create ${entityName}: ${error.message}`);
        }
      });
    } else {
      createSublocationMutation.mutate(data as Partial<Sublocation>, {
        onSuccess: (response: { data: SublocationResponse }) => {
          logger.debug(`✅ ${entityName} created successfully:`, response);
          toast.success(`${entityName} created successfully`);
          onSuccess?.();
        },
        onError: (error: Error) => {
          logger.error(`❌ Failed to create ${entityName}:`, error);
          toast.error(`Failed to create ${entityName}: ${error.message}`);
        }
      });
    }
  }, [type, createLocationMutation, createSublocationMutation, entityName, onSuccess]);

  // UPDATE
  const update = useCallback((data: PhysicalLocationData | SublocationData) => {
    logger.debug(`📝 Starting ${entityName} update:`, { data });

    if (type === 'physical') {
      updateLocationMutation.mutate(data as Partial<PhysicalLocation>, {
        onSuccess: (response) => {
          logger.debug(`✅ ${entityName} updated successfully:`, response);
          toast.success(`${entityName} updated successfully`);
          onSuccess?.();
        },
        onError: (error) => {
          logger.error(`❌ Failed to update ${entityName}:`, error);
          toast.error(`Failed to update ${entityName}: ${error.message}`);
        }
      });
    } else {
      updateSublocationMutation.mutate(data as Partial<Sublocation>, {
        onSuccess: (response) => {
          logger.debug(`✅ ${entityName} updated successfully:`, response);
          toast.success(`${entityName} updated successfully`);
          onSuccess?.();
        },
        onError: (error) => {
          logger.error(`❌ Failed to update ${entityName}:`, error);
          toast.error(`Failed to update ${entityName}: ${error.message}`);
        }
      });
    }
  }, [type, updateLocationMutation, updateSublocationMutation, entityName, onSuccess]);

  // DELETE
  const deleteItem = useCallback((itemId: string) => {
    logger.debug(`🗑️ Starting ${entityName} deletion:`, { itemId });

    if (type === 'physical') {
      deleteLocationMutation.mutate(itemId, {
        onSuccess: (response) => {
          logger.debug(`✅ ${entityName} deleted successfully:`, response);
          toast.success(`${entityName} deleted successfully`);
          onSuccess?.();
        },
        onError: (error) => {
          logger.error(`❌ Failed to delete ${entityName}:`, error);
          toast.error(`Failed to delete ${entityName}: ${error.message}`);
        }
      });
    } else {
      deleteSublocationMutation.mutate(itemId, {
        onSuccess: (response) => {
          logger.debug(`✅ ${entityName} deleted successfully:`, response);
          toast.success(`${entityName} deleted successfully`);
          onSuccess?.();
        },
        onError: (error) => {
          logger.error(`❌ Failed to delete ${entityName}:`, error);
          toast.error(`Failed to delete ${entityName}: ${error.message}`);
        }
      });
    }
  }, [type, deleteLocationMutation, deleteSublocationMutation, entityName, onSuccess]);

  // Status indicators
  const isCreating = type === 'physical'
    ? createLocationMutation.isPending
    : createSublocationMutation.isPending;

  const isUpdating = type === 'physical'
    ? updateLocationMutation.isPending
    : updateSublocationMutation.isPending;

  const isDeleting = type === 'physical'
    ? deleteLocationMutation.isPending
    : deleteSublocationMutation.isPending;

  return {
    create,
    update,
    delete: deleteItem,
    isCreating,
    isUpdating,
    isDeleting
  };
}