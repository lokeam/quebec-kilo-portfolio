/**
 * Online Services API
 *
 * For API standards and best practices, see:
 * @see {@link ../../../docs/api-standards.md}
 */

import { axiosInstance } from '@/core/api/client/axios-instance';
import type { OnlineService } from '@/features/dashboard/lib/types/online-services/services';
import type { SelectableItem } from '@/shared/components/ui/ResponsiveCombobox/ResponsiveCombobox';
import { logger } from '@/core/utils/logger/logger';
import {
  DIGITAL_SERVICE_DEFAULTS,
} from '@/core/api/types/api.types';
import { mapToApiServiceType, type ServiceType } from '@/shared/constants/service.constants';


export interface CreateOnlineServiceRequest {
  id?: string;
  name: string;
  parentId: string | null;
  type: 'basic' | 'subscription';
  parentLocationId: string;
  is_active?: boolean;
  metadata: {
    service: OnlineService | null;
    expenseType?: string;
    cost: number;
    billingPeriod?: string;
    nextPaymentDate?: Date;
    paymentMethod?: SelectableItem;
  };
}

export interface OnlineServiceResponse {
  id: string;
  name: string;
  service: OnlineService;
  expenseType?: string;
  cost: number;
  billingPeriod?: string;
  nextPaymentDate?: Date;
  paymentMethod?: SelectableItem;
  createdAt: Date;
  updatedAt: Date;
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
 *
 * @example
 * ```typescript
 * try {
 *   const service = await createOnlineService({
 *     name: 'Netflix',
 *     type: 'digital',
 *     parentId: null,
 *     parentLocationId: 'root',
 *     metadata: {
 *       service: netflixService,
 *       cost: 15.99,
 *       expenseType: '1 month'
 *     }
 *   });
 *   console.log(service);
 * } catch (error) {
 *   console.error('Creation failed:', error);
 * }
 * ```
 */
export const createOnlineService = async (serviceData: CreateOnlineServiceRequest): Promise<OnlineServiceResponse> => {
  logger.debug('📝 Creating new online service', { serviceData });

  // Transform service data to backend format
  const digitalLocation = {
    name: serviceData.name,
    service_type: mapToApiServiceType(serviceData.type as ServiceType),
    is_active: serviceData.is_active !== undefined ? serviceData.is_active : true,
    url: serviceData.metadata.service?.url && serviceData.metadata.service.url !== '#'
        ? serviceData.metadata.service.url
        : DIGITAL_SERVICE_DEFAULTS.URL,
    ...(serviceData.metadata.service?.isSubscriptionService && {
        subscription: {
            billing_cycle: serviceData.metadata.billingPeriod || serviceData.metadata.expenseType || '1 month',
            cost_per_cycle: serviceData.metadata.cost || 0,
            next_payment_date: serviceData.metadata.nextPaymentDate?.toISOString() || new Date().toISOString(),
            payment_method: serviceData.metadata.paymentMethod?.displayName || serviceData.metadata.service?.billing?.paymentMethod || 'None'
        }
    })
} ;

  try {
    const response = await axiosInstance.post<OnlineServiceResponse>(
      '/v1/locations/digital',
      digitalLocation
    );

    // Debug: Log the full response
    logger.debug('🔄 Full backend response:', { response });

    // Ensure the response has the required fields
    if (!response || typeof response !== 'object') {
      throw new Error(`Invalid response: Expected an object, got ${JSON.stringify(response)}`);
    }

    return response;
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
 * @param {CreateOnlineServiceRequest} serviceData - The data for updating the online service
 * @returns {Promise<OnlineServiceResponse>} A promise that resolves to the updated online service
 */
export const updateOnlineService = async (serviceData: CreateOnlineServiceRequest): Promise<OnlineServiceResponse> => {
  logger.debug('📝 Updating online service', { serviceData });

  const digitalLocation = {
    id : serviceData.id,
    name: serviceData.name,
    service_type: serviceData.type,
    is_active: serviceData.is_active !== undefined ? serviceData.is_active : true,
    url: serviceData.metadata.service?.url && serviceData.metadata.service.url !== '#'
      ? serviceData.metadata.service.url
      : DIGITAL_SERVICE_DEFAULTS.URL,
    ...(serviceData.metadata.service?.isSubscriptionService && {
      subscription: {
        billing_cycle: serviceData.metadata.billingPeriod || serviceData.metadata.expenseType || '1  month',
        cost_per_cycle: serviceData.metadata.cost || 0,
        next_payment_date: serviceData.metadata.nextPaymentDate?.toISOString() || new Date().toISOString(),
        payment_method: serviceData.metadata.paymentMethod?.label || 'Unknown'
      }
    })
  };

  try {
    const response = await axiosInstance.put<OnlineServiceResponse>(
      `/v1/locations/digital/${serviceData.id}`,
      digitalLocation
    );

    logger.debug('🔄 Full backend response:', { response });

    if (!response || typeof response !== 'object') {
      throw new Error(`Invalid response: Expected an object, got ${JSON.stringify(response)}`);
    }

    return response;
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
