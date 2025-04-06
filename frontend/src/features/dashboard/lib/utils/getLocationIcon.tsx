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

export const getLocationIcon = (
  location: PhysicalLocation | DigitalLocation,
  type: 'physical' | 'digital',
  domainMaps: DomainMapsResult,
) => {
  const { games, location: locationIcons } = domainMaps;

  if (type === 'physical') {
    const physicalLocation = location as PhysicalLocation;
    const IconComponent = locationIcons[physicalLocation.locationType];
    return IconComponent ? <IconComponent className="h-4 w-4" /> : null;
  } else {
    const digitalLocation = location as DigitalLocation;
    const LogoComponent = games[digitalLocation.label];
    return LogoComponent ? <LogoComponent className="h-4 w-4" /> : null;
  }
}

// Add a new function specifically for handling location types
export const getLocationTypeIcon = (
  locationType: string,
  domainMaps: DomainMapsResult
) => {
  const { location: locationIcons } = domainMaps;
  const IconComponent = locationIcons[locationType];
  return IconComponent ? <IconComponent className="h-4 w-4 mr-1" /> : null;
}
