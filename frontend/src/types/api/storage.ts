import type { LocationIconBgColor } from '@/types/domain/location-types';

/**
 * Storage Analytics Types
 * These types exactly match the shape of data returned by the /v1/analytics?domains=storage endpoint
 */

export interface StorageAnalytics {
  storage: {
    digital_locations: DigitalLocation[];
    physical_locations: PhysicalLocation[];
    total_digital_locations: number;
    total_physical_locations: number;
  };
}

export interface DigitalLocation {
  id: string;
  created_at: string;
  is_active: boolean;
  is_subscription: boolean;
  item_count: number;
  location_type: 'subscription' | 'basic';
  monthly_cost: number;
  name: string;
  updated_at: string;
  url: string;
}

export interface PhysicalLocation {
  id: string;
  created_at: string;
  item_count: number;
  location_type: 'house' | 'apartment' | 'office' | 'warehouse' | 'vehicle';
  bgColor?: LocationIconBgColor;
  name: string;
  sublocations: Sublocation[];
  updated_at: string;
}

export interface Sublocation {
  id: string;
  created_at: string;
  location_type: 'shelf' | 'console' | 'cabinet' | 'closet' | 'drawer' | 'box' | 'device';
  name: string;
  stored_items: number;
  updated_at: string;
}

/**
 * Location counts for media storage metadata
 */
export interface LocationCounts {
  /** Total count of storage locations */
  total: number;
  /** Number of physical storage locations */
  physical: number;
  /** Number of digital storage locations */
  digital: number;
}

/**
 * Item counts for individual storage locations
 */
export interface LocationItemCounts {
  /** Total number of items in this location */
  total: number;
  /** Number of items stored in sub-locations */
  inSublocations: number;
}

/**
 * Comprehensive item counting interface
 */
export interface ItemCounts {
  /** Total number of items across all locations */
  total: number;
  /** Number of physical items */
  physical: number;
  /** Number of digital items */
  digital: number;
  /** Map of location IDs to item counts */
  byLocation: Record<string, LocationItemCounts>;
}

/**
 * Top-level metadata interface for media storage system
 */
export interface MediaStorageMetadata {
  /** Aggregated count info */
  counts: {
    /** Location-related counts */
    locations: LocationCounts;
    /** Item-related counts */
    items: ItemCounts;
  };
  /** Timestamp of most recent metadata update */
  lastUpdated: Date;
  /** Schema version for backwards compatibility */
  version: string;
}