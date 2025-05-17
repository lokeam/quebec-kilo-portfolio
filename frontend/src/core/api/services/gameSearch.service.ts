/**
 * Game Search API Service
 *
 * This service handles API requests related to game search functionality.
 * It provides a unified interface for searching games with various filters and sorting options.
 */

import { axiosInstance } from '@/core/api/client/axios-instance';
import { apiRequest } from '@/core/api/utils/apiRequest';
import type { SearchResponse } from '@/types/api/search';
import type { SearchCriteria } from '@/types/domain/search';

/**
 * Searches for games based on the provided criteria
 *
 * Uses apiRequest helper to wrap the axios call with:
 *  - async/await syntax
 *  - pre‑call debug log
 *  - post‑call success log
 *  - catch block with error log + optional Sentry/metrics
 *  - retry logic (if configured)
 *
 * Usage:
 *   return apiRequest('getUserById', () => axios.get(...));
 */
export const searchGames = (criteria: SearchCriteria): Promise<SearchResponse> =>
  apiRequest('searchGames', () =>
    axiosInstance
      .post<SearchResponse>('/v1/search',
        {
          query: criteria.query,
          filters: criteria.filters,
          sort_order: criteria.sortOrder,
          page: criteria.page ?? 0,
          limit: criteria.limit ?? 20,
        }
      )
      .then(response => response.data)
  );
