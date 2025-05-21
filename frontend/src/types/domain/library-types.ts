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
  physicalLocations: Array<{
    name: string;
    type: string;
    sublocation: {
      name: string;
      type: string;
      bgColor: string;
    };
    platform: string;
  }>;
  digitalLocations: any[]; // We can type this properly when we need it
}