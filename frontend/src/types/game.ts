// src/types/game.ts
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