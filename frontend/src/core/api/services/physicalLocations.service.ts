/**
 * Physical Locations API Service
 *
 * This service handles API requests related to physical storage locations
 * It interfaces with the backend endpoints for physical locations
 */

import { axiosInstance } from '@/core/api/client/axios-instance';
import type { PhysicalLocation } from '@/features/dashboard/lib/types/media-storage/physical';
import { logger } from '@/core/utils/logger/logger';

interface PhysicalLocationsResponse {
  success: boolean;
  user_id: string;
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

    const response = await axiosInstance.get<PhysicalLocationsResponse>(
      '/v1/locations/physical',
      config
    );

    logger.debug('Physical locations fetched successfully', {
      count: response.locations?.length || 0
    });

    if (!response || !response.success) {
      throw new Error('Invalid response from physical locations API');
    }

    return response.locations || [];
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

    if (!response || !response.success) {
      throw new Error(`Failed to fetch physical location with ID: ${locationId}`);
    }

    return response.location;
  } catch (error) {
    logger.error('Failed to fetch physical location', { locationId, error });
    throw error;
  }
};