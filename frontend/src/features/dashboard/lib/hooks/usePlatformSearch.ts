import { useMemo } from 'react';
import { useLibraryGames } from '@/features/dashboard/lib/stores/libraryStore';
import { type PlatformOption, CONSOLE_PLATFORMS } from '@/features/dashboard/lib/constants/filterOptions/library/platform.filterOptions';
import { type LibraryItem } from '@/features/dashboard/lib/types/page.types';
import { useMultiMatchSearch } from '@/features/dashboard/lib/hooks/useMultiMatchSearch';

interface PlatformWithSearch extends PlatformOption {
  searchString: string;
  manufacturer: string;
}

interface UsePlatformSearchReturn {
  availablePlatforms: Record<string, Omit<PlatformWithSearch, 'manufacturer'>[]>;
  handleSearch: (value: string) => void;
}

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


export function usePlatformSearch(): UsePlatformSearchReturn {
  const userGames = useLibraryGames();

  // Get platform keys from user's games
  const platformKeys = useMemo(() =>
    new Set(userGames.map((game: LibraryItem) =>
      (game.platformVersion ?? '').toLowerCase()
    )),
    [userGames]
  );

  /* Create enhanced platforms with search strings */
  const enhancedPlatforms = useMemo(() => {
    const platforms = Object.entries(CONSOLE_PLATFORMS)
      .flatMap(([manufacturer, platforms]) =>
        platforms
          .filter(platform => platformKeys.has(platform.key.toLowerCase()))
          .map(platform => ({
            ...platform,
            manufacturer,
            searchString: createPlatformSearchString(platform)
          }))
      );

      return platforms;
  }, [platformKeys]);

  /* Use generic search hook */
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