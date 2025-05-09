/**
 * LOGO SYSTEM DOCUMENTATION
 * ========================
 *
 * This file handles icon/logo rendering for both digital and physical locations.
 *
 * HOW TO CHANGE A LOGO
 * --------------------
 * 1. For Digital Locations (e.g., Steam, PlayStation):
 *    - The logo is determined by the digital location's name
 *    - The name is normalized using getLogo() in service-utils.ts
 *    - The normalized name must match a key in LOGO_MAP.games
 *
 *    Example:
 *    - API returns: { name: "PlayStation Network" }
 *    - getLogo() normalizes to: "playstation"
 *    - LOGO_MAP.games must have: { playstation: PlayStationLogo }
 *
 * 2. For Physical Locations:
 *    - The logo is determined by the location's type
 *    - The type must match a key in LOGO_MAP.location
 *
 * DEBUGGING LOGOS
 * --------------
 * The function includes debug logging that shows:
 * - Original name/type
 * - Normalized name (for digital)
 * - Available logos in the map
 *
 * COMMON ISSUES
 * ------------
 * 1. Logo not showing up?
 *    - Check the console logs
 *    - Verify the name/type matches a key in LOGO_MAP
 *    - Check if getLogo() is normalizing correctly
 *
 * 2. Wrong logo showing?
 *    - Check if the name/type is being normalized correctly
 *    - Verify the LOGO_MAP key matches the normalized name
 *
 * 3. Default icon showing?
 *    - This means no matching logo was found
 *    - Check the console logs to see what was looked up
 *
 * RELATED FILES
 * ------------
 * - service-utils.ts: Contains getLogo() for name normalization
 * - useDomainMaps.ts: Contains LOGO_MAP with all available logos
 */

import type { PhysicalLocation } from '@/types/domain/physical-location';
import type { DigitalLocation } from '@/types/domain/digital-location';
import type { ComponentType } from 'react';
import { getLogo } from './service-utils';

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
 * Gets the appropriate icon component for a location
 *
 * Flow for Digital Locations:
 * 1. Takes the digital location's name
 * 2. Normalizes it using getLogo() (e.g., "PlayStation Network" -> "playstation")
 * 3. Looks up the normalized name in LOGO_MAP.games
 * 4. Returns the matching logo component or null
 *
 * Flow for Physical Locations:
 * 1. Takes the physical location's type
 * 2. Converts it to lowercase
 * 3. Looks it up in LOGO_MAP.location
 * 4. Returns the matching icon component or null
 */
export const getDigitalOrPhysicalLocationIcon = (
  location: PhysicalLocation | DigitalLocation,
  type: 'physical' | 'digital',
  domainMaps: DomainMapsResult,
) => {
  const { games, location: locationIcons } = domainMaps;

  if (type === 'physical') {
    const physicalLocation = location as PhysicalLocation;
    const IconComponent = locationIcons[physicalLocation.type?.toLowerCase() || ''];
    console.log('Physical Location Icon Debug:', {
      originalName: physicalLocation.name,
      type: physicalLocation.type,
      availableIcons: Object.keys(locationIcons),
    });
    return IconComponent ? <IconComponent className="h-4 w-4" /> : null;
  } else {
    const digitalLocation = location as DigitalLocation;
    const normalizedName = getLogo(digitalLocation.name);
    console.log('Digital Location Logo Debug:', {
      originalName: digitalLocation.name,
      normalizedName,
      availableLogos: Object.keys(games),
    });
    const LogoComponent = normalizedName ? games[normalizedName] : null;
    return LogoComponent ? <LogoComponent className="h-4 w-4" /> : null;
  }
}

// Add a new function specifically for handling location types
export const getLocationTypeIcon = (
  locationType: string | undefined,
  domainMaps: DomainMapsResult
) => {
  const { location: locationIcons } = domainMaps;
  if (!locationType) return null;
  const IconComponent = locationIcons[locationType.toLowerCase()];
  return IconComponent ? <IconComponent className="h-4 w-4" /> : null;
}
