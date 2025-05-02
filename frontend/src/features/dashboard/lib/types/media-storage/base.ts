import type { GameItem } from '@/features/dashboard/lib/types/media-storage/items';
import type { LocationType } from '@/features/dashboard/lib/types/media-storage/constants';

/**
 * Base interface for all storage locations in the media storage system.
 * This serves as the foundation for both physical and digital locations.
 *
 * @interface BaseLocation
 */
export interface BaseLocation {
  /** Unique identifier for the location */
  id: string;

  /** Human-readable name of the location */
  name: string;

  /** URL-friendly identifier for the location */
  label: string;

  /** Discriminator to distinguish between physical and digital locations */
  type: LocationType;

  /** Optional detailed description of the location */
  description?: string;

  /** List of media items stored directly in this location */
  items?: GameItem[];

  /** Timestamp of when the location was created */
  createdAt: Date;

  /** Timestamp of when the location was last modified */
  updatedAt: Date;

  /** Additional flexible metadata for the location */
  metadata?: Record<string, unknown>;

  /** Number of items stored in sub-locations */
  itemsInSublocations: number;
}

/**
 * Metrics interface for tracking storage location statistics.
 * Used for analytics and reporting purposes.
 *
 * @interface LocationMetrics
 */
export interface LocationMetrics {
  /** Total number of items in this location */
  totalItems: number;

  /** Number of items stored in sub-locations */
  itemsInSublocations: number;

  /** Timestamp of the last metrics update */
  lastUpdated: Date;
}

/**
 * Simplified location data for list views and summaries.
 * Contains essential information for location overview displays.
 *
 * @interface LocationSummary
 */
export interface LocationSummary {
  /** Unique identifier for the location */
  id: string;

  /** Display name of the location */
  name: string;

  /** Type of storage location */
  type: LocationType;

  /** Total number of items in this location */
  itemCount: number;

  /** Timestamp of the last modification */
  lastUpdated: Date;
}

/**
 * Parameters for querying and filtering storage locations.
 * Supports pagination, sorting, and flexible filtering options.
 *
 * @interface StorageQueryParams
 */
export interface StorageQueryParams {
  /** Page number for pagination */
  page?: number;

  /** Number of items per page */
  limit?: number;

  /** Field to sort results by */
  sortBy?: string;

  /** Sort direction */
  sortOrder?: 'asc' | 'desc';

  /** Key-value pairs for filtering results */
  filterBy?: Record<string, string>;
}
