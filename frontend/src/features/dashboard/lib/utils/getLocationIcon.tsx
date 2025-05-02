import type { PhysicalLocation } from '@/features/dashboard/lib/types/media-storage/physical';
import type { DigitalLocation } from '@/features/dashboard/lib/types/media-storage/digital';
import type { ComponentType } from 'react';

type DomainMapsResult = {
  games: Record<string, ComponentType<{ className?: string }>>;
  location: Record<string, ComponentType<{ className?: string }>>;
  // NOTE: Add other properties as needed
  movies: Record<string, ComponentType<{ className?: string }>>;
  music: Record<string, ComponentType<{ className?: string }>>;
  platforms: Record<string, ComponentType<{ className?: string }>>;
  sublocation: Record<string, ComponentType<{ className?: string }>>;
  physicalMedia: Record<string, ComponentType<{ className?: string }>>;
  digitalMedia: Record<string, ComponentType<{ className?: string }>>;
  misc: Record<string, ComponentType<{ className?: string }>>;
  notifications: Record<string, ComponentType<{ className?: string }>>;
}

/**
 * Type definition that allows both camelCase and snake_case property access
 * during the API standardization transition period
 */
interface PhysicalLocationWithSnakeCase extends PhysicalLocation {
  location_type?: string;
  map_coordinates?: string;
  sub_locations?: Array<{
    id: string;
    name: string;
    description?: string;
    [key: string]: unknown;
  }>;
  created_at?: string;
  updated_at?: string;
  user_id?: string;
}

/**
 * Gets the appropriate icon component for a location
 *
 * IMPORTANT: This function handles both snake_case and camelCase properties
 * for compatibility during the API standardization transition.
 */
export const getLocationIcon = (
  location: PhysicalLocation | DigitalLocation,
  type: 'physical' | 'digital',
  domainMaps: DomainMapsResult,
) => {
  const { games, location: locationIcons } = domainMaps;

  if (type === 'physical') {
    const physicalLocation = location as PhysicalLocationWithSnakeCase;

    // Handle both snake_case and camelCase properties
    const locationType =
      // Try camelCase first (preferred)
      physicalLocation.locationType ||
      // Fall back to snake_case (legacy)
      physicalLocation.location_type;

    const IconComponent = locationIcons[locationType?.toLowerCase() || ''];
    return IconComponent ? <IconComponent className="h-4 w-4" /> : null;
  } else {
    const digitalLocation = location as DigitalLocation;
    const LogoComponent = games[digitalLocation.label?.toLowerCase() || ''];
    return LogoComponent ? <LogoComponent className="h-4 w-4" /> : null;
  }
}

// Add a new function specifically for handling location types
export const getLocationTypeIcon = (
  locationType: string,
  domainMaps: DomainMapsResult
) => {
  const { location: locationIcons } = domainMaps;
  const IconComponent = locationIcons[locationType.toLowerCase()];
  return IconComponent ? <IconComponent className="h-4 w-4 mr-1" /> : null;
}
