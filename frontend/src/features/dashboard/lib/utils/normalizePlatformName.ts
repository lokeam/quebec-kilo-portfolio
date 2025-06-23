import { chartConfig } from '@/features/dashboard/components/organisms/HomePage/ItemsByPlatformCard/itemsByPlatformChard.const';

export type PlatformInfo = {
  key: keyof typeof chartConfig;
  displayName: string;
};

/**
 * Takes a raw platform name from the API and returns a standardized object
 * containing a simplified key for config lookups and a clean display name.
 * @param rawName - The platform name string from the API.
 * @returns A PlatformInfo object.
 */
export function normalizePlatformName(rawName: string): PlatformInfo {
  const lowerCaseName = rawName.toLowerCase().trim();

  // PC
  if (lowerCaseName.includes('pc') || lowerCaseName.includes('windows')) {
    return { key: 'pc', displayName: 'PC' };
  }

  // PlayStation
  if (lowerCaseName.includes('playstation 4') || lowerCaseName.includes('ps4')) {
    return { key: 'ps4', displayName: 'PlayStation 4' };
  }
  if (lowerCaseName.includes('playstation 3') || lowerCaseName.includes('ps3')) {
    return { key: 'ps3', displayName: 'PlayStation 3' };
  }
  if (lowerCaseName.includes('playstation 2') || lowerCaseName.includes('ps2')) {
    return { key: 'ps2', displayName: 'PlayStation 2' };
  }

  // Xbox
  if (lowerCaseName.includes('xbox series')) {
    return { key: 'xboxseriesx', displayName: 'Xbox Series X/S' };
  }
  if (lowerCaseName.includes('xbox one')) {
    return { key: 'xboxone', displayName: 'Xbox One' };
  }
  if (lowerCaseName.includes('xbox 360')) {
    return { key: 'xbox360', displayName: 'Xbox 360' };
  }

  // Nintendo
  if (lowerCaseName.includes('switch')) {
    return { key: 'switch', displayName: 'Switch' };
  }

  // Others
  if (lowerCaseName.includes('mobile')) {
    return { key: 'mobile', displayName: 'Mobile' };
  }
  if (lowerCaseName.includes('arcade')) {
    return { key: 'arcade', displayName: 'Arcade' };
  }

  // Fallback for any unmatched platforms
  return { key: 'pc', displayName: 'Other' };
}