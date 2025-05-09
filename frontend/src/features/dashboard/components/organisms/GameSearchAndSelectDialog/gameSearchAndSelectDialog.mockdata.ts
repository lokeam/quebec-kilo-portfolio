import type { Game } from '@/types/domain/game';

export const gameSearchAndSelectDialogMockData = {
  games: [
    {
      id: 1,
      name: 'The Legend of Zelda: Breath of the Wild',
      summary: 'An action-adventure game set in an open world environment.',
      coverUrl: 'https://example.com/zelda.jpg',
      rating: 9.5,
      platformNames: ['Nintendo Switch'],
      genreNames: ['Action', 'Adventure'],
      themeNames: ['Fantasy'],
      isInLibrary: false,
      isInWishlist: false
    },
    {
      id: 2,
      name: 'Super Mario Odyssey',
      summary: 'A 3D platformer featuring Mario and his new companion Cappy.',
      coverUrl: 'https://example.com/mario.jpg',
      rating: 9.3,
      platformNames: ['Nintendo Switch'],
      genreNames: ['Platformer'],
      themeNames: ['Fantasy'],
      isInLibrary: false,
      isInWishlist: false
    }
  ] as Game[],
  total: 2
};
