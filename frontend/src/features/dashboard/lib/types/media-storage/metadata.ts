/**
 * Aggregate counts for different types of storage locations.
 * Used for system-wide statistics and reporting.
 *
 * @interface LocationCounts
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
 * Item counts for individual storage locations.
 * Tracks both direct items + items in sub-locations.
 *
 * @interface LocationItemCounts
 */
export interface LocationItemCounts {
  /** Total number of items in this location */
  total: number;

  /** Number of items stored in sub-locations */
  inSubLocations: number;
}

/**
 * Comprehensive item counting interface.
 * Provides both system-wide totals + per-location breakdowns.
 *
 * @interface ItemCounts
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
 * Top-level metadata interface for media storage system.
 * Contains aggregated statistics + system information.
 *
 * @interface MediaStorageMetadata
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
  lastUpdated?: Date;

  /** Schema version for backwards compatibility */
  version?: string;
}