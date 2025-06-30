import type {
  PhysicalLocationResponse,
  DigitalLocationResponse
} from '@/types/domain/library-types';

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