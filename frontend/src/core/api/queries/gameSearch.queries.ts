/**
 * Game Search Query Hooks
 *
 * Provides React Query hooks for searching games with various filters and sorting options.
 */

import { keepPreviousData } from '@tanstack/react-query';

// utils
import { logger } from '@/core/utils/logger/logger';

// services
import { searchGames } from '@/core/api/services/gameSearch.service';

// base query hook
import { useAPIQuery } from '@/core/api/queries/useAPIQuery';

// adapters
import { gameSearchResultAdapter } from '@/core/api/adapters/gameSearchResult.adapter';

// types
import type { SearchCriteria, SearchResult, SearchMetadata } from '@/types/domain/search';

// constants
import { QUERY_GARBAGE_COLLECTION_TIME, QUERY_STALE_TIME } from '@/core/api/config';

/**
 * Query key factory for game search queries
 */
export const gameSearchKeys = {
  all: ['game-search'] as const,
  lists: () => [...gameSearchKeys.all, 'list'] as const,
  list: (criteria: SearchCriteria) => [...gameSearchKeys.lists(), criteria] as const,
};

/**
 * Hook for searching games with React Query
 *
 * @function useGameSearch
 * @param {SearchCriteria} criteria - The search criteria including query, filters, and sorting options
 * @returns {UseQueryResult} React Query result object containing the search results and metadata
 *
 * @remarks
 * Responsibilities:
 * - Keeps previous data while loading new results
 * - Only runs when query has non-empty value
 * - Configurable stale time + garbage collection
 * - Handles pagination and filtering
 * - Transforms API response to domain types using the adapter
 *
 * @example
 * ```typescript
 * function SearchComponent() {
 *   const [searchTerm, setSearchTerm] = useState('');
 *   const { data, isLoading, error } = useGameSearch({
 *     query: searchTerm,
 *     filters: {
 *       platforms: ['Nintendo Switch'],
 *       rating: 4.5
 *     },
 *     sortBy: 'rating',
 *     sortOrder: 'desc'
 *   });
 *
 *   if (isLoading) return <div>Loading...</div>;
 *   if (error) return <div>Error: {error.message}</div>;
 *
 *   return (
 *     <div>
 *       {data?.results.map(result => (
 *         <GameCard key={result.game.id} game={result.game} />
 *       ))}
 *     </div>
 *   );
 * }
 * ```
 */
export function useGameSearch(criteria: SearchCriteria) {
  const hasQueryValue = criteria.query.trim().length > 0;

  return useAPIQuery<{ results: SearchResult[]; metadata: SearchMetadata }>({
    queryKey: gameSearchKeys.list(criteria),
    queryFn: async () => {
      try {
        // Get raw API response
        const response = await searchGames(criteria);
        logger.debug('🔍 Search succeeded - useGameSearch raw response:', {
          resultCount: response.games.length,
          total: response.total
        });

        // Transform API response to domain types using the adapter
        const transformedResponse = gameSearchResultAdapter.toDomain(response);
        logger.debug('🔍 Search succeeded - useGameSearch transformed response:', {
          resultCount: transformedResponse.results.length,
          metadata: transformedResponse.metadata
        });

        return transformedResponse;
      } catch (error) {
        logger.error('🔍 Search failed - useGameSearch error:', { error });
        throw error;
      }
    },
    enabled: hasQueryValue,
    placeholderData: keepPreviousData,
    staleTime: QUERY_STALE_TIME,
    gcTime: QUERY_GARBAGE_COLLECTION_TIME,
  });
}
