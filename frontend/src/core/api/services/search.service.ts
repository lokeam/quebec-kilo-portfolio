import { axiosInstance } from '@/core/api/client/axios-instance';
import type { WishlistItem } from '@/features/dashboard/lib/types/wishlist/base';
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
    const { data } = await axiosInstance.get<WishlistItem[]>(
      API_ROUTES.SEARCH.GAMES,
      { params: { query } }
    );

    return data;
  } catch (error) {
    logger.error('🚨 Search failed:', { error });

    throw error;
  }
};