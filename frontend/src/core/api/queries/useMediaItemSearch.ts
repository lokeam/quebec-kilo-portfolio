import { keepPreviousData } from '@tanstack/react-query';
import { useBackendQuery } from '@/core/api/hooks/useBackendQuery';
import { searchKeys } from '@/core/api/constants/query-keys/search';
import { searchMediaItems } from '@/core/api/services/search.service';
import { QUERY_GARBAGE_COLLECTION_TIME, QUERY_STALE_TIME } from '@/core/api/config';

/**
 * Custom hook for searching media items with React Query
 *
 * @function useMediaItemSearch
 * @param {string} query - The search query string to find media items
 * @returns {UseQueryResult} React Query result object containing the search results + status
 *
 * @remarks
 * Responsibilities:
 * - Keeps previous data while loading new results
 * - Only runs when query has non-empty value
 * - Configurable stale time + garbage collection
 *
 * @example
 * ```typescript
 * function SearchComponent() {
 *   const [searchTerm, setSearchTerm] = useState('');
 *   const { data, isLoading, error } = useMediaItemSearch(searchTerm);
 *
 *   if (isLoading) return <div>Loading...</div>;
 *   if (error) return <div>Error: {error.message}</div>;
 *
 *   return <div>{data?.map(item => <MediaItem key={item.id} {...item} />)}</div>;
 * }
 * ```
 */
export function useMediaItemSearch(query: string) {
  const hasQueryValue = query.trim().length > 0;

  return useBackendQuery({
    queryKey: searchKeys.results({ query }),
    queryFn: () => searchMediaItems(query),
    enabled: hasQueryValue,
    placeholderData: keepPreviousData,
    staleTime: QUERY_STALE_TIME,
    gcTime: QUERY_GARBAGE_COLLECTION_TIME,
  });
}