/**
 * Primary location type discriminator.
 * Used to distinguish between physical and digital storage locations.
 *
 * @constant LocationType
 */
export const LocationType = {
  PHYSICAL: 'physical',
  DIGITAL: 'digital'
} as const;

export type LocationType = typeof LocationType[keyof typeof LocationType];

/**
 * Types of physical storage locations.
 * Represents different real-world storage environments.
 *
 * @constant PhysicalLocationType
 */
export const PhysicalLocationType = {
  HOUSE: 'house',
  APARTMENT: 'apartment',
  OFFICE: 'office',
  WAREHOUSE: 'warehouse',
} as const;

export type PhysicalLocationType = typeof PhysicalLocationType[keyof typeof PhysicalLocationType];

/**
 * Types of storage subdivisions within physical locations.
 * Represents specific storage units or furniture.
 *
 * @constant SublocationType
 */
export const SublocationType = {
  SHELF: 'shelf',
  CONSOLE: 'console',
  CABINET: 'cabinet',
  CLOSET: 'closet',
  DRAWER: 'drawer',
  BOX: 'box'
} as const;

export type SublocationType = typeof SublocationType[keyof typeof SublocationType];

/**
 * Digital gaming platforms.
 * Represents major digital distribution services for games.
 *
 * @constant GamePlatform
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
 * Maps platform identifiers to their consumer-facing names.
 *
 * @constant PLATFORM_DISPLAY_NAMES
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
 * Maps platform identifiers to their primary web storefronts.
 *
 * @constant PLATFORM_URLS
 */
export const PLATFORM_URLS: Record<GamePlatform, string> = {
  [GamePlatform.STEAM]: 'https://store.steampowered.com',
  [GamePlatform.EPIC]: 'https://store.epicgames.com',
  [GamePlatform.GOG]: 'https://www.gog.com',
  [GamePlatform.PLAYSTATION]: 'https://store.playstation.com',
  [GamePlatform.XBOX]: 'https://www.xbox.com/games',
  [GamePlatform.NINTENDO]: 'https://www.nintendo.com/store'
};
