/**
 * Online Services API
 *
 * For API standards and best practices, see:
 * @see {@link ../../../docs/api-standards.md}
 */

import { axiosInstance } from '@/core/api/client/axios-instance';
import { logger } from '@/core/utils/logger/logger';

export interface CreateOnlineServiceRequest {
  name: string;
  isActive: boolean;
  url: string;
  billingCycle: string;
  costPerCycle: number;
  nextPaymentDate: string;
  paymentMethod: string;
}

export interface OnlineServiceResponse {
  id: string;
  name: string;
  isActive: boolean;
  url: string;
  billingCycle: string;
  costPerCycle: number;
  nextPaymentDate: string;
  paymentMethod: string;
  createdAt: string;
  updatedAt: string;
}

/**
 * Creates a new online service
 *
 * @async
 * @function createOnlineService
 * @param {CreateOnlineServiceRequest} serviceData - The data for creating a new online service
 * @returns {Promise<OnlineServiceResponse>} A promise that resolves to the created online service
 *
 * @throws {Error} If the API request fails or returns an invalid response
 */
export const createOnlineService = async (serviceData: CreateOnlineServiceRequest): Promise<OnlineServiceResponse> => {
  logger.debug('📝 Creating new online service', { serviceData });

  try {
    const response = await axiosInstance.post<OnlineServiceResponse>(
      '/v1/locations/digital',
      serviceData
    );

    logger.debug('🔄 Full backend response:', { response });

    if (!response.data || typeof response.data !== 'object') {
      throw new Error(`Invalid response: Expected an object, got ${JSON.stringify(response.data)}`);
    }

    return response.data;
  } catch (error) {
    logger.error('🚨 Service creation failed:', { error });
    throw error;
  }
};

/**
 * Updates an existing online service
 *
 * @async
 * @function updateOnlineService
 * @param {string} id - The ID of the service to update
 * @param {CreateOnlineServiceRequest} serviceData - The data for updating the online service
 * @returns {Promise<OnlineServiceResponse>} A promise that resolves to the updated online service
 */
export const updateOnlineService = async (id: string, serviceData: CreateOnlineServiceRequest): Promise<OnlineServiceResponse> => {
  logger.debug('📝 Updating online service', { id, serviceData });

  try {
    const response = await axiosInstance.put<OnlineServiceResponse>(
      `/v1/locations/digital/${id}`,
      serviceData
    );

    logger.debug('🔄 Full backend response:', { response });

    if (!response.data || typeof response.data !== 'object') {
      throw new Error(`Invalid response: Expected an object, got ${JSON.stringify(response.data)}`);
    }

    return response.data;
  } catch (error) {
    logger.error('🚨 Service update failed:', { error });
    throw error;
  }
};

/**
 * Deletes an online service
 *
 * @async
 * @function deleteOnlineService
 * @param {string} serviceId - The ID of the service to delete
 * @returns {Promise<void>}
 */
export const deleteOnlineService = async (serviceId: string): Promise<void> => {
  logger.debug('🗑️ Deleting online service', { serviceId });

  try {
    await axiosInstance.delete(`/v1/locations/digital/${serviceId}`);
    logger.debug('✅ Service deleted successfully');
  } catch (error) {
    logger.error('🚨 Service deletion failed:', { error });
    throw error;
  }
};
