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
  totalGames: number;
  monthlySubscriptionCost: number;
  totalDigitalLocations: number;
  totalPhysicalLocations: number;
}

export interface FinancialStats {
  annualSubscriptionCost: number;
  totalServices: number;
  renewalsThisMonth: number;
  services: ServiceDetails[];
}

export interface ServiceDetails {
  name: string;
  monthlyFee: number;
  billingCycle: string;
  nextPayment: string;
}

export interface StorageStats {
  totalPhysicalLocations: number;
  totalDigitalLocations: number;
  digitalLocations: LocationSummary[];
  physicalLocations: LocationSummary[];
}

export interface LocationSummary {
  id: string;
  name: string;
  itemCount: number;
  locationType: string;
  isSubscription?: boolean;
  monthlyCost?: number;
}

export interface InventoryStats {
  totalItemCount: number;
  newItemCount: number;
  platformCounts: PlatformItemCount[];
}

export interface PlatformItemCount {
  platform: string;
  itemCount: number;
}

export interface WishlistStats {
  totalWishlistItems: number;
  itemsOnSale: number;
  starredItem?: string;
  starredItemPrice?: number;
  cheapestSaleDiscount?: number;
}

export interface AnalyticsResponse {
  success: boolean;
  userId: string;
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