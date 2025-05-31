import type { ServiceProviderRecord } from "../types/services";

/**
 * Physical item conditions for inventory tracking.
 * Applies to hardware and physical media.
 * Used in MediaStorage and SpendTracking services.
 */
export type ItemCondition =
  | 'new'
  | 'used'
  | 'refurbished';

/**
 * Maps service providers to their display names and internal identifiers
 * Used for consistent provider identification across the application
 */
export const ONLINE_SERVICE_PROVIDERS: ServiceProviderRecord = {
  APPLE: {
    displayName: 'Apple Arcade',
    id: 'apple',
  },
  EA: {
    displayName: 'EA Play',
    id: 'ea',
  },
  EPIC: {
    displayName: 'Epic Games Store',
    id: 'epic',
  },
  FANATICAL: {
    displayName: 'Fanatical',
    id: 'fanatical',
  },
  GOG: {
    displayName: 'GOG.com',
    id: 'gog',
  },
  GOOGLE: {
    displayName: 'Google Play Pass',
    id: 'playpass',
  },
  GREENMAN: {
    displayName: 'Green Man Gaming',
    id: 'greenman',
  },
  HUMBLE: {
    displayName: 'Humble Bundle',
    id: 'humble',
  },
  META: {
    displayName: 'Meta Quest+',
    id: 'meta',
  },
  MICROSOFT: {
    displayName: 'Xbox Network',
    id: 'xbox',
  },
  NETFLIX: {
    displayName: 'Netflix Games',
    id: 'netflix',
  },
  NINTENDO: {
    displayName: 'Nintendo',
    id: 'nintendo',
  },
  NVIDIA: {
    displayName: 'GeForce Now',
    id: 'nvidia',
  },
  PRIME: {
    displayName: 'Prime Gaming',
    id: 'prime',
  },
  SONY: {
    displayName: 'Playstation Plus',
    id: 'playstation',
  },
  STEAM: {
    displayName: 'Steam',
    id: 'steam',
  },
} as const;

/**
 * Display names for online service providers as they appear in the UI.
 * Extracted from ONLINE_SERVICE_PROVIDERS.displayName.
 */
export type OnlineServiceProviderDisplay =
  typeof ONLINE_SERVICE_PROVIDERS[keyof typeof ONLINE_SERVICE_PROVIDERS]['displayName'];

/**
 * Unique identifiers for online service providers used in internal operations.
 * Extracted from ONLINE_SERVICE_PROVIDERS.id.
 */
export type OnlineServiceProviderId =
  typeof ONLINE_SERVICE_PROVIDERS[keyof typeof ONLINE_SERVICE_PROVIDERS]['id'];

/**
 * Status codes for service operational states.
 * Used to track service availability and health.
 */
export const SERVICE_STATUS_CODES = {
  ACTIVE: 'active',
  INACTIVE: 'inactive',
  ERROR: 'error',
} as const;

export type ServiceStatusCode = typeof SERVICE_STATUS_CODES[keyof typeof SERVICE_STATUS_CODES];


/**
 * Categories of services supported by the application.
 * Used for service classification and filtering.
 */
export const SERVICE_TYPES = {
  ONLINE: 'online', // Digital services accessed via internet
  PHYSICAL: 'physical', // Physical items, such as hardware or physical media
  SUBSCRIPTION: 'subscription', // Recurring payment services
} as const;

/**
 * Unique identifiers for service types used in internal operations.
 * Extracted from SERVICE_TYPES.id.
 */
export type ServiceType = typeof SERVICE_TYPES[keyof typeof SERVICE_TYPES];


/**
 * API service types used specifically for backend communication.
 * These match the exact string values expected by the backend API.
 */
export const API_SERVICE_TYPES = {
  BASIC: 'basic',
  SUBSCRIPTION: 'subscription',
} as const;

export type ApiServiceType = typeof API_SERVICE_TYPES[keyof typeof API_SERVICE_TYPES];

/**
 * UPDATE: 4/22 - Maps frontend service types to their corresponding backend API service types.
 * Use this when sending data to the backend to ensure correct type values.
 */
export function mapToApiServiceType(frontendType: ServiceType): ApiServiceType {
  switch (frontendType) {
    case SERVICE_TYPES.ONLINE:
      return API_SERVICE_TYPES.BASIC;
    case SERVICE_TYPES.SUBSCRIPTION:
      return API_SERVICE_TYPES.SUBSCRIPTION;
    case SERVICE_TYPES.PHYSICAL:
      // Assuming physical services map to basic in your API
      return API_SERVICE_TYPES.BASIC;
    default:
      // Type safety - should never happen with proper TypeScript usage
      return API_SERVICE_TYPES.BASIC;
  }
}

/**
 * Maps API service types back to frontend service types.
 * Use this when receiving data from the backend.
 */
export function mapFromApiServiceType(apiType: ApiServiceType): ServiceType {
  switch (apiType) {
    case API_SERVICE_TYPES.SUBSCRIPTION:
      return SERVICE_TYPES.SUBSCRIPTION;
    case API_SERVICE_TYPES.BASIC:
    default:
      return SERVICE_TYPES.ONLINE; // Assuming basic maps to online in your frontend
  }
}
