import type { OnlineServiceProviderId } from "../constants/service.constants";
import type { StorageSize } from "./storage";

/**
 * Base location types shared across features
 * @see https://martinfowler.com/articles/modularizing-domain-models.html
 */

export interface BaseLocation {
  type: LocationCategory;
  name: string;
}

/**
 * Physical location types supported by the system
 * Used for categorizing where physical media is stored
 */
export type LocationCategory =
  | 'apartment'
  | 'house'
  | 'office'
  | 'warehouse'
  | 'vehicle';

/**
 * Sublocation types supported by the system
 * Used for categorizing where physical media is stored within a physical location
 */
export type StorageUnit =
  | 'box'
  | 'console'
  | 'cabinet'
  | 'closet'
  | 'drawer'
  | 'shelf';

/**
 * Complete information about a physical storage location
 * Example: { name: "Main Apartment", category: "apartment", unit: "shelf" }
 */
export interface PhysicalLocation {
  name: string;
  category: LocationCategory;
  subname?: string;
  sublocation?: StorageUnit;
}

/**
 * Complete information about a digital storage location
 * Example: { service: "Steam", diskSize: "50GB" }
 */
export interface DigitalLocation {
  service: OnlineServiceProviderId;
  diskSize?: StorageSize;
}
