import type { PhysicalLocationType } from './location-types';
import type { Sublocation } from './sublocation';

/**
 * Represents a physical location where media items are stored
 */
export interface PhysicalLocationMetadata {
  address?: string;
  room?: string;
  notes?: string;
}

export interface MapCoordinates {
  coords: string;
  googleMapsLink: string;
}

export interface PhysicalLocation {
  id: string;
  name: string;
  type: PhysicalLocationType;
  description?: string;
  metadata?: PhysicalLocationMetadata;
  sublocations?: Sublocation[];
  createdAt: Date;
  updatedAt: Date;
  bgColor?: string;
  mapCoordinates?: MapCoordinates;
}

/**
 * Request type for creating a new physical location
 */
export interface CreatePhysicalLocationRequest {
  name: string;
  type: PhysicalLocationType;
  description?: string;
  metadata?: PhysicalLocationMetadata;
}