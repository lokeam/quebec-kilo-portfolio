/**
 * Digital Location Service
 *
 * Provides functions for managing digital locations through the backend API.
 */

import { axiosInstance } from '../client/axios-instance';
import { logger } from '@/core/utils/logger/logger';
import type { DigitalLocation, CreateDigitalLocationRequest } from '@/features/dashboard/lib/types/media-storage/digital-location.types';

interface ApiResponse<T> {
  success: boolean;
  data: T;
  metadata: {
    timestamp: string;
    request_id: string;
  };
}

interface LocationsResponse {
  locations: DigitalLocation[];
}

export const digitalLocationService = {
  /**
   * Fetches all digital locations for the current user
   */
  async getAllLocations(): Promise<DigitalLocation[]> {
    try {
      logger.debug('Fetching all digital locations');

      const response = await axiosInstance.get<LocationsResponse>(
        '/v1/locations/digital',
        {
          headers: {
            'Content-Type': 'application/json',
            'Accept': 'application/json'
          }
        }
      );

      if (!response.data?.locations) {
        throw new Error('Failed to fetch digital locations');
      }

      logger.debug('Digital locations fetched successfully', {
        count: response.data.locations.length
      });

      return response.data.locations;
    } catch (error) {
      logger.error('Error fetching digital locations', { error });
      throw error;
    }
  },

  /**
   * Creates a new digital location
   */
  async createLocation(input: CreateDigitalLocationRequest): Promise<DigitalLocation> {
    try {
      logger.debug('Creating new digital location', { input });

      const response = await axiosInstance.post<ApiResponse<DigitalLocation>>(
        '/v1/locations/digital',
        input,
        {
          headers: {
            'Content-Type': 'application/json',
            'Accept': 'application/json'
          }
        }
      );

      if (!response.data?.success) {
        throw new Error('Failed to create digital location');
      }

      logger.debug('Digital location created successfully', {
        id: response.data.data.id
      });

      return response.data.data;
    } catch (error) {
      logger.error('Error creating digital location', { error, input });
      throw error;
    }
  },

  /**
   * Updates an existing digital location
   */
  async updateLocation(id: string, input: Partial<CreateDigitalLocationRequest>): Promise<DigitalLocation> {
    try {
      logger.debug('Updating digital location', { id, input });

      const response = await axiosInstance.put<ApiResponse<DigitalLocation>>(
        `/v1/locations/digital/${id}`,
        input,
        {
          headers: {
            'Content-Type': 'application/json',
            'Accept': 'application/json'
          }
        }
      );

      if (!response.data?.success) {
        throw new Error('Failed to update digital location');
      }

      logger.debug('Digital location updated successfully', {
        id: response.data.data.id
      });

      return response.data.data;
    } catch (error) {
      logger.error('Error updating digital location', { error, id, input });
      throw error;
    }
  },

  /**
   * Deletes a digital location
   */
  async deleteLocation(id: string): Promise<void> {
    try {
      logger.debug('Deleting digital location', { id });

      const response = await axiosInstance.delete<ApiResponse<void>>(
        `/v1/locations/digital/${id}`,
        {
          headers: {
            'Content-Type': 'application/json',
            'Accept': 'application/json'
          }
        }
      );

      if (!response.data?.success) {
        throw new Error('Failed to delete digital location');
      }

      logger.debug('Digital location deleted successfully', { id });
    } catch (error) {
      logger.error('Error deleting digital location', { error, id });
      throw error;
    }
  }
};