import type { PhysicalLocation } from '@/types/domain/physical-location';
import type { DigitalLocation } from '@/types/domain/digital-location';

/**
 * Type guard functions for media storage locations.
 * These functions provide runtime type checking + TypeScript type narrowing.
 */

/**
 * Determines if storage location is a brick and mortar.
 *
 * @param {PhysicalLocation | DigitalLocation} location - The location to check
 * @returns {boolean} True if the location is a physical location
 *
 * @example
 * ```typescript
 * if (isPhysicalLocation(location)) {
 *   // location is typed as PhysicalLocation
 *   console.log(location.type);
 * }
 * ```
 */
export const isPhysicalLocation = (
  location: PhysicalLocation | DigitalLocation
): location is PhysicalLocation => {
  return 'sublocations' in location;
};

/**
 * Determines if a storage location is a digital location.
 *
 * @param {PhysicalLocation | DigitalLocation} location - The location to check
 * @returns {boolean} True if the location is a digital location
 *
 * @example
 * ```typescript
 * if (isDigitalLocation(location)) {
 *   // location is typed as DigitalLocation
 *   console.log(location.type);
 * }
 * ```
 */
export const isDigitalLocation = (
  location: PhysicalLocation | DigitalLocation
): location is DigitalLocation => {
  return 'items' in location;
};
