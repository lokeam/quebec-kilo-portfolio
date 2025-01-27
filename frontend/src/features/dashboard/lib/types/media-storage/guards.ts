import type { PhysicalLocation } from '@/features/dashboard/lib/types/media-storage/physical';
import type { DigitalLocation } from '@/features/dashboard/lib/types/media-storage/digital';
import type { StorageLocation } from '@/features/dashboard/lib/types/media-storage/aggregates';
import { LocationType } from '@/features/dashboard/lib/types/media-storage/constants';

/**
 * Type guard functions for media storage locations.
 * These functions provide runtime type checking + TypeScript type narrowing.
 */

/**
 * Determines if storage location is a brick and mortar.
 *
 * @param {StorageLocation} location - The location to check
 * @returns {boolean} True if the location is a physical location
 *
 * @example
 * ```typescript
 * if (isPhysicalLocation(location)) {
 *   // location is typed as PhysicalLocation
 *   console.log(location.locationType);
 * }
 * ```
 */
export const isPhysicalLocation = (
  location: StorageLocation
): location is PhysicalLocation => {
  return location.type === LocationType.PHYSICAL;
};

/**
 * Determines if a storage location is a digital location.
 *
 * @param {StorageLocation} location - The location to check
 * @returns {boolean} True if the location is a digital location
 *
 * @example
 * ```typescript
 * if (isDigitalLocation(location)) {
 *   // location is typed as DigitalLocation
 *   console.log(location.platform);
 * }
 * ```
 */
export const isDigitalLocation = (
  location: StorageLocation
): location is DigitalLocation => {
  return location.type === LocationType.DIGITAL;
};
