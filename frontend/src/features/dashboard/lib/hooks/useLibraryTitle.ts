import { useMemo } from 'react';
import { CONSOLE_PLATFORMS } from '@/features/dashboard/lib/constants/filter-options/library/platform.filterOptions';
/**
 * Interface for useLibraryTitle hook options
 * @interface UseLibraryTitleOptions
 */
interface UseLibraryTitleOptions {
  /** Base title to display when no platform filter is active */
  baseTitle: string;
  /** Number of games that match the current filter */
  filteredCount: number;
  /** Current platform filter value (e.g., 'pc', 'ps4') or null if no filter */
  platformFilter: string | null;
}

/**
 * Hook that generates formatted title and count text for the library view
 *
 * @param {Object} options - The options object
 * @param {string} options.baseTitle - Base title to display when no platform filter is active (defaults to 'All Games')
 * @param {number} options.filteredCount - Number of games that match the current filter
 * @param {string|null} options.platformFilter - Current platform filter value or null if no filter
 *
 * @returns {Object} An object containing:
 *   - title: The formatted title (e.g., "All PC Games" or "All PlayStation 4 Games")
 *   - countText: The formatted count text (e.g., "(5 games)" or "(1 game)")
 *
 * @example
 * ```tsx
 * const { title, countText } = useLibraryTitle({
 *   baseTitle: "All Games",
 *   filteredCount: filteredGames.length,
 *   platformFilter: currentPlatform
 * });
 *
 * return (
 *   <h1>
 *     {title}
 *     <span>{countText}</span>
 *   </h1>
 * );
 * ```
 */
export function useLibraryTitle({
  baseTitle = 'All Games',
  filteredCount,
  platformFilter
}: UseLibraryTitleOptions) {
  /* Generate title based on platform filter */
  const title = useMemo(() => {
    if (!platformFilter) return baseTitle;

    const platformLabel = Object.values(CONSOLE_PLATFORMS)
      .flatMap(platforms => platforms)
      .find(platform => platform.key.toLowerCase() === platformFilter.toLowerCase())
      ?.label ?? platformFilter;

      return `All ${platformLabel} Games`;
  }, [baseTitle, platformFilter]);

  /* Generate count text with proper pluralization */
  const countText = useMemo(() => {
    return `(${filteredCount} ${filteredCount === 1 ? 'game' : 'games'})`;
  }, [filteredCount]);

  return { title, countText };
};
