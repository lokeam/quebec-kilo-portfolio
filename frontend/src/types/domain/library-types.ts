export interface CreateLibraryGameRequest {
  /** IGDB game ID */
  gameId: number;

  /** Game name */
  gameName: string;

  /** Game cover URL */
  gameCoverUrl: string;

  /** Game release date */
  gameFirstReleaseDate: number;

  /** Game rating */
  gameRating: number;

  /** Game type */
  gameType: {
    displayText: string;
    normalizedText: string;
  };

  /** Game theme names */
  gameThemeNames: string[];

  /** Type of digital location */
  gamesByPlatformAndLocation: CreateLibraryGameRequestLocationEntry[];
}

/**
 * Metadata for a digital location
 */
export interface CreateLibraryGameRequestLocationEntry {
  /** Platform-specific name */
  platformName?: string;

  /** Platform ID according to IGDB */
  platformId: number;

  /** Platform-specific type. May be either digital or physical */
  type: 'digital' | 'physical';

  /** Platform-specific location ID */
  location: {
    sublocationId?: string;
    digitalLocationId?: string;
  };
}

export interface LibraryGameItem {
  id: number;
  name: string;
  coverUrl: string;
  firstReleaseDate: number;
  rating: number;
  isInLibrary: boolean;
  isInWishlist: boolean;
  gameType: {
    displayText: string;
    normalizedText: string;
  };
  favorite: boolean;
  gamesByPlatformAndLocation: Array<{
    gameId: number;
    platformId: number;
    platformName: string;
    type: 'physical' | 'digital';
    locationId: string;
    locationName: string;
    locationType: string;
    sublocationId?: string;
    sublocationName?: string;
    sublocationType?: string;
    sublocationBgColor?: string;
    isActive?: boolean;
  }>;
}

export interface GamePlatformLocation {
  GameID: number;
  PlatformID: number;
  PlatformName: string;
  Type: string;
  LocationID: string;
  LocationName: string;
  LocationType: string;
  SublocationID: string;
  SublocationName: string;
  SublocationType: string;
  SublocationBgColor: string;
  IsActive: boolean | null;
}

export interface GameType {
  displayText: string;
  normalizedText: string;
}

export interface LibraryGame {
  id: number;
  name: string;
  coverUrl: string;
  firstReleaseDate: number;
  rating: number;
  themeNames: string[] | null;
  isInLibrary: boolean;
  isInWishlist: boolean;
  gameType: GameType;
  favorite: boolean;
  gamesByPlatformAndLocation: GamePlatformLocation[];
}