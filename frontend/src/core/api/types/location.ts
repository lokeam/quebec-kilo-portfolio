// Base interface for all location types
export interface BaseLocation {
  id?: string;
  name: string;
  locationType: string;
  mapCoordinates?: string;
}

// Physical location is just a BaseLocation with no additional fields
export type PhysicalLocation = BaseLocation;

// Sublocation specific interface
export interface Sublocation extends BaseLocation {
  parentLocationId: string;
  bgColor?: string;
}

// Type guard to check if a location is a physical location
export function isPhysicalLocation(location: BaseLocation): location is PhysicalLocation {
  return !('parentLocationId' in location);
}

// Type guard to check if a location is a sublocation
export function isSublocation(location: BaseLocation): location is Sublocation {
  return 'parentLocationId' in location;
}