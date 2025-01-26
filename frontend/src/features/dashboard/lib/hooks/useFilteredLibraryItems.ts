import { useMemo } from 'react';
import { type LibraryItem } from '@/features/dashboard/lib/types/library/items';
import { isPhysicalLibraryItem } from '@/features/dashboard/lib/types/library/guards';

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
  services: LibraryItem[],
  platformFilter: string | null,
  searchQuery: string | null
): TransformedLibraryItem[] {
  return useMemo(() => {
    let filtered = services;

    /* Dropdown Filter by platform */
    if (platformFilter) {
      filtered = filtered.filter(game =>
        game.platform?.category?.toLowerCase() === platformFilter.toLowerCase()
      );
    }

    /* Search Bar Filter by title */
    if (searchQuery) {
      filtered = filtered.filter(game =>
        game.title.toLowerCase().includes(searchQuery.toLowerCase())
      );
    }

    return filtered.map((item, index) => ({
        index,
        title: item.title,
        imageUrl: item.imageUrl,
        favorite: item.favorite,
        platform: item.platform.category,
        ...(isPhysicalLibraryItem(item)
          ? {
              physicalLocation: item.location.name,
              physicalLocationType: item.location.category,
              physicalSublocation: item.location.subname,
              physicalSublocationType: item.location.sublocation,
              type: item.type
          } : {
              digitalLocation: item.location.service,
              diskSize: `${item.location.diskSize.value} ${item.location.diskSize.unit}`,
              type: item.type
          }
        )
    }));
  }, [services, platformFilter, searchQuery]);
}
