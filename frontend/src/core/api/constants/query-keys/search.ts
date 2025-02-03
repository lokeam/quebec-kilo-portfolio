/**
 * Parameters for media item search operations
 * @interface SearchParams
 */
interface SearchParams {
  /** Search query string to match against media items */
  query: string;

  /**
   * Type of media to filter results
   * @optional
   */
  type?: 'movie' | 'game' | 'book'; // Note: movie + game to be added post launch

  /**
   * Release year to filter results
   * @optional
   */
  year?: number;

  /**
   * Page number for paginated results
   * @optional
   * @default 1
   */
  page?: number;
}

/**
 * Tuple type representing a search query key
 * Used for React Query cache management
 */
type SearchKey = readonly ['search', SearchParams];

/**
 * Query key factory for search-related operations
 * Used to maintain consistent cache keys in React Query
 */
export const searchKeys = Object.freeze({
  /** Base key for all search-related queries */
  base: ['search'] as const,

  /**
   * Generates a unique query key for specific search parameters
   * @param params - Search parameters to include in the query key
   * @returns A tuple containing the base key and search parameters
   * @example
   * ```typescript
   * const queryKey = searchKeys.results({ query: 'Batman', type: 'movie' });
   * ```
   */
  results: (params: SearchParams): SearchKey => ['search', params]
});
