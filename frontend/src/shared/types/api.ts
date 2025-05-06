/**
 * API Response Types
 * These types exactly match the shape of data returned by our backend API
 */

/**
 * Digital Location as returned by the API
 */
export interface ApiDigitalLocation {
  id: string;
  name: string;
  locationType: string;
  itemCount: number;
  isSubscription: boolean;
  monthlyCost: number;
}

/**
 * Physical Location as returned by the API
 */
export interface ApiPhysicalLocation {
  id: string;
  name: string;
  locationType: string;
  itemCount: number;
}

/**
 * Storage Analytics Response
 */
export interface ApiStorageAnalytics {
  storage: {
    totalDigitalLocations: number;
    totalPhysicalLocations: number;
    digitalLocations: ApiDigitalLocation[];
    physicalLocations: ApiPhysicalLocation[];
  };
}

/**
 * Analytics Response
 * This is the top-level response type for the /v1/analytics endpoint
 */
export interface ApiAnalyticsResponse {
  storage?: ApiStorageAnalytics['storage'];
  // Add other analytics domains as needed
  // general?: ApiGeneralAnalytics;
  // financial?: ApiFinancialAnalytics;
  // inventory?: ApiInventoryAnalytics;
  // wishlist?: ApiWishlistAnalytics;
}