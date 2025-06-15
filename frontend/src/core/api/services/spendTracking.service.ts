/**
 * Spend Tracking Service
 *
 * Provides functions for managing spend tracking through the backend API.
 */

import { axiosInstance } from '@/core/api/client/axios-instance';
import { apiRequest } from '@/core/api/utils/apiRequest';

// Types
import type { SpendItem, YearlySpending } from '@/types/domain/spend-tracking';

const SPEND_TRACKING_ENDPOINT = '/v1/spend-tracking';
const SPEND_TRACKING_BFF_ENDPOINT = '/v1/spend-tracking/bff';

interface SpendTrackingOperationResponseWrapper {
  success: boolean;
  spendTracking: {
    id: number;
    message: string;
  };
  metadata: {
    timestamp: string;
    request_id: string;
  };
}

interface SpendTrackingItemResponseWrapper {
  success: boolean;
  spendTracking: {
    item: SpendItem;
  };
  metadata: {
    timestamp: string;
    request_id: string;
  };
}

interface SpendTrackingBFFResponseWrapper {
  success: boolean;
  spendTracking: {
    currentMonthItems: SpendItem[];
    nextMonthItems: SpendItem[];
    yearlyTotals: YearlySpending[];
  };
  metadata: {
    timestamp: string;
    request_id: string;
  };
}

type SpendTrackingOperationResponse = {
  id: number;
  message: string;
};

type SpendTrackingBFFResponse = {
  currentMonthItems: SpendItem[];
  nextMonthItems: SpendItem[];
  yearlyTotals: YearlySpending[];
};

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
    console.log('[DEBUG] getSpendTrackingPageBFFResponse: Making API request');
    const response = await axiosInstance.get<SpendTrackingBFFResponseWrapper>(SPEND_TRACKING_BFF_ENDPOINT);
    console.log('[DEBUG] getSpendTrackingPageBFFResponse: Raw API response:', response.data);

    if (!response.data.spendTracking) {
      console.error('[DEBUG] getSpendTrackingPageBFFResponse: No spend tracking data in response:', response.data);
    }

    console.log('[DEBUG] getSpendTrackingPageBFFResponse: Successfully extracted spend tracking data:', response.data.spendTracking);
    return response.data.spendTracking;
  });

/**
 * Fetches a specific spend item by ID
 */
export const getSpendTrackingItemById = (id: string): Promise<SpendItem> =>
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
export const createSpendTrackingItem = (data: Omit<SpendItem, 'id'>): Promise<SpendTrackingOperationResponse> =>
  apiRequest('createSpendItem', () =>
    axiosInstance
      .post<SpendTrackingOperationResponseWrapper>(SPEND_TRACKING_ENDPOINT, data)
      .then(response => {
        const item = response.data.spendTracking;
        if (!item) {
          throw new Error('Failed to create spend item');
        }
        return item;
      })
  );

/**
 * Updates an existing spend item
 */
export const updateSpendTrackingItem = (id: string, data: Partial<SpendItem>): Promise<SpendTrackingOperationResponse> =>
  apiRequest(`updateSpendItem(${id})`, () =>
    axiosInstance
      .put<SpendTrackingOperationResponseWrapper>(`${SPEND_TRACKING_ENDPOINT}/${id}`, data)
      .then(response => {
        const item = response.data.spendTracking;
        if (!item) {
          throw new Error(`Failed to update spend item with id ${id}`);
        }
        return item;
      })
  );

/**
 * Deletes an existing spend item
 */
export const deleteSpendTrackingItem = (id: string): Promise<SpendTrackingOperationResponse> =>
  apiRequest(`deleteSpendItem(${id})`, () =>
    axiosInstance
      .delete<SpendTrackingOperationResponseWrapper>(`${SPEND_TRACKING_ENDPOINT}/${id}`)
      .then(response => response.data.spendTracking)
  );