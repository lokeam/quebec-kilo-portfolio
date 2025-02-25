import { axiosInstance } from '@/core/api/client/axios-instance';
import type { SearchResponse, WishlistItem } from '@/features/dashboard/lib/types/wishlist/base';
import { API_ROUTES } from '@/core/api/constants/routes';

// Debug
import { logger } from '@/core/utils/logger/logger';

/**
 * Searches for media items based on a query string
 *
 * @async
 * @function searchMediaItems
 * @param {string} query - The search query to find matching media items
 * @returns {Promise<WishlistItem[]>} A promise that resolves to an array of matching wishlist items
 *
 * @throws {Error} If the API request fails or returns an invalid response
 *
 * @example
 * ```typescript
 * try {
 *   const items = await searchMediaItems('Harry Potter');
 *   console.log(items);
 * } catch (error) {
 *   console.error('Search failed:', error);
 * }
 * ```
 */
export const searchMediaItems = async (query: string): Promise<WishlistItem[]> => {
  logger.debug('🔍 Searching for media items', { query });

  try {
    const response = await axiosInstance.post<SearchResponse>(
      '/api/v1/search',
      { query }
    ) as unknown as SearchResponse;

    // Debug: Log the full response
    logger.debug('🔄 Full backend response:', { response });

    // Ensure the response has the required fields
    if (!response || typeof response !== 'object') {
      throw new Error(`Invalid response: Expected an object, got ${JSON.stringify(response)}`);
    }

    // Since your Axios interceptor returns response.data, 'response' is already a SearchResponse
    if (!response.games || !Array.isArray(response.games)) {
      throw new Error(`Invalid 'games' field: Expected an array, got ${JSON.stringify(response.games)}`);
    }

    return response.games; // Return the games array
  } catch (error) {
    logger.error('🚨 Search failed:', { error });
    throw error;
  }
};