/**
 * Game Search API Service
 *
 * This service handles API requests related to game search functionality.
 * It provides a unified interface for searching games with various filters and sorting options.
 */

import { axiosInstance } from '@/core/api/client/axios-instance';
import { apiRequest } from '@/core/api/utils/apiRequest';

// Types
import type { SearchResponse } from '@/types/api/search';
import type { SearchCriteria, AddGameFormStorageLocationsResponse, AddGameFormDigitalLocationsResponse, AddGameFormPhysicalLocationsResponse } from '@/types/domain/search';

// Constants
import { API_BASE_PATH } from '@/core/api/config';

const SEARCH_ENDPOINT = `${API_BASE_PATH}/search`;
const SEARCH_BFF_ENDPOINT = `${API_BASE_PATH}/search/bff`;

// Response wrapper for BFF endpoint
interface AddGameFormBFFResponseWrapper {
  storageLocations: {
    success: boolean;
    physicalLocations: AddGameFormPhysicalLocationsResponse[];
    digitalLocations: AddGameFormDigitalLocationsResponse[];
  };

  metadata: {
    timestamp: string;
    requestId: string;
  };
}

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
      .post<SearchResponse>(SEARCH_ENDPOINT,
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

/**
 * Fetches all storage locations (physical and digital) for the add game form
 *
 * @returns Promise<AddGameFormStorageLocationsResponse>
 * @throws Error if the request fails or response is invalid
 */
export const getAddGameFormStorageLocationsBFF = (): Promise<AddGameFormStorageLocationsResponse> =>
  apiRequest('getAddGameFormStorageLocationsBFF', () =>
    axiosInstance
      .get<AddGameFormBFFResponseWrapper>(SEARCH_BFF_ENDPOINT)
      .then(response => {
        if (!response.data.storageLocations) {
          throw new Error('Invalid response structure from server');
        }
        return response.data.storageLocations;
      })
  );
