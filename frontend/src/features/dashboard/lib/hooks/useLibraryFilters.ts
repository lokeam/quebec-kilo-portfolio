import { useMemo } from 'react';
import type { LibraryGameItemRefactoredResponse } from '@/types/domain/library-types';

interface FilterOption {
  key: string;
  label: string;
}

interface LibraryFilterOptions {
  platforms: FilterOption[];
  locations: FilterOption[];
}

export function useLibraryFilters(libraryItems: LibraryGameItemRefactoredResponse[]): LibraryFilterOptions {
  return useMemo(() => {
    if (!libraryItems || libraryItems.length === 0) {
      return { platforms: [], locations: [] };
    }

    // Extract unique platforms from library items
    const uniquePlatforms = Array.from(new Set(
      libraryItems.flatMap(item => [
        // Extract platforms from physical locations
        ...(item.physicalLocations || []).flatMap(location =>
          (location.gamePlatformVersions || []).map(platform => platform.platformName)
        ),
        // Extract platforms from digital locations
        ...(item.digitalLocations || []).flatMap(location =>
          (location.gamePlatformVersions || []).map(platform => platform.platformName)
        )
      ])
    ))
    .sort()
    .map(platform => ({
      key: platform,
      label: platform
    }));

    // Extract unique locations from library items
    const locationSet = new Set<string>();

    libraryItems.forEach(item => {
      // Add physical locations
      (item.physicalLocations || []).forEach(location => {
        // Add sublocation name if it exists
        if (location.sublocationName) {
          locationSet.add(location.sublocationName);
        }
        // Add parent location name if it exists
        if (location.parentLocationName) {
          locationSet.add(location.parentLocationName);
        }
      });

      // Add digital locations
      (item.digitalLocations || []).forEach(location => {
        if (location.digitalLocationName) {
          locationSet.add(location.digitalLocationName);
        }
      });
    });

    const uniqueLocations = Array.from(locationSet)
      .sort()
      .map(location => ({
        key: location,
        label: location
      }));

    return {
      platforms: uniquePlatforms,
      locations: uniqueLocations
    };
  }, [libraryItems]);
}