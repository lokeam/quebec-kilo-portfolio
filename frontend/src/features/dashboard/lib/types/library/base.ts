import type { ISO8601Date } from '@/shared/types/types';
import type { BaseLocation, StorageUnit } from '@/shared/types/location';

/**
 * Base interface for media items in the library
 * Common properties shared across all library items
 */
export interface BaseMediaItem {
  id: string;
  title: string;
  imageUrl: string;
  favorite: boolean;
  dateAdded: ISO8601Date;
}

/**
 * Generic type for items with optional metadata
 * T represents additional properties specific to the item type
 */
export type WithMetadata<T> = {
  metadata?: T;
};

/**
 * Represents a physical location where library items are stored
 * Extends BaseLocation with additional properties specific to the library context
 */
export interface LibraryLocation extends BaseLocation {
  gameCount?: number;
  sublocation?: {
    type: StorageUnit;
    name: string;
  };
}
