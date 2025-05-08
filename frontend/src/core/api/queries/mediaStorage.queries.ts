import { useMutation, useQueryClient } from '@tanstack/react-query';
import type { PhysicalLocation as DomainPhysicalLocation } from '@/types/domain/physical-location';
import type { PhysicalLocation as UIPhysicalLocation } from '@/types/api/storage';
import type { PhysicalLocationMetadata } from '@/types/domain/physical-location';
import type { Sublocation as DomainSublocation } from '@/types/domain/sublocation';
import type { Sublocation as UISublocation } from '@/types/api/storage';
import type { SublocationMetadata } from '@/types/domain/sublocation';
import type { PhysicalLocationType, SublocationType } from '@/types/domain/location-types';
import { mediaStorageKeys } from '../constants/query-keys/mediaStorage';
import { useAPIQuery } from './useAPIQuery';
import { adaptDomainToUIPhysicalLocations, adaptDomainToUIPhysicalLocation } from '../adapters/mediaStorage.adapter';
import {
  getPhysicalLocations,
  getPhysicalLocationById,
  createPhysicalLocation,
  updatePhysicalLocation,
  deletePhysicalLocation,
  createSublocation,
  updateSublocation,
  deleteSublocation,
} from '../services/mediaStorage.service.refactor';

// Request types
interface CreatePhysicalLocationRequest {
  name: string;
  type: PhysicalLocationType;
  description?: string;
  metadata?: PhysicalLocationMetadata;
}

interface CreateSublocationRequest {
  name: string;
  parentLocationId: string;
  type: SublocationType;
  description?: string;
  metadata?: SublocationMetadata;
}

interface DeleteSublocationContext {
  parentLocationId: string;
}

// Physical Location Queries

/**
 * Hook for fetching all physical locations
 */
export const usePhysicalLocations = () => {
  return useAPIQuery<UIPhysicalLocation[]>({
    queryKey: mediaStorageKeys.locations.all,
    queryFn: async () => {
      const locations = await getPhysicalLocations();
      return adaptDomainToUIPhysicalLocations(locations);
    },
  });
};

/**
 * Hook for fetching a specific physical location by ID
 */
export const usePhysicalLocation = (id: string) => {
  return useAPIQuery<UIPhysicalLocation>({
    queryKey: mediaStorageKeys.locations.byId(id),
    queryFn: async () => {
      const location = await getPhysicalLocationById(id);
      return adaptDomainToUIPhysicalLocations([location])[0];
    },
    enabled: !!id,
  });
};

/**
 * Hook for creating a new physical location
 */
export const useCreatePhysicalLocation = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreatePhysicalLocationRequest) => createPhysicalLocation(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: mediaStorageKeys.locations.all });
    },
  });
};

/**
 * Hook for updating a physical location
 */
export const useUpdatePhysicalLocation = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<DomainPhysicalLocation> }) =>
      updatePhysicalLocation(id, data),
    onSuccess: (data) => {
      queryClient.invalidateQueries({ queryKey: mediaStorageKeys.locations.all });
      queryClient.invalidateQueries({ queryKey: mediaStorageKeys.locations.byId(data.id) });
    },
  });
};

/**
 * Hook for deleting a physical location
 */
export const useDeletePhysicalLocation = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: deletePhysicalLocation,
    onSuccess: (_, id) => {
      queryClient.invalidateQueries({ queryKey: mediaStorageKeys.locations.all });
      queryClient.invalidateQueries({ queryKey: mediaStorageKeys.locations.byId(id) });
    },
  });
};

// Sublocation Queries

/**
 * Hook for fetching all sublocations for a physical location
 */
export const useSublocations = (parentLocationId: string) => {
  return useAPIQuery<UISublocation[]>({
    queryKey: mediaStorageKeys.sublocations.byLocation(parentLocationId),
    queryFn: async () => {
      const location = await getPhysicalLocationById(parentLocationId);
      const uiLocation = adaptDomainToUIPhysicalLocation(location);
      return uiLocation.sublocations || [];
    },
    enabled: !!parentLocationId,
  });
};

/**
 * Hook for creating a new sublocation
 */
export const useCreateSublocation = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateSublocationRequest) => createSublocation(data),
    onSuccess: (data) => {
      queryClient.invalidateQueries({
        queryKey: mediaStorageKeys.sublocations.byLocation(data.parentLocationId),
      });
      queryClient.invalidateQueries({
        queryKey: mediaStorageKeys.locations.byId(data.parentLocationId),
      });
    },
  });
};

/**
 * Hook for updating a sublocation
 */
export const useUpdateSublocation = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<DomainSublocation> }) =>
      updateSublocation(id, data),
    onSuccess: (data) => {
      // Invalidate both the sublocation list and the parent location
      queryClient.invalidateQueries({
        queryKey: mediaStorageKeys.sublocations.byLocation(data.parentLocationId),
      });
      queryClient.invalidateQueries({
        queryKey: mediaStorageKeys.locations.byId(data.parentLocationId),
      });
    },
  });
};

/**
 * Hook for deleting a sublocation
 */
export const useDeleteSublocation = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => deleteSublocation(id),
    onSuccess: (_, id, context: DeleteSublocationContext) => {
      queryClient.invalidateQueries({
        queryKey: mediaStorageKeys.sublocations.byLocation(context.parentLocationId),
      });
      queryClient.invalidateQueries({
        queryKey: mediaStorageKeys.locations.byId(context.parentLocationId),
      });
    },
  });
};