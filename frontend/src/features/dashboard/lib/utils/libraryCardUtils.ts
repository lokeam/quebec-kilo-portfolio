import type {
  PhysicalLocationResponse,
  DigitalLocationResponse
} from '@/types/domain/library-types';

/**
 * Data structure for displaying physical location information
 */
export interface PhysicalLocationDisplayData {
  parentLocationName: string;
  parentLocationType: string;
  parentLocationBgColor: string;
  sublocationName: string;
  sublocationType: string;
  platforms: string[];
}

/**
 * Data structure for displaying digital location information
 */
export interface DigitalLocationDisplayData {
  digitalLocationName: string;
  platforms: string[];
}

/**
 * Formats a release date from milliseconds to a readable format
 * @param releaseDate - Release date in milliseconds (Unix timestamp)
 * @returns Formatted date string (e.g., "January 15, 2023")
 */
export function formatReleaseDate(releaseDate: number): string {
  if (!releaseDate || releaseDate <= 0) {
    return 'Release date unknown';
  }

  try {
    const date = new Date(releaseDate * 1000); // Convert seconds to milliseconds
    return date.toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric'
    });
  } catch {
    return 'Release date unknown';
  }
}

/**
 * Removes duplicate platform names and returns unique array
 * @param platforms - Array of platform names
 * @returns Array of unique platform names
 */
export function getUniquePlatformNames(platforms: string[]): string[] {
  return [...new Set(platforms)];
}

/**
 * Extracts platform names from physical locations only
 * @param physicalLocations - Array of physical locations
 * @returns Array of platform names from physical locations
 */
export function extractPhysicalPlatformNamesFromLocations(
  physicalLocations: PhysicalLocationResponse[]
): string[] {
  return getUniquePlatformNames(
    physicalLocations.flatMap(location =>
      location.gamePlatformVersions.map(platform => platform.platformName)
    )
  );
}

/**
 * Extracts platform names from digital locations only
 * @param digitalLocations - Array of digital locations
 * @returns Array of platform names from digital locations
 */
export function extractDigitalPlatformNamesFromLocations(
  digitalLocations: DigitalLocationResponse[]
): string[] {
  return getUniquePlatformNames(
    digitalLocations.flatMap(location =>
      location.gamePlatformVersions.map(platform => platform.platformName)
    )
  );
}

/**
 * Extracts physical location data for display
 * @param physicalLocations - Array of physical locations
 * @returns Array of physical location display data
 */
export function extractPhysicalLocationData(
  physicalLocations: PhysicalLocationResponse[]
): PhysicalLocationDisplayData[] {
  return physicalLocations.map(location => ({
    parentLocationName: location.parentLocationName,
    parentLocationType: location.parentLocationType,
    parentLocationBgColor: location.parentLocationBgColor,
    sublocationName: location.sublocationName,
    sublocationType: location.sublocationType,
    platforms: location.gamePlatformVersions.map(p => p.platformName)
  }));
}

/**
 * Extracts digital location data for display
 * @param digitalLocations - Array of digital locations
 * @returns Array of digital location display data
 */
export function extractDigitalLocationData(
  digitalLocations: DigitalLocationResponse[]
): DigitalLocationDisplayData[] {
  return digitalLocations.map(location => ({
    digitalLocationName: location.digitalLocationName,
    platforms: location.gamePlatformVersions.map(p => p.platformName)
  }));
}

/**
 * Formats platform names for display (handles edge cases)
 * @param platforms - Array of platform names
 * @param maxDisplayCount - Maximum number of platforms to show before truncating
 * @returns Object with display platforms and overflow count
 */
export function formatPlatformsForDisplay(
  platforms: string[],
  maxDisplayCount: number = 3
): { displayPlatforms: string[]; overflowCount: number } {
  const uniquePlatforms = getUniquePlatformNames(platforms);
  const displayPlatforms = uniquePlatforms.slice(0, maxDisplayCount);
  const overflowCount = Math.max(0, uniquePlatforms.length - maxDisplayCount);

  return { displayPlatforms, overflowCount };
}