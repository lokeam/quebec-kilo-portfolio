/**
 * Search Service API
 *
 * For API standards and best practices, see:
 * @see {@link ../../../docs/api-standards.md}
 */

import { axiosInstance } from '@/core/api/client/axios-instance';


// Type Refactor
import type { Game } from '@/types/game';
import type { SearchResponse } from '@/types/api';

// Legacy Types
//import type { SearchResponse } from '@/core/api/types/search.types';
//import type { WishlistItem } from '@/features/dashboard/lib/types/wishlist/base';

// Debug
import { logger } from '@/core/utils/logger/logger';

const IGDB_COVER_URL_BASE = 'https://images.igdb.com/igdb/image/upload/t_cover_big/';

/**
 * Searches for media items based on a query string
 *
 * @async
 * @function searchMediaItems
 * @param {string} query - The search query to find matching media items
 * @returns {Promise<Game[]>} A promise that resolves to an array of matching wishlist items
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
export const searchMediaItems = async (query: string): Promise<Game[]> => {
  logger.debug('🔍 Searching for media items', { query });

  try {
    const { data } = await axiosInstance.post<{ data: SearchResponse }>(
      '/v1/search',
      { query }
    );

    // Debug: Log the full response
    logger.debug('🔄 Full backend response:', { data });

    // Ensure the response has the required fields
    if (!data || typeof data !== 'object') {
      throw new Error(`Invalid response: Expected an object, got ${JSON.stringify(data)}`);
    }

    // The response should have a games array
    if (!data.games || !Array.isArray(data.games)) {
      throw new Error(`Invalid 'games' field: Expected an array, got ${JSON.stringify(data.games)}`);
    }

    // Transform the response to ensure proper image URLs
    const games = data.games.map((game: Game) => ({
      ...game,
      cover_url: game.coverUrl ? `${IGDB_COVER_URL_BASE}${game.coverUrl}` : ''
    }));

    return games;
  } catch (error) {
    logger.error('🚨 Search failed:', { error });
    throw error;
  }
};