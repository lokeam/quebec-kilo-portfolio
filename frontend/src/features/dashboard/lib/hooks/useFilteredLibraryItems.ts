import { useMemo } from 'react';
import type { LibraryGameItem } from '@/types/domain/library-types';

/**
 * Interface for the transformed library item structure used in the UI
 */
interface TransformedLibraryItem {
  index: number;
  title: string;
  imageUrl: string;
  favorite: boolean;
  platform: string;
  type: 'physical' | 'digital';
  physicalLocation?: string;
  physicalLocationType?: string;
  physicalSublocation?: string;
  physicalSublocationType?: string;
  digitalLocation?: string;
  diskSize?: string;
}

/**
 * Custom hook that filters and transforms library items based on platform and search criteria
 *
 * @param services - Array of library items to filter
 * @param platformFilter - Optional platform category to filter by
 * @param searchQuery - Optional search term to filter titles
 * @returns Array of transformed library items for UI display
 *
 * @example
 * ```tsx
 * function LibraryView() {
 *   const items = useLibraryGames();
 *   const platformFilter = useLibraryPlatformFilter();
 *   const searchQuery = useLibrarySearchQuery();
 *
 *   const filteredItems = useFilteredLibraryItems(
 *     items,
 *     platformFilter,
 *     searchQuery
 *   );
 *
 *   return filteredItems.map(item => <LibraryItem {...item} />);
 * }
 * ```
 */
export function useFilteredLibraryItems(
  services: LibraryGameItem[],
  platformFilter: string | null,
  searchQuery: string | null
): TransformedLibraryItem[] {
  return useMemo(() => {
    let filtered = services;

    /* Dropdown Filter by platform */
    if (platformFilter) {
      filtered = filtered.filter(game =>
        game.gamesByPlatformAndLocation.some(loc =>
          loc.platformName.toLowerCase() === platformFilter.toLowerCase()
        )
      );
    }

    /* Search Bar Filter by title */
    if (searchQuery) {
      filtered = filtered.filter(game =>
        game.name.toLowerCase().includes(searchQuery.toLowerCase())
      );
    }

    return filtered.map((item, index) => {
      // Get the first platform location for display
      const firstLocation = item.gamesByPlatformAndLocation[0];

      return {
        index,
        title: item.name,
        imageUrl: item.coverUrl,
        favorite: item.favorite,
        platform: firstLocation.platformName,
        type: firstLocation.type,
        ...(firstLocation.type === 'physical'
          ? {
              physicalLocation: firstLocation.locationName,
              physicalLocationType: firstLocation.locationType,
              physicalSublocation: firstLocation.sublocationName,
              physicalSublocationType: firstLocation.sublocationType,
            }
          : {
              digitalLocation: firstLocation.locationName,
              // Note: diskSize is not available in the new API response
              // We'll need to add this if it's required
            }
        )
      };
    });
  }, [services, platformFilter, searchQuery]);
}
