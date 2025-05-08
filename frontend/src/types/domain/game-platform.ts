/**
 * Digital gaming platforms.
 * Represents major digital distribution services for games.
 */
export const GamePlatform = {
  STEAM: 'steam',
  EPIC: 'epic',
  GOG: 'gog',
  PLAYSTATION: 'sony',
  XBOX: 'xbox',
  NINTENDO: 'nintendo'
} as const;

export type GamePlatform = typeof GamePlatform[keyof typeof GamePlatform];

/**
 * Human-readable display names for gaming platforms.
 */
export const PLATFORM_DISPLAY_NAMES: Record<GamePlatform, string> = {
  [GamePlatform.STEAM]: 'Steam',
  [GamePlatform.EPIC]: 'Epic Games',
  [GamePlatform.GOG]: 'GOG',
  [GamePlatform.PLAYSTATION]: 'PlayStation Network',
  [GamePlatform.XBOX]: 'Xbox Live',
  [GamePlatform.NINTENDO]: 'Nintendo eShop'
};

/**
 * Store URLs for gaming platforms.
 */
export const PLATFORM_URLS: Record<GamePlatform, string> = {
  [GamePlatform.STEAM]: 'https://store.steampowered.com',
  [GamePlatform.EPIC]: 'https://store.epicgames.com',
  [GamePlatform.GOG]: 'https://www.gog.com',
  [GamePlatform.PLAYSTATION]: 'https://store.playstation.com',
  [GamePlatform.XBOX]: 'https://www.xbox.com/games',
  [GamePlatform.NINTENDO]: 'https://www.nintendo.com/store'
};