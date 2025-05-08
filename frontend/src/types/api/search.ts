import type { Game } from '../domain/game';

/**
 * Response type for search API endpoints
 */
export interface SearchResponse {
  games: Game[];
  total: number;
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