// src/types/game.ts

interface GameType {
  display_text: string;
  normalized_text: string;
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
  platforms?: number[];
  genres?: number[];
  themes?: number[];
  gameTypeId?: number;
  gameType?: GameType;
  isInLibrary?: boolean;
  isInWishlist?: boolean;
  platformNames?: string[];
  genreNames?: string[];
  themeNames?: string[];
  platform?: string;
  platformVersion?: string;
  acquiredDate?: string;
  condition?: string;
  hasOriginalCase?: boolean;
  hasManual?: boolean;
}