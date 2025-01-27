import type { BaseLocation } from '@/features/dashboard/lib/types/media-storage/base';
import type { PhysicalLocationType, SubLocationType } from '@/features/dashboard/lib/types/media-storage/constants';
import type { GameItem } from '@/features/dashboard/lib/types/media-storage/items';

/**
 * Configuration for subdivisions within physical storage locations.
 * Represents organizational units such as shelves, cabinets + boxes.
 *
 * @interface SubLocation
 */
export interface SubLocation {
  /** Unique identifier for the sub-location */
  id: string;

  /** Display name of the sub-location */
  name: string;

  /** Detailed description of the sublocation -- NOTE: determine if this is worthwhile in UAT*/
  description: string;

  /** Type of storage unit */
  locationType: SubLocationType;

  /** Items stored in this sub-location */
  items?: GameItem[];

  /** Maximum number of items this sub-location can hold */
  capacity?: number;

  /** Whether the sub-location is currently accessible */
  isAccessible: boolean;
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

  /** Street address of the location */
  address?: string;

  /** List of storage subdivisions within this location */
  subLocations?: SubLocation[];

  /** Hours during which the location can be accessed */
  accessHours?: string;

  /** Environmental conditions monitoring */
  climate?: {
    /** Temperature in Celsius */
    temperature?: number;
    /** Relative humidity percentage */
    humidity?: number;
  };
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
