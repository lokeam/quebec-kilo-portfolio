/**
 * Media Storage API Service
 *
 * This service handles API requests related to physical storage locations and their sublocations.
 * It provides a unified interface for managing physical media storage operations.
 */

import { axiosInstance } from '@/core/api/client/axios-instance';
import { logger } from '@/core/utils/logger/logger';
import type { PhysicalLocation, PhysicalLocationMetadata } from '@/types/domain/physical-location';
import type { Sublocation, SublocationMetadata } from '@/types/domain/sublocation';
import type { PhysicalLocationType, SublocationType } from '@/types/domain/location-types';

// Standard response type for all API calls
interface MediaStorageResponse<T> {
  success: boolean;
  data: T;
  metadata: {
    timestamp: string;
    request_id: string;
  };
}

// Request types for mutations
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

interface PhysicalLocationsResponse {
  success: boolean;
  userId: string;
  locations: PhysicalLocation[];
}

// Physical Location Operations

/**
 * Fetches all physical locations for the current user
 *
 * @async
 * @function getPhysicalLocations
 * @param {string} [token] - Optional auth token
 * @returns {Promise<PhysicalLocation[]>} A promise that resolves to an array of physical locations
 * @throws {Error} If the API request fails
 */
export const getPhysicalLocations = async (token?: string): Promise<PhysicalLocation[]> => {
  logger.debug('Fetching physical locations');

  try {
    const config = {
      headers: token ? {
        Authorization: `Bearer ${token}`,
        'Content-Type': 'application/json',
        'Accept': 'application/json'
      } : undefined
    };

    const response = await axiosInstance.get<PhysicalLocationsResponse>(
      '/v1/locations/physical',
      config
    );

    logger.debug('Physical locations fetched successfully', {
      count: response.data.locations.length
    });

    return response.data.locations;
  } catch (error) {
    logger.error('Failed to fetch physical locations', { error });
    throw error;
  }
};

/**
 * Fetches a specific physical location by ID
 *
 * @async
 * @function getPhysicalLocationById
 * @param {string} id - The ID of the physical location to fetch
 * @param {string} [token] - Optional auth token
 * @returns {Promise<PhysicalLocation>} A promise that resolves to the physical location
 * @throws {Error} If the API request fails or the location isn't found
 */
export const getPhysicalLocationById = async (id: string, token?: string): Promise<PhysicalLocation> => {
  logger.debug('Fetching physical location by ID', { id });

  try {
    const config = {
      headers: token ? {
        Authorization: `Bearer ${token}`,
        'Content-Type': 'application/json',
        'Accept': 'application/json'
      } : undefined
    };

    const response = await axiosInstance.get<MediaStorageResponse<PhysicalLocation>>(
      `/v1/locations/physical/${id}`,
      config
    );

    logger.debug('Physical location fetched successfully', { id });
    return response.data.data;
  } catch (error) {
    logger.error('Failed to fetch physical location', { id, error });
    throw error;
  }
};

/**
 * Creates a new physical location
 *
 * @async
 * @function createPhysicalLocation
 * @param {CreatePhysicalLocationRequest} data - The data for creating a new physical location
 * @param {string} [token] - Optional auth token
 * @returns {Promise<PhysicalLocation>} A promise that resolves to the created physical location
 * @throws {Error} If the API request fails
 */
export const createPhysicalLocation = async (
  data: CreatePhysicalLocationRequest,
  token?: string
): Promise<PhysicalLocation> => {
  logger.debug('Creating physical location', { data });

  try {
    const config = {
      headers: token ? {
        Authorization: `Bearer ${token}`,
        'Content-Type': 'application/json',
        'Accept': 'application/json'
      } : undefined
    };

    const response = await axiosInstance.post<MediaStorageResponse<PhysicalLocation>>(
      '/v1/locations/physical',
      data,
      config
    );

    logger.debug('Physical location created successfully', { id: response.data.data.id });
    return response.data.data;
  } catch (error) {
    logger.error('Failed to create physical location', { error });
    throw error;
  }
};

/**
 * Updates an existing physical location
 *
 * @async
 * @function updatePhysicalLocation
 * @param {string} id - The ID of the physical location to update
 * @param {Partial<CreatePhysicalLocationRequest>} data - The data to update
 * @param {string} [token] - Optional auth token
 * @returns {Promise<PhysicalLocation>} A promise that resolves to the updated physical location
 * @throws {Error} If the API request fails
 */
export const updatePhysicalLocation = async (
  id: string,
  data: Partial<CreatePhysicalLocationRequest>,
  token?: string
): Promise<PhysicalLocation> => {
  logger.debug('Updating physical location', { id, data });

  try {
    const config = {
      headers: token ? {
        Authorization: `Bearer ${token}`,
        'Content-Type': 'application/json',
        'Accept': 'application/json'
      } : undefined
    };

    const response = await axiosInstance.put<MediaStorageResponse<PhysicalLocation>>(
      `/v1/locations/physical/${id}`,
      data,
      config
    );

    logger.debug('Physical location updated successfully', { id });
    return response.data.data;
  } catch (error) {
    logger.error('Failed to update physical location', { id, error });
    throw error;
  }
};

/**
 * Deletes a physical location
 *
 * @async
 * @function deletePhysicalLocation
 * @param {string} id - The ID of the physical location to delete
 * @param {string} [token] - Optional auth token
 * @returns {Promise<void>}
 * @throws {Error} If the API request fails
 */
export const deletePhysicalLocation = async (id: string, token?: string): Promise<void> => {
  logger.debug('Deleting physical location', { id });

  try {
    const config = {
      headers: token ? {
        Authorization: `Bearer ${token}`,
        'Content-Type': 'application/json',
        'Accept': 'application/json'
      } : undefined
    };

    await axiosInstance.delete(`/v1/locations/physical/${id}`, config);
    logger.debug('Physical location deleted successfully', { id });
  } catch (error) {
    logger.error('Failed to delete physical location', { id, error });
    throw error;
  }
};

// Sublocation Operations

/**
 * Fetches all sublocations for a physical location
 *
 * @async
 * @function getSublocations
 * @param {string} parentLocationId - The ID of the parent physical location
 * @param {string} [token] - Optional auth token
 * @returns {Promise<Sublocation[]>} A promise that resolves to an array of sublocations
 * @throws {Error} If the API request fails
 */
export const getSublocations = async (
  parentLocationId: string,
  token?: string
): Promise<Sublocation[]> => {
  logger.debug('Fetching sublocations', { parentLocationId });

  try {
    const config = {
      headers: token ? {
        Authorization: `Bearer ${token}`,
        'Content-Type': 'application/json',
        'Accept': 'application/json'
      } : undefined
    };

    const response = await axiosInstance.get<MediaStorageResponse<Sublocation[]>>(
      `/v1/locations/physical/${parentLocationId}/sublocations`,
      config
    );

    logger.debug('Sublocations fetched successfully', {
      parentLocationId,
      count: response.data.data.length
    });

    return response.data.data;
  } catch (error) {
    logger.error('Failed to fetch sublocations', { parentLocationId, error });
    throw error;
  }
};

/**
 * Creates a new sublocation
 *
 * @async
 * @function createSublocation
 * @param {CreateSublocationRequest} data - The data for creating a new sublocation
 * @param {string} [token] - Optional auth token
 * @returns {Promise<Sublocation>} A promise that resolves to the created sublocation
 * @throws {Error} If the API request fails
 */
export const createSublocation = async (
  data: CreateSublocationRequest,
  token?: string
): Promise<Sublocation> => {
  logger.debug('Creating sublocation', { data });

  try {
    const config = {
      headers: token ? {
        Authorization: `Bearer ${token}`,
        'Content-Type': 'application/json',
        'Accept': 'application/json'
      } : undefined
    };

    const response = await axiosInstance.post<MediaStorageResponse<Sublocation>>(
      `/v1/locations/physical/${data.parentLocationId}/sublocations`,
      data,
      config
    );

    logger.debug('Sublocation created successfully', { id: response.data.data.id });
    return response.data.data;
  } catch (error) {
    logger.error('Failed to create sublocation', { error });
    throw error;
  }
};

/**
 * Updates an existing sublocation
 *
 * @async
 * @function updateSublocation
 * @param {string} id - The ID of the sublocation to update
 * @param {Partial<CreateSublocationRequest>} data - The data to update
 * @param {string} [token] - Optional auth token
 * @returns {Promise<Sublocation>} A promise that resolves to the updated sublocation
 * @throws {Error} If the API request fails
 */
export const updateSublocation = async (
  id: string,
  data: Partial<CreateSublocationRequest>,
  token?: string
): Promise<Sublocation> => {
  logger.debug('Updating sublocation', { id, data });

  try {
    const config = {
      headers: token ? {
        Authorization: `Bearer ${token}`,
        'Content-Type': 'application/json',
        'Accept': 'application/json'
      } : undefined
    };

    const response = await axiosInstance.put<MediaStorageResponse<Sublocation>>(
      `/v1/locations/physical/sublocations/${id}`,
      data,
      config
    );

    logger.debug('Sublocation updated successfully', { id });
    return response.data.data;
  } catch (error) {
    logger.error('Failed to update sublocation', { id, error });
    throw error;
  }
};

/**
 * Deletes a sublocation
 *
 * @async
 * @function deleteSublocation
 * @param {string} id - The ID of the sublocation to delete
 * @param {string} [token] - Optional auth token
 * @returns {Promise<void>}
 * @throws {Error} If the API request fails
 */
export const deleteSublocation = async (id: string, token?: string): Promise<void> => {
  logger.debug('Deleting sublocation', { id });

  try {
    const config = {
      headers: token ? {
        Authorization: `Bearer ${token}`,
        'Content-Type': 'application/json',
        'Accept': 'application/json'
      } : undefined
    };

    await axiosInstance.delete(`/v1/locations/physical/sublocations/${id}`, config);
    logger.debug('Sublocation deleted successfully', { id });
  } catch (error) {
    logger.error('Failed to delete sublocation', { id, error });
    throw error;
  }
};
