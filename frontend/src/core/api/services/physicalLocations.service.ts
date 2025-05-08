/**
 * Physical Locations API Service
 *
 * This service handles API requests related to physical storage locations
 * It interfaces with the backend endpoints for physical locations
 */

import { axiosInstance } from '@/core/api/client/axios-instance';
import type { PhysicalLocation } from '@/features/dashboard/lib/types/media-storage/physical';
import { logger } from '@/core/utils/logger/logger';
import { LocationType } from '@/features/dashboard/lib/types/media-storage/constants';
import type { AxiosResponse } from 'axios';

interface PhysicalLocationsResponse {
  success: boolean;
  userId: string;
  locations: PhysicalLocation[];
}

/**
 * Fetches all physical locations for the current user
 *
 * @async
 * @function getUserPhysicalLocations
 * @param {string} [token] - Optional auth token
 * @returns {Promise<PhysicalLocation[]>} A promise that resolves to an array of physical locations
 *
 * @throws {Error} If the API request fails
 */
export const getUserPhysicalLocations = async (token?: string): Promise<PhysicalLocation[]> => {
  logger.debug('Fetching user physical locations');

  try {
    const config = {
      headers: token ? {
        Authorization: `Bearer ${token}`,
        'Content-Type': 'application/json',
        'Accept': 'application/json'
      } : undefined
    };

    const response: AxiosResponse<PhysicalLocationsResponse> = await axiosInstance.get(
      '/v1/locations/physical',
      config
    );

    logger.debug('Physical locations fetched successfully', {
      count: response.data.locations?.length || 0,
      firstLocation: response.data.locations?.[0],
      firstLocationSublocations: response.data.locations?.[0]?.sublocations
    });

    if (!response.data || !response.data.success) {
      throw new Error('Invalid response from physical locations API');
    }

    // Ensure we have an array of locations, defaulting to empty array if undefined
    const locations = response.data.locations || [];

    // Ensure each location has the correct type
    const transformedLocations = locations.map((location: PhysicalLocation) => ({
      ...location,
      type: LocationType.PHYSICAL
    }));

    logger.debug('Transformed locations', {
      count: transformedLocations.length,
      firstLocation: transformedLocations[0],
      firstLocationSublocations: transformedLocations[0]?.sublocations
    });

    return transformedLocations;
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
 * @param {string} locationId - The ID of the physical location to fetch
 * @param {string} [token] - Optional auth token
 * @returns {Promise<PhysicalLocation>} A promise that resolves to the physical location
 *
 * @throws {Error} If the API request fails or the location isn't found
 */
export const getPhysicalLocationById = async (locationId: string, token?: string): Promise<PhysicalLocation> => {
  logger.debug('Fetching physical location by ID', { locationId });

  try {
    const config = {
      headers: token ? {
        Authorization: `Bearer ${token}`,
        'Content-Type': 'application/json',
        'Accept': 'application/json'
      } : undefined
    };

    const response = await axiosInstance.get<{ success: boolean; location: PhysicalLocation }>(
      `/v1/locations/physical/${locationId}`,
      config
    );

    logger.debug('Physical location fetched successfully', { locationId });

    if (!response.data || !response.data.success) {
      throw new Error(`Failed to fetch physical location with ID: ${locationId}`);
    }

    // Ensure the location has the correct type
    return {
      ...response.data.location,
      type: LocationType.PHYSICAL
    };
  } catch (error) {
    logger.error('Failed to fetch physical location', { locationId, error });
    throw error;
  }
};