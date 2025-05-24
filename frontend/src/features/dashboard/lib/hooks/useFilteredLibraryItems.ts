import { useMemo } from 'react';
import type { LibraryGameItem } from '@/types/domain/library-types';

/**
 * Interface for the transformed library item structure used in the UI
 */
type TransformedLibraryItem = LibraryGameItem;

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

    return filtered;
  }, [services, platformFilter, searchQuery]);
}
