/**
 * Digital Services Catalog Service
 *
 * Provides functions for fetching the digital services catalog from the backend.
 * This catalog contains all available digital services that users can add to their locations.
 */

import { axiosInstance } from '../client/axios-instance';
import { logger } from '@/core/utils/logger/logger';
import type { DigitalServiceItem } from '@/types/services';

interface DigitalServicesCatalogResponse {
  success: boolean;
  data: {
    catalog: DigitalServiceItem[];
  };
  metadata: {
    timestamp: string;
    request_id: string;
  };
}

/**
 * Fetches the digital services catalog from the backend.
 *
 * @returns Promise resolving to an array of digital service items
 * @throws Error if the request fails or the response is invalid
 */
export async function getDigitalServicesCatalog(): Promise<DigitalServiceItem[]> {
  try {
    logger.debug('Fetching digital services catalog');

    const response = await axiosInstance.get<DigitalServicesCatalogResponse>(
      '/v1/services/catalog',
      {
        headers: {
          'Content-Type': 'application/json',
          'Accept': 'application/json'
        }
      }
    );

    if (!response.data?.success) {
      throw new Error('Failed to fetch digital services catalog');
    }

    logger.debug('Digital services catalog fetched successfully', {
      itemCount: response.data.data.catalog.length
    });

    return response.data.data.catalog;
  } catch (error) {
    logger.error('Error fetching digital services catalog', { error });
    throw error;
  }
}