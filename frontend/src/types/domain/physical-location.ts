import type { PhysicalLocationType, LocationIconBgColor } from './location-types';
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
  locationType: PhysicalLocationType;  // 'house', 'apartment', etc.
  description?: string;
  metadata?: PhysicalLocationMetadata;
  sublocations?: Sublocation[];
  createdAt: Date;
  updatedAt: Date;
  bgColor?: LocationIconBgColor;
  mapCoordinates?: MapCoordinates;
}

/**
 * Request type for creating a new physical location
 */
export interface CreatePhysicalLocationRequest {
  name: string;
  bgColor?: LocationIconBgColor;
  type: PhysicalLocationType;
  description?: string;
  locationType: PhysicalLocationType
  mapCoordinates?: string;
}

/* yet another fucking refactor ----- bff response types ------ */
export interface MapCoordinatesResponse {
  coords: string;
  googleMapsLink: string;
}

export interface LocationsBFFPhysicalLocationResponse {
  physicalLocationId: string;
  name: string;
  physicalLocationType: string;
  mapCoordinates: MapCoordinatesResponse;
  bgColor: string;
  createdAt: string; // or Date, depending on your usage
  updatedAt: string; // or Date
}

export interface LocationsBFFStoredGameResponse {
  id: string;
  name: string;
  platform: string;
  isUniqueCopy: boolean;
  hasDigitalCopy: boolean;
}

export interface LocationsBFFSublocationResponse {
  sublocationId: string;
  sublocationName: string;
  sublocationType: string;
  storedItems: number;
  parentLocationId: string;
  parentLocationName: string;
  parentLocationType: string;
  parentLocationBgColor: LocationIconBgColor;
  mapCoordinates: MapCoordinatesResponse;
  createdAt: string; // or Date
  updatedAt: string; // or Date
  storedGames: LocationsBFFStoredGameResponse[];
}

export interface LocationsBFFResponse {
  physicalLocations: LocationsBFFPhysicalLocationResponse[];
  sublocations: LocationsBFFSublocationResponse[];
}