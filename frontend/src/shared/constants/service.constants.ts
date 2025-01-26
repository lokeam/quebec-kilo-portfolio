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
    displayName: 'Playstation Network',
    id: 'playstation',
  },
  STEAM: {
    displayName: 'Steam',
    id: 'steam',
  },
} as const;

export type OnlineServiceProviderDisplay =
  typeof ONLINE_SERVICE_PROVIDERS[keyof typeof ONLINE_SERVICE_PROVIDERS]['displayName'];

export type OnlineServiceProviderId =
  typeof ONLINE_SERVICE_PROVIDERS[keyof typeof ONLINE_SERVICE_PROVIDERS]['id'];
