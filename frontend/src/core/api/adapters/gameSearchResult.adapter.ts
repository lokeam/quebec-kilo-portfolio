import type { Game } from '@/types/domain/game';
import type { SearchResult, SearchMetadata } from '@/types/domain/search';
import type { SearchResponse } from '@/types/api/search';

/**
 * Adapter for transforming game search data between API and domain types
 */
export const gameSearchResultAdapter = {
  /**
   * Transforms API search response to domain search results
   */
  toDomain: (response: SearchResponse): { results: SearchResult[]; metadata: SearchMetadata } => {
    const results: SearchResult[] = response.games.map(game => ({
      game,
      relevance: 1.0, // Default relevance score
      matchType: 'exact', // Default match type
      matchedFields: ['name'], // Default matched fields
    }));

    const metadata: SearchMetadata = {
      totalResults: response.total,
      pageSize: response.games.length,
      currentPage: 0, // Default to first page
      totalPages: Math.ceil(response.total / response.games.length),
      executionTime: 0, // This should come from the API
      timestamp: new Date().toISOString(),
    };

    return { results, metadata };
  },

  /**
   * Creates a library mutation payload from a game
   */
  toLibraryPayload: (game: Game) => ({
    id: Number(game.id),
    name: game.name,
    cover_url: game.coverUrl,
    rating: game.rating ? Number(game.rating) : undefined,
    theme_names: game.themeNames ? [...game.themeNames] : undefined,
  }),

  /**
   * Creates a wishlist mutation payload from a game
   */
  toWishlistPayload: (game: Game) => ({
    id: Number(game.id),
    name: game.name,
    cover_url: game.coverUrl,
    rating: game.rating ? Number(game.rating) : undefined,
    theme_names: game.themeNames ? [...game.themeNames] : undefined,
  }),

  /**
   * Transforms a game to a search result display model
   * Used by the SearchResult component
   */
  toDisplayModel: (game: Game) => ({
    id: game.id,
    name: game.name,
    coverUrl: game.coverUrl,
    isInLibrary: game.isInLibrary ?? false,
    isInWishlist: game.isInWishlist ?? false,
  }),
};