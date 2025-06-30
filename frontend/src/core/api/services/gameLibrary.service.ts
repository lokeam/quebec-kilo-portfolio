/**
 * Game Library Service
 *
 * Provides functions for managing game library through the backend API.
 */

import { axiosInstance } from '@/core/api/client/axios-instance';
import { apiRequest } from '@/core/api/utils/apiRequest';

// types
// import type { ApiResponse } from '@/types/api/api.types';
import type { CreateLibraryGameRequest, LibraryGameItemRefactoredResponse, LibraryGameItemRefactoredResponse, LibraryGameItemResponse } from '@/types/domain/library-types';


const LIBRARY_ENDPOINT = '/v1/library';
const LIBRARY_BFF_ENDPOINT = '/v1/library/bff';

interface LibraryOperationResponseWrapper {
  success: boolean;
  library: {
    id: number;
    message: string;
  };
  metadata: {
    timestamp: string;
    request_id: string;
  };
}

interface LibraryGameResponseWrapper {
  success: boolean;
  library: {
    game: LibraryGameItemResponse;
  };
  metadata: {
    timestamp: string;
    request_id: string;
  };
}

type LibraryOperationResponse = {
  id: number;
  message: string;
};

interface LibraryBFFResponseWrapper {
  success: boolean;
  library: {
    libraryItems: LibraryGameItemResponse[];
    recentlyAdded: LibraryGameItemResponse[];
  };
  metadata: {
    timestamp: string;
    request_id: string;
  };
}

type LibraryBFFResponse = {
  libraryItems: LibraryGameItemResponse[];
  recentlyAdded: LibraryGameItemResponse[];
};

// REFACTORED RESPONSE - REMOVE UNUSED types WHEN COMPLETE
interface LibraryBFFRefactoredResponseWrapper {
  success: boolean;
  library: {
    libraryItems: LibraryGameItemRefactoredResponse[];
    recentlyAdded: LibraryGameItemRefactoredResponse[];
  };
  metadata: {
    timestamp: string;
    request_id: string;
  };
}

type LibraryBFFRefactoredResponse = {
  libraryItems: LibraryGameItemRefactoredResponse[];
  recentlyAdded: LibraryGameItemRefactoredResponse[];
};


/**
 * Fetches all games for the current user.
 *
 * Uses apiRequest helper to wrap the axios call with:
 *  - async/await syntax
 *  - pre‑call debug log
 *  - post‑call success log
 *  - catch block with error log + optional Sentry/metrics
 *  - retry logic (if configured)
 *
 * Usage:
 *   return apiRequest('getAllGames', () => axios.get(...));
 */

export const getLibraryPageBFFResponse = (): Promise<LibraryBFFRefactoredResponse> =>
  apiRequest('getLibraryPageBFFResponse', async () => {
    console.log('[DEBUG] getLibraryPageBFFResponse: Making API request');
    const response = await axiosInstance.get<LibraryBFFRefactoredResponseWrapper>(LIBRARY_BFF_ENDPOINT);
    console.log('[DEBUG] getLibraryPageBFFResponse: Raw API response:', response.data);

    if (!response.data.library) {
      console.error('[DEBUG] getLibraryPageBFFResponse: No library data in response:', response.data);
    }

    console.log('[DEBUG] getLibraryPageBFFResponse: Successfully extracted library data:', response.data.library);
    return response.data.library;
  });

// Legacy: DO NOT USE. This was the old way of fetching all library games.
// It was replaced with the BFF endpoint above.
// export const getAllLibraryGames = (): Promise<LibraryGameItemResponse[]> =>
//   apiRequest('getAllLibraryGames', () =>
//     axiosInstance
//       .get<LibraryResponseWrapper>(LIBRARY_ENDPOINT)
//       .then(response => response.data.library.games || [])
//   );

/**
 * Fetches a specific game by ID
 *
 * Uses apiRequest helper to wrap the axios call with:
 *  - async/await syntax
 *  - pre‑call debug log
 *  - post‑call success log
 *  - catch block with error log + optional Sentry/metrics
 *  - retry logic (if configured)
 *
 * Usage:
 *   return apiRequest('getGameById', () => axios.get(...));
 */
export const getLibraryGameById = (id: string): Promise<LibraryGameItemResponse> =>
  apiRequest(`getGameById(${id})`, () =>
    axiosInstance
      .get<LibraryGameResponseWrapper>(`${LIBRARY_ENDPOINT}/${id}`)
      .then(response => {
        const game = response.data.library.game;
        if (!game) {
          throw new Error(`Game with id ${id} not found`);
        }
        return game;
      })
  );


/**
 * Creates a new game in the library
 *
 * Uses apiRequest helper to wrap the axios call with:
 *  - async/await syntax
 *  - pre‑call debug log
 *  - post‑call success log
 *  - catch block with error log + optional Sentry/metrics
 *  - retry logic (if configured)
 *
 * Usage:
 *   return apiRequest('createGame', () => axios.post(...));
 */
export const createLibraryGame = (data: CreateLibraryGameRequest): Promise<LibraryOperationResponse> =>
  apiRequest('createGame', () =>
    axiosInstance
      .post<LibraryOperationResponseWrapper>(LIBRARY_ENDPOINT, data)
      .then(response => {
        const game = response.data.library;
        if (!game) {
          throw new Error('Failed to create game');
        }
        return game;
      })
  );


/**
 * Updates an existing game in the library
 *
 * Uses apiRequest helper to wrap the axios call with:
 *  - async/await syntax
 *  - pre‑call debug log
 *  - post‑call success log
 *  - catch block with error log + optional Sentry/metrics
 *  - retry logic (if configured)
 *
 * Usage:
 *   return apiRequest('updateGame', () => axios.put(...));
 */
export const updateLibraryGame = (id: string, data: Partial<CreateLibraryGameRequest>): Promise<LibraryOperationResponse> =>
  apiRequest(`updateGame(${id})`, () =>
    axiosInstance
      .put<LibraryOperationResponseWrapper>(`${LIBRARY_ENDPOINT}/${id}`, data)
      .then(response => {
        const game = response.data.library;
        if (!game) {
          throw new Error(`Failed to update game with id ${id}`);
        }
        return game;
      })
  );


/**
 * Deletes an existing game in the library
 *
 * Uses apiRequest helper to wrap the axios call with:
 *  - async/await syntax
 *  - pre‑call debug log
 *  - post‑call success log
 *  - catch block with error log + optional Sentry/metrics
 *  - retry logic (if configured)
 *
 * Usage:
 *   return apiRequest('deleteGame', () => axios.put(...));
 */
export const deleteLibraryGame = (id: string): Promise<LibraryOperationResponse> =>
  apiRequest(`deleteGame(${id})`, () =>
    axiosInstance
      .delete<LibraryOperationResponseWrapper>(`${LIBRARY_ENDPOINT}/${id}`)
      .then(response => response.data.library)
);
