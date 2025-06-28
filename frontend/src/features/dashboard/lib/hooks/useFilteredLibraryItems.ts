import { useMemo } from 'react';
import type { LibraryGameItem, LibraryGameItemResponse } from '@/types/domain/library-types';

/**
 * Interface for the transformed library item structure used in the UI
 */
type TransformedLibraryItem = LibraryGameItemResponse;

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
  services: LibraryGameItemResponse[],
  platformFilter: string | null,
  searchQuery: string | null
): TransformedLibraryItem[] {
  return useMemo(() => {
    let filtered = services;

    /* Dropdown Filter by platform */
    if (platformFilter) {
      filtered = filtered.filter(game =>
        game.gamesByPlatformAndLocation.some(loc =>
          getPlatformCategory(loc.platformName) === platformFilter
        )
      );
    }

    /* Search Bar Filter by title */
    if (searchQuery) {
      filtered = filtered.filter(game =>
        game.name.toLowerCase().includes(searchQuery.toLowerCase())
      );
    }

    return filtered;
  }, [services, platformFilter, searchQuery]);
}

function getPlatformCategory(platformName: string): string {
  const lowerName = platformName.toLowerCase();

  if (lowerName.includes('playstation') ||
      lowerName.includes('xbox') ||
      lowerName.includes('nintendo') ||
      lowerName.includes('sega') ||
      lowerName.includes('atari') ||
      lowerName.includes('nec')) {
    return 'Console';
  }

  if (lowerName.includes('android') ||
      lowerName.includes('ios') ||
      lowerName.includes('mobile')) {
    return 'Mobile';
  }

  if (lowerName.includes('pc') ||
      lowerName.includes('windows') ||
      lowerName.includes('mac')) {
    return 'PC';
  }

  return 'Console'; // Default fallback
}