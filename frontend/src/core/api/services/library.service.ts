/**
 * Library Service API
 *
 * For API standards and best practices, see:
 * @see {@link ../../../docs/api-standards.md}
 */

import { axiosInstance } from "@/core/api/client/axios-instance";

const LIBRARY_ENDPOINT = '/api/v1/library/add';
const WISHLIST_ENDPOINT = '/api/v1/wishlist/add';

export interface AddGameRequest {
  id: number;
  name: string;
  summary?: string;
  cover_url?: string;
  rating?: number;
  platform_names?: string[];
  genre_names?: string[];
  theme_names?: string[];
}

export interface AddGameResponse {
  success: boolean;
  game: {
    id: number;
    name: string;
  };
}

export const addToLibrary = async(gameData: AddGameRequest): Promise<AddGameResponse> => {
  return await axiosInstance.post(LIBRARY_ENDPOINT, gameData);
};

export const addToWishList = async(gameData: AddGameRequest): Promise<AddGameResponse> => {
  return await axiosInstance.post(WISHLIST_ENDPOINT, gameData);
};

