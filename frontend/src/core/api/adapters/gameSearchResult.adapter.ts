import type { SearchResponse } from '@/types/api/search';
import type { Game } from '@/types/domain/game';
import type { SearchResult, SearchMetadata } from '@/types/domain/search';
import { logger } from '@/core/utils/logger/logger';

/**
 * Adapter for transforming game search data between API and domain types
 */
export const gameSearchResultAdapter = {
  /**
   * Transforms API search response to domain search results
   */
  toDomain: (response: SearchResponse): { results: SearchResult[]; metadata: SearchMetadata } => {
    console.log('Raw API response:', response);

    const results = response.games.map(game => {
      console.log('Transforming game:', game);
      const transformedGame: Game = {
        id: game.id,
        name: game.name,
        coverUrl: game.coverUrl,
        firstReleaseDate: game.firstReleaseDate,
        rating: game.rating,
        themeNames: game.themeNames,
        platforms: game.platforms,
        platformNames: game.platforms.map(p => p.name),
        isInLibrary: game.isInLibrary,
        isInWishlist: game.isInWishlist,
        gameType: game.gameType,
      };
      console.log('Transformed game:', transformedGame);

      return {
        game: transformedGame,
        relevance: 1.0,
        matchType: 'exact' as const,
        matchedFields: ['name'],
      };
    });

    const metadata: SearchMetadata = {
      totalResults: response.total,
      pageSize: response.pageSize,
      currentPage: response.page,
      totalPages: response.totalPages,
      executionTime: 0,
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
  toDisplayModel: (game: Game) => {
    const displayModel = {
      id: game.id,
      name: game.name,
      coverUrl: game.coverUrl,
      firstReleaseDate: game.firstReleaseDate,
      isInLibrary: game.isInLibrary ?? false,
      isInWishlist: game.isInWishlist ?? false,
      gameType: game.gameType,
    };

    logger.debug('🔍 Created display model', {
      original: game,
      displayModel,
      firstReleaseDate: displayModel.firstReleaseDate,
    });

    return displayModel;
  },
};