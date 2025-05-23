/**
 * Game Library Service
 *
 * Provides functions for managing game library through the backend API.
 */

import { axiosInstance } from '@/core/api/client/axios-instance';
import { apiRequest } from '@/core/api/utils/apiRequest';

// types
import type { ApiResponse } from '@/types/api/api.types';
import type { Game } from '@/types/game';
import type { CreateLibraryGameRequest, LibraryGameItem } from '@/types/domain/library-types';


const LIBRARY_ENDPOINT = '/v1/library';

interface LibraryResponseWrapper {
  library: {
    success: boolean;
    games?: LibraryGameItem[];
    game?: LibraryGameItem;
  };
}

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
export const getAllLibraryGames = (): Promise<LibraryGameItem[]> =>
  apiRequest('getAllLibraryGames', () =>
    axiosInstance
      .get<LibraryResponseWrapper>(LIBRARY_ENDPOINT)
      .then(response => response.data.library.games || [])
  );

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
export const getLibraryGameById = (id: string): Promise<Game> =>
  apiRequest(`getGameById(${id})`, () =>
    axiosInstance
      .get<LibraryResponseWrapper>(`${LIBRARY_ENDPOINT}/${id}`)
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
export const createLibraryGame = (data: CreateLibraryGameRequest): Promise<Game> =>
  apiRequest('createGame', () =>
    axiosInstance
      .post<LibraryResponseWrapper>(LIBRARY_ENDPOINT, data)
      .then(response => {
        const game = response.data.library.game;
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
export const updateLibraryGame = (id: string, data: Partial<CreateLibraryGameRequest>): Promise<Game> =>
  apiRequest(`updateGame(${id})`, () =>
    axiosInstance
      .put<LibraryResponseWrapper>(`${LIBRARY_ENDPOINT}/${id}`, data)
      .then(response => {
        const game = response.data.library.game;
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
export const deleteLibraryGame = (id: string): Promise<void> =>
  apiRequest(`deleteGame(${id})`, () =>
    axiosInstance
      .delete<ApiResponse<void>>(`${LIBRARY_ENDPOINT}/${id}`)
      .then(response => response.data.data)
);
