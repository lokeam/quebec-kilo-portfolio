import { useMemo } from 'react';
import { useDebounce } from '@/shared/hooks/useDebounce';

// Hooks
import { useLibraryGames } from '@/features/dashboard/lib/stores/libraryStore';

// Types + Constants
import { type PlatformOption, CONSOLE_PLATFORMS } from '@/features/dashboard/lib/constants/filterOptions/library/platform.filterOptions';
import { type LibraryItem } from '@/features/dashboard/lib/types/page.types';


interface PlatformWithSearch extends PlatformOption {
  searchString: string;
};

interface UsePlatformSearchReturn {
  availablePlatforms: Record<string, PlatformWithSearch[]>;
  handleSearch: (value: string) => void;
}

export function usePlatformSearch(): UsePlatformSearchReturn {
  const userGames = useLibraryGames();

  // Memoize platform keys
  const platformKeys = useMemo(() =>
    new Set(userGames.map((game: LibraryItem) => (game.platformVersion ?? '').toLowerCase())),
    [userGames]
  );

  // Memoize available platforms with pre-computed search strings
  const availablePlatforms = useMemo(() => {
    return Object.fromEntries(
      Object.entries(CONSOLE_PLATFORMS)
        .map(([manufacturer, platforms]) => [
          manufacturer,
          platforms
            .filter(platform => platformKeys.has(platform.key.toLowerCase()))
            .map(platform => ({
              ...platform,
              searchString: [
                platform.key,
                platform.label.toLowerCase(),
                ...platform.searchTerms
              ].join(' ')
            }))
        ])
        .filter(([_, platforms]) => platforms.length > 0)
    ) as Record<string, PlatformWithSearch[]>;
  }, [platformKeys]);

  // Debounce search to prevent unnecessary re-renders
  const handleSearch = useDebounce(
    (value: string) => value?.toLowerCase() ?? '',
    350
  );

  return {
    availablePlatforms,
    handleSearch
  };
}