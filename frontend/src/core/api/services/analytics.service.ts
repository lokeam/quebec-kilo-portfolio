/**
 * Analytics API Service
 *
 * This service handles API requests related to analytics data
 * It interfaces with the backend GET /api/v1/analytics endpoint
 */

import { axiosInstance } from '@/core/api/client/axios-instance';
import { logger } from '@/core/utils/logger/logger';

// Analytics domain constants
export const ANALYTICS_DOMAINS = {
  GENERAL: 'general',
  FINANCIAL: 'financial',
  STORAGE: 'storage',
  INVENTORY: 'inventory',
  WISHLIST: 'wishlist'
} as const;

export type AnalyticsDomain = typeof ANALYTICS_DOMAINS[keyof typeof ANALYTICS_DOMAINS];

// Types from backend analytics_models.go
export interface GeneralStats {
  total_games: number;
  monthly_subscription_cost: number;
  total_digital_locations: number;
  total_physical_locations: number;
}

export interface FinancialStats {
  annual_subscription_cost: number;
  total_services: number;
  renewals_this_month: number;
  services: ServiceDetails[];
}

export interface ServiceDetails {
  name: string;
  monthly_fee: number;
  billing_cycle: string;
  next_payment: string;
}

export interface StorageStats {
  total_physical_locations: number;
  total_digital_locations: number;
  digital_locations: LocationSummary[];
  physical_locations: LocationSummary[];
}

export interface LocationSummary {
  id: string;
  name: string;
  item_count: number;
  location_type: string;
  is_subscription?: boolean;
  monthly_cost?: number;
}

export interface InventoryStats {
  total_item_count: number;
  new_item_count: number;
  platform_counts: PlatformItemCount[];
}

export interface PlatformItemCount {
  platform: string;
  item_count: number;
}

export interface WishlistStats {
  total_wishlist_items: number;
  items_on_sale: number;
  starred_item?: string;
  starred_item_price?: number;
  cheapest_sale_discount?: number;
}

export interface AnalyticsResponse {
  success: boolean;
  user_id: string;
  data: {
    general?: GeneralStats;
    financial?: FinancialStats;
    storage?: StorageStats;
    inventory?: InventoryStats;
    wishlist?: WishlistStats;
  };
}

/**
 * Fetches analytics data for specified domains
 *
 * @async
 * @function getAnalyticsData
 * @param {AnalyticsDomain[]} domains - Array of analytics domains to fetch
 * @param {string} [token] - Optional auth token
 * @returns {Promise<AnalyticsResponse>} A promise that resolves to analytics data
 *
 * @throws {Error} If the API request fails
 */
export const getAnalyticsData = async (
  domains: AnalyticsDomain[],
  token?: string
): Promise<AnalyticsResponse> => {
  logger.debug('Fetching analytics data', { domains });

  try {
    // Build query params for domains
    const queryParams = domains.map(domain => `domains=${domain}`).join('&');
    const url = `/v1/analytics${queryParams ? `?${queryParams}` : ''}`;

    const config = {
      headers: token
        ? {
            Authorization: `Bearer ${token}`,
            'Content-Type': 'application/json',
            Accept: 'application/json'
          }
        : undefined
    };

    const response = await axiosInstance.get<AnalyticsResponse>(url, config);

    logger.debug('Analytics data fetched successfully', {
      domains,
      dataKeys: Object.keys(response.data || {})
    });

    if (!response || !response.success) {
      throw new Error('Invalid response from analytics API');
    }

    return response;
  } catch (error) {
    logger.error('Failed to fetch analytics data', { domains, error });
    throw error;
  }
};