/**
 * Analytics API Service
 *
 * This service handles API requests related to analytics data
 * It interfaces with the backend GET /api/v1/analytics endpoint
 */

import { axiosInstance } from '@/core/api/client/axios-instance';
import { apiRequest } from '@/core/api/utils/apiRequest';

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

export interface LocationSummary {
  id: string;
  name: string;
  itemCount: number;
  locationType: string;
  mapCoordinates?: string;
  createdAt: string;
  updatedAt: string;
  isSubscription?: boolean;
  monthlyCost?: number;
  sublocations?: {
    id: string;
    name: string;
    locationType: string;
    bgColor?: string;
    storedItems: number;
    createdAt: string;
    updatedAt: string;
    items?: Array<{
      id: number;
      name: string;
      platform: string;
      platformVersion: string;
      acquiredDate: string;
    }>;
  }[];
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

// to delete
export interface AnalyticsResponse {
  general?: GeneralStats;
  financial?: FinancialStats;
  storage?: StorageStats;
  inventory?: InventoryStats;
  wishlist?: WishlistStats;
}

export interface AnalyticsResponseWrapper {
  storage?: {
    digitalLocations: LocationSummary[];
    physicalLocations: LocationSummary[];
    totalDigitalLocations: number;
    totalPhysicalLocations: number;
  };
  general?: GeneralStats;
  financial?: FinancialStats;
  inventory?: InventoryStats;
  wishlist?: WishlistStats;
}

/**
 * Fetches analytics data for specified domains
 *
 */
export const getAnalyticsData = (domains?: string[]): Promise<AnalyticsResponseWrapper> =>
  apiRequest('getAnalytics', () =>
    axiosInstance
      .get<{ analytics: AnalyticsResponseWrapper }>('/v1/analytics', {
        params: domains ? { domains: domains.join(',') } : undefined
      })
      .then(response => response.data.analytics)
  );
