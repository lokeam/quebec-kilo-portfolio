/**
 * Dashboard Service
 *
 * Provides functions for managing dashboard data through the backend API.
 */

import { axiosInstance } from '@/core/api/client/axios-instance';
import { apiRequest } from '@/core/api/utils/apiRequest';
import { logger } from '@/core/utils/logger/logger';

// Types
export type MediaTypeDomain = "games" | "movies" | "oneTimePurchase" | "hardware" | "dlc" | "inGamePurchase" | "subscription";

export interface GameStats {
  title: string;
  icon: string;
  value: number;
  lastUpdated: number;
}

export interface LocationStats {
  title: string;
  icon: string;
  value: number;
  lastUpdated: number;
}

export interface DigitalLocation {
  logo: string;
  name: string;
  url: string;
  billingCycle: string;
  monthlyFee: number;
  storedItems: number;
}

export interface Sublocation {
  sublocationId: string;
  sublocationName: string;
  sublocationType: string;
  storedItems: number;
  parentLocationId: string;
  parentLocationName: string;
  parentLocationType: string;
  parentLocationBgColor: string;
}

export interface PlatformDistribution {
  platform: string;
  itemCount: number;
}

export interface MonthlyExpenditure {
  date: string;
  oneTimePurchase: number;
  hardware: number;
  dlc: number;
  inGamePurchase: number;
  subscription: number;
}

export interface DashboardResponse {
  // Basic Statistics
  gameStats: GameStats;
  subscriptionStats: GameStats;
  digitalLocationStats: LocationStats;
  physicalLocationStats: LocationStats;

  // Digital Locations
  subscriptionTotal: number;
  subscriptionRecurringNextMonth: number;
  digitalLocations: DigitalLocation[];

  // Storage Locations
  sublocations: Sublocation[];

  // Platform Distribution
  newItemsThisMonth: number;
  platformList: PlatformDistribution[];

  // Monthly Spending
  mediaTypeDomains: MediaTypeDomain[];
  monthlyExpenditures: MonthlyExpenditure[];
}

interface DashboardResponseWrapper {
  success: boolean;
  dashboard: DashboardResponse;
  metadata: {
    timestamp: string;
    request_id: string;
  };
}

// Endpoints
const DASHBOARD_BFF_ENDPOINT = '/v1/dashboard/bff';

/**
 * Fetches all dashboard data for the current user.
 *
 * This function retrieves all dashboard data in a single request, including:
 * - Basic statistics (games, subscriptions, locations)
 * - Digital and physical storage locations
 * - Platform distribution
 * - Monthly spending data
 *
 * Uses apiRequest helper to wrap the axios call with:
 *  - async/await syntax
 *  - pre‑call debug log
 *  - post‑call success log
 *  - catch block with error log + optional Sentry/metrics
 *  - retry logic (if configured)
 *
 * @returns Promise<DashboardResponse> The dashboard data
 * @throws Error if no dashboard data is found in the response
 */
export const getDashboardBFFResponse = (): Promise<DashboardResponse> =>
  apiRequest('getDashboardBFFResponse', async () => {
    logger.debug('getDashboardBFFResponse: Making API request');
    const response = await axiosInstance.get<DashboardResponseWrapper>(DASHBOARD_BFF_ENDPOINT);
    logger.debug('getDashboardBFFResponse: Raw API response:', response.data);

    if (!response.data.dashboard) {
      logger.error('getDashboardBFFResponse: No dashboard data in response:', response.data);
      throw new Error('No dashboard data in response');
    }

    logger.debug('getDashboardBFFResponse: Successfully extracted dashboard data:', response.data.dashboard);
    return response.data.dashboard;
  });