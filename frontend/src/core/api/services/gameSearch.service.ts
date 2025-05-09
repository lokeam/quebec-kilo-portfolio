/**
 * Game Search API Service
 *
 * This service handles API requests related to game search functionality.
 * It provides a unified interface for searching games with various filters and sorting options.
 */

import { axiosInstance } from '@/core/api/client/axios-instance';
import { logger } from '@/core/utils/logger/logger';
import type { SearchResponse } from '@/types/api/search';
import type { SearchCriteria } from '@/types/domain/search';

// Standard response type for all API calls
interface ApiResponse<T> {
  success: boolean;
  data: T;
  metadata: {
    timestamp: string;
    request_id: string;
  };
}

/**
 * Searches for games based on the provided criteria
 *
 * @async
 * @function searchGames
 * @param {SearchCriteria} criteria - The search criteria including query, filters, and sorting options
 * @returns {Promise<SearchResponse>} A promise that resolves to the raw API response
 * @throws {Error} If the API request fails or returns an invalid response
 *
 * @example
 * ```typescript
 * try {
 *   const response = await searchGames({
 *     query: 'Zelda',
 *     filters: {
 *       platforms: ['Nintendo Switch'],
 *       rating: 4.5
 *     },
 *     sortBy: 'rating',
 *     sortOrder: 'desc'
 *   });
 *   console.log(response.games);
 * } catch (error) {
 *   console.error('Search failed:', error);
 * }
 * ```
 */
export const searchGames = async (
  criteria: SearchCriteria
): Promise<SearchResponse> => {
  logger.debug('🔍 Searching for games', { criteria });

  try {
    const response = await axiosInstance.post<ApiResponse<SearchResponse>>(
      '/v1/search',
      {
        query: criteria.query,
        filters: criteria.filters,
        sort_by: criteria.sortBy,
        sort_order: criteria.sortOrder,
        page: criteria.page || 0,
        limit: criteria.limit || 20
      }
    );

    logger.debug('🔄 Search response received', {
      totalResults: response.data.data.total,
      resultCount: response.data.data.games.length
    });

    return response.data.data;
  } catch (error) {
    logger.error('🚨 Game search failed', { error });
    throw error;
  }
};

/**
 * Searches for games by platform
 *
 * @async
 * @function searchGamesByPlatform
 * @param {string} platform - The platform to search for
 * @returns {Promise<SearchResponse>} A promise that resolves to the raw API response
 */
export const searchGamesByPlatform = async (
  platform: string
): Promise<SearchResponse> => {
  return searchGames({
    query: '',
    filters: { platforms: [platform] }
  });
};

/**
 * Searches for games by genre
 *
 * @async
 * @function searchGamesByGenre
 * @param {string} genre - The genre to search for
 * @returns {Promise<SearchResponse>} A promise that resolves to the raw API response
 */
export const searchGamesByGenre = async (
  genre: string
): Promise<SearchResponse> => {
  return searchGames({
    query: '',
    filters: { genres: [genre] }
  });
};