/**
 * Domain type for Game entities
 * Represents a game in the system with all its properties
 */

export interface GameType {
  displayText: string;
  normalizedText: string;
}

export interface PlatformInfo {
  id: number;
  name: string;
}

export interface Game {
  id: number;
  name: string;
  label?: string;
  summary?: string;
  coverId?: number;
  coverUrl?: string;
  firstReleaseDate?: number;
  rating?: number;
  platforms: PlatformInfo[];  // Array of platform objects with id and name
  genres?: number[];
  themes?: number[];
  gameType?: GameType;
  isInLibrary?: boolean;
  isInWishlist?: boolean;
  platformNames?: string[];  // Keep for backward compatibility
  genreNames?: string[];
  themeNames?: string[];
  platform?: string;
  platformVersion?: string;
  acquiredDate?: string;
  condition?: string;
  hasOriginalCase?: boolean;
  hasManual?: boolean;
}