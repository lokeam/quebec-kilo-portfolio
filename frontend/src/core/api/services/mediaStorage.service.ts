/**
 * Media Storage API Service
 *
 * This service handles API requests related to digital storage locations
 * It interfaces with the backend GET /api/v1/locations/digital endpoint
 */

import { axiosInstance } from '@/core/api/client/axios-instance';
import type { DigitalLocation } from '@/features/dashboard/lib/types/media-storage/digital-location.types';
import { logger } from '@/core/utils/logger/logger';

interface DigitalLocationsResponse {
  success: boolean;
  user_id: string;
  locations: DigitalLocation[];
}

/**
 * Fetches all digital locations for the current user
 *
 * @async
 * @function getUserDigitalLocations
 * @param {string} [token] - Optional auth token
 * @returns {Promise<DigitalLocation[]>} A promise that resolves to an array of digital locations
 *
 * @throws {Error} If the API request fails
 */
export const getUserDigitalLocations = async (token?: string): Promise<DigitalLocation[]> => {
  logger.debug('Fetching user digital locations');

  try {
    const config = {
      headers: token ? {
        Authorization: `Bearer ${token}`,
        'Content-Type': 'application/json',
        'Accept': 'application/json'
      } : undefined
    };

    const response = await axiosInstance.get<DigitalLocationsResponse>(
      '/v1/locations/digital',
      config
    );

    logger.debug('Digital locations fetched successfully', {
      count: response.locations?.length || 0
    });

    if (!response || !response.success) {
      throw new Error('Invalid response from digital locations API');
    }

    return response.locations || [];
  } catch (error) {
    logger.error('Failed to fetch digital locations', { error });
    throw error;
  }
};

/**
 * Fetches a specific digital location by ID
 *
 * @async
 * @function getDigitalLocationById
 * @param {string} locationId - The ID of the digital location to fetch
 * @param {string} [token] - Optional auth token
 * @returns {Promise<DigitalLocation>} A promise that resolves to the digital location
 *
 * @throws {Error} If the API request fails or the location isn't found
 */
export const getDigitalLocationById = async (locationId: string, token?: string): Promise<DigitalLocation> => {
  logger.debug('Fetching digital location by ID', { locationId });

  try {
    const config = {
      headers: token ? {
        Authorization: `Bearer ${token}`,
        'Content-Type': 'application/json',
        'Accept': 'application/json'
      } : undefined
    };

    const response = await axiosInstance.get<{ success: boolean; location: DigitalLocation }>(
      `/v1/locations/digital/${locationId}`,
      config
    );

    logger.debug('Digital location fetched successfully', { locationId });

    if (!response || !response.success) {
      throw new Error(`Failed to fetch digital location with ID: ${locationId}`);
    }

    return response.location;
  } catch (error) {
    logger.error('Failed to fetch digital location', { locationId, error });
    throw error;
  }
};