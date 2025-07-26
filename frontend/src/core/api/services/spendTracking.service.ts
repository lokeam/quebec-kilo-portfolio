/**
 * Spend Tracking Service
 *
 * Provides functions for managing spend tracking through the backend API.
 */

import { axiosInstance } from '@/core/api/client/axios-instance';
import { apiRequest } from '@/core/api/utils/apiRequest';

// Constants
import { API_BASE_PATH } from '@/core/api/config';

// Types
import type {
  CreateOneTimePurchaseRequest,
  SpendingItemBFFResponse,
  SpendTrackingBFFResponse,
  SpendTrackingDeleteResponse,
} from '@/types/domain/spend-tracking';

const SPEND_TRACKING_ENDPOINT = `${API_BASE_PATH}/spend-tracking`;
const SPEND_TRACKING_BFF_ENDPOINT = `${API_BASE_PATH}/spend-tracking/bff`;

interface SpendTrackingDeleteResponseWrapper {
  success: boolean;
  spendTracking: SpendTrackingDeleteResponse;
  metadata: {
    timestamp: string;
    request_id: string;
  };
}

interface SpendTrackingItemResponseWrapper {
  success: boolean;
  spendTracking: {
    item: SpendingItemBFFResponse;
  };
  metadata: {
    timestamp: string;
    request_id: string;
  };
}

interface SpendTrackingBFFResponseWrapper {
  success: boolean;
  spendTracking: SpendTrackingBFFResponse;
  metadata: {
    timestamp: string;
    request_id: string;
  };
}

/**
 * Fetches all spend tracking data for the BFF page.
 *
 * Uses apiRequest helper to wrap the axios call with:
 *  - async/await syntax
 *  - pre‑call debug log
 *  - post‑call success log
 *  - catch block with error log + optional Sentry/metrics
 *  - retry logic (if configured)
 */
export const getSpendTrackingPageBFFResponse = (): Promise<SpendTrackingBFFResponse> =>
  apiRequest('getSpendTrackingPageBFFResponse', async () => {
    // console.log('[DEBUG] getSpendTrackingPageBFFResponse: Making API request');
    const response = await axiosInstance.get<SpendTrackingBFFResponseWrapper>(SPEND_TRACKING_BFF_ENDPOINT);
    // console.log('[DEBUG] getSpendTrackingPageBFFResponse: Raw API response:', response.data);

    if (!response.data.spendTracking) {
      console.error('[DEBUG] getSpendTrackingPageBFFResponse: No spend tracking data in response:', response.data);
    }

    // console.log('[DEBUG] getSpendTrackingPageBFFResponse: Successfully extracted spend tracking data:', response.data.spendTracking);
    return response.data.spendTracking;
  });

/**
 * Fetches a specific spend item by ID
 */
export const getSpendTrackingItemById = (id: string): Promise<SpendingItemBFFResponse> =>
  apiRequest(`getSpendItemById(${id})`, () =>
    axiosInstance
      .get<SpendTrackingItemResponseWrapper>(`${SPEND_TRACKING_ENDPOINT}/${id}`)
      .then(response => {
        const item = response.data.spendTracking.item;
        if (!item) {
          throw new Error(`Spend item with id ${id} not found`);
        }
        return item;
      })
  );

/**
 * Creates a new spend item
 */
export const createOneTimePurchase = (input: CreateOneTimePurchaseRequest): Promise<SpendingItemBFFResponse> =>
  apiRequest('createSpendTrackingItem', async () => {
    // console.log('[DEBUG] createSpendTrackingItem: Making API request');
    // console.log('[DEBUG] createSpendTrackingItem: Request payload:', input);

    const response = await axiosInstance.post<SpendTrackingItemResponseWrapper>(SPEND_TRACKING_ENDPOINT, input);
    // console.log('[DEBUG] createSpendTrackingItem: Raw API response:', response.data);

    if (!response.data.spendTracking) {
      console.error('[DEBUG] createSpendTrackingItem: No spend tracking data in response:', response.data);
    }

    // console.log('[DEBUG] createSpendTrackingItem: Successfully extracted spend tracking data:', response.data.spendTracking.item);
    return response.data.spendTracking.item;
  });

/**
 * Updates an existing spend item
 */
export const updateSpendTrackingItem = (id: string, data: Partial<CreateOneTimePurchaseRequest>): Promise<SpendingItemBFFResponse> =>
  apiRequest(`updateSpendItem(${id})`, () =>
    axiosInstance
.put<SpendTrackingItemResponseWrapper>(`${SPEND_TRACKING_ENDPOINT}/${id}`, data)
  .then(response => {
    const item = response.data.spendTracking.item;
    if (!item) {
      throw new Error(`Failed to update spend item with id ${id}`);
    }
    return item;
  })
  );

/**
 * Deletes an array existing spend items (frontend UI currently only allows for single item deletion. Expand upon this feature in future)
 */
export const deleteSpendTrackingItems = (ids: string[]): Promise<SpendTrackingDeleteResponse> =>
  apiRequest(`deleteSpendItems(${ids.join(',')})`, () =>
    axiosInstance
      .delete<SpendTrackingDeleteResponseWrapper>(`${SPEND_TRACKING_ENDPOINT}?ids=${ids.join(',')}`)
      .then(response => response.data.spendTracking)
  );