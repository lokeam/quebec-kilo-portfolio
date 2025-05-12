import type { Game } from '../domain/game';

/**
 * Response type for search API endpoints
 */
export interface SearchResponse {
  games: Array<{
    id: number;
    name: string;
    coverUrl?: string;
    firstReleaseDate?: number;
    rating?: number;
    themeNames?: string[];
    platformNames?: string[];
    isInLibrary?: boolean;
    isInWishlist?: boolean;
    gameType?: {
      displayText: string;
      normalizedText: string;
    };
  }>;
  total: number;
  page: number;
  pageSize: number;
  totalPages: number;
}

/**
 * Props type for SearchResult component
 */
export interface SearchResultProps {
  game: Game;
  onAction?: () => void;
}

/**
 * Payload type for library/wishlist mutations
 */
export interface SearchMutationPayload {
  id: number;
  name: string;
  cover_url: string;
  rating?: number;
  theme_names?: string[];
}