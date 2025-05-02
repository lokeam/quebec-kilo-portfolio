import type { BaseLocation } from '@/features/dashboard/lib/types/media-storage/base';
import type { PhysicalLocationType, SublocationType } from '@/features/dashboard/lib/types/media-storage/constants';
import type { GameItem } from '@/features/dashboard/lib/types/media-storage/items';

/**
 * Configuration for subdivisions within physical storage locations.
 * Represents organizational units such as shelves, cabinets + boxes.
 *
 * @interface Sublocation
 */
export interface Sublocation {
  /** Unique identifier for the sub-location */
  id: string;

  /** Display name of the sub-location */
  name: string;

  /** Detailed description of the sublocation */
  description?: string;

  /** Type of storage unit */
  locationType: SublocationType;

  /** Items stored in this sub-location */
  items?: GameItem[];

  /** Maximum number of items this sub-location can hold */
  capacity?: number;

  /** Whether the sub-location is currently accessible */
  isAccessible?: boolean;

  /** Background color for the sub-location */
  bgColor?: string;

  /** Number of stored items */
  storedItems?: number;

  /** Timestamp of when the sub-location was created */
  createdAt?: Date;

  /** Timestamp of when the sub-location was last modified */
  updatedAt?: Date;
}

/**
 * Configuration for physical storage locations.
 * Extends BaseLocation with properties specific to real-world locations.
 *
 * @interface PhysicalLocation
 * @extends {BaseLocation}
 */
export interface PhysicalLocation extends BaseLocation {
  /** Type of physical location (ex: house, apartment) */
  locationType: PhysicalLocationType;

  /** Geographic coordinates of the location */
  mapCoordinates?: string;

  /** List of storage subdivisions within this location */
  sublocations?: Sublocation[];

  /** Parent location ID */
  parentLocationId?: string;
}

/**
 * Input type for creating new physical locations.
 * Omits server-generated fields from PhysicalLocation.
 *
 * @type CreatePhysicalLocationInput
 */
export type CreatePhysicalLocationInput = Omit<
  PhysicalLocation,
  'id' | 'createdAt' | 'updatedAt'
>;

/**
 * Input type for updating existing physical locations.
 * Makes all fields optional except id.
 *
 * @type UpdatePhysicalLocationInput
 */
export type UpdatePhysicalLocationInput = Partial<
  Omit<PhysicalLocation, 'id' | 'createdAt' | 'updatedAt'>
>;
