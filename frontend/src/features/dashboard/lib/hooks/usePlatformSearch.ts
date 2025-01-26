import { useMemo } from 'react';
import { useLibraryGames } from '@/features/dashboard/lib/stores/libraryStore';
import { type PlatformOption, CONSOLE_PLATFORMS } from '@/features/dashboard/lib/constants/filter-options/library/platform.filterOptions';
import { type LibraryItem } from '@/features/dashboard/lib/types/library/items';
import { useMultiMatchSearch } from '@/features/dashboard/lib/hooks/useMultiMatchSearch';

/**
 * Enhanced platform type that includes search metadata
 * @internal
 */
interface PlatformWithSearch extends PlatformOption {
  /** Combined searchable string for the platform */
  searchString: string;
  /** Platform manufacturer (e.g., 'sony', 'microsoft') */
  manufacturer: string;
}

/**
 * Return type for the usePlatformSearch hook
 */
interface UsePlatformSearchReturn {
  /**
   * Platforms grouped by manufacturer, filtered by search term
   * Only includes platforms that exist in user's game library
   */
  availablePlatforms: Record<string, Omit<PlatformWithSearch, 'manufacturer'>[]>;
  /** Function to update the platform search term */
  handleSearch: (value: string) => void;
}

/**
 * Creates a searchable string for a platform by combining its key, label, and search terms
 *
 * @param platform - Platform to create search string for
 * @returns Combined, lowercase, space-separated search string
 *
 * @example
 * ```ts
 * const searchString = createPlatformSearchString({
 *   key: 'ps4',
 *   label: 'PlayStation 4',
 *   searchTerms: ['sony', 'console']
 * });
 * // Returns: "ps4 playstation 4 sony console"
 * ```
 *
 * @internal
 */
function createPlatformSearchString(platform: PlatformOption): string {
  const searchParts = [
    platform.key,
    platform.label,
    ...(platform.searchTerms || [])
  ].map(term => term.toLowerCase());

  const searchString = [...new Set(searchParts)].join(' ');
  console.log('Created search string for platform:', {
    key: platform.key,
    label: platform.label,
    searchTerms: platform.searchTerms,
    finalSearchString: searchString
  });
  return searchString;
}

/**
 * Hook that provides searchable platform filtering functionality for the library
 *
 * This hook is specifically designed for the platform filter in the library view.
 * It combines data from:
 * 1. The user's game library (to determine available platforms)
 * 2. The platform configuration (for display names and search terms)
 * 3. Search functionality (to filter platforms based on user input)
 *
 * Key features:
 * - Only shows platforms that exist in the user's library
 * - Provides fuzzy search across platform names and aliases
 * - Groups platforms by manufacturer
 * - Memoizes results for performance
 *
 * @returns {UsePlatformSearchReturn} Object containing:
 *   - availablePlatforms: Platforms grouped by manufacturer
 *   - handleSearch: Function to update search term
 *
 * @example
 * ```tsx
 * // In PlatformCombobox.tsx
 * function PlatformCombobox() {
 *   const { availablePlatforms, handleSearch } = usePlatformSearch();
 *
 *   return (
 *     <Command.List>
 *       {Object.entries(availablePlatforms).map(([manufacturer, platforms]) => (
 *         <Command.Group key={manufacturer} heading={manufacturer}>
 *           {platforms.map(platform => (
 *             <Command.Item key={platform.key} value={platform.key}>
 *               {platform.label}
 *             </Command.Item>
 *           ))}
 *         </Command.Group>
 *       ))}
 *     </Command.List>
 *   );
 * }
 * ```
 *
 * @see {@link PlatformOption} for platform data structure
 * @see {@link useMultiMatchSearch} for search functionality
 * @see {@link useLibraryGames} for game library access
 */
export function usePlatformSearch(): UsePlatformSearchReturn {
  const userGames = useLibraryGames();

  /* Get unique platform keys from user's game library */
  const platformKeys = useMemo(() =>
    new Set(userGames.map((game: LibraryItem) =>
      game.platform.category
    )),
    [userGames]
  );

  /* Create searchable objects for game console platforms in user's library */
  const enhancedPlatforms = useMemo(() => {
    const platforms = Object.entries(CONSOLE_PLATFORMS)
      .flatMap(([manufacturer, platforms]) =>
        platforms
          .filter(platform => platformKeys.has(platform.key))
          .map(platform => ({
            ...platform,
            manufacturer,
            searchString: createPlatformSearchString(platform)
          }))
      );

      return platforms;
  }, [platformKeys]);

  /* Apply search functionality using useMultiMatchSearch hook */
  const { filteredItems, setSearchTerm } = useMultiMatchSearch(enhancedPlatforms, {
    minFuzzyLength: 3,
    caseSensitive: false
  });


  /* Group filtered results by manufacturer */
  const availablePlatforms = useMemo(() => {
    const grouped = filteredItems.reduce<Record<string, Omit<PlatformWithSearch, 'manufacturer'>[]>>(
      (acc, platform) => {
        const { manufacturer, ...platformData } = platform;

        if (!acc[manufacturer]) {
          acc[manufacturer] = [];
        }

        acc[manufacturer].push(platformData);
        return acc;
      },
      {}
    );

    return grouped;
  }, [filteredItems]);

  return {
    availablePlatforms,
    handleSearch: setSearchTerm,
  };
}
