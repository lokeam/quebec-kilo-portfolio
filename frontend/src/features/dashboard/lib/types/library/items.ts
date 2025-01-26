import type { BaseMediaItem } from '@/features/dashboard/lib/types/library/base';
import type { ItemCondition } from '@/shared/constants/service.constants';
import type { PhysicalLocation, DigitalLocation } from '@/shared/types/location';
import type { Platform } from '@/shared/types/platform';
import type { StorageSize } from '@/shared/types/storage';
import { MEDIA_ITEM_TYPES } from './constants';

/**
 * Represents a physical item in the library (e.g., game cartridge, disc)
 * Extends BaseMediaItem with physical-specific properties
 */
export interface PhysicalLibraryItem extends BaseMediaItem {
  type: typeof MEDIA_ITEM_TYPES.PHYSICAL;
  location: PhysicalLocation;
  platform: Platform;
  condition?: ItemCondition;
}

/**
 * Represents a digital item in the library (e.g., downloaded game)
 * Extends BaseMediaItem with digital-specific properties
 */
export interface DigitalLibraryItem extends BaseMediaItem {
  type: typeof MEDIA_ITEM_TYPES.DIGITAL;
  location: DigitalLocation & {
    diskSize: StorageSize;  // Required for digital items
  };
  platform: Platform;
}

export type LibraryItem = PhysicalLibraryItem | DigitalLibraryItem;
