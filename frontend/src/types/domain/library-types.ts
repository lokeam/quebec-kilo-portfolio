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

// THIS TYPE TO BE SINGLE POINT OF TRUTH FOR ALL LIBRARY GAME DATA
export interface LibraryGameItemResponse {
  id: number;
  name: string;
  coverUrl: string;
  firstReleaseDate: number;
  rating: number;
  isInLibrary: boolean;
  isInWishlist: boolean;
  isUniqueCopy: boolean;
  gameType: {
    displayText: string;
    normalizedText: string;
  };
  favorite: boolean;
  themeNames?: string[];
  genreNames?: string[];
  gamesByPlatformAndLocation: Array<{
    gameId: number;
    platformId: number;
    platformName: string;
    platformCategory?: string;
    type: 'physical' | 'digital';
    locationId: string;
    sublocationId?: string;
    sublocationName?: string;
    sublocationType?: string;
    parentLocationId?: string;
    parentLocationType?: string;
    parentLocationBgColor?: string;
    parentLocationName?: string;
  }>;
}

export interface GamePlatformLocationResponse {
  gameId: number;
  platformId: number;
  platformName: string;
  platformCategory?: string;
  type: 'digital' | 'physical';
  locationId: string;
  sublocationId?: string;
  sublocationName?: string;
  sublocationType?: string;
  parentLocationId?: string;
  parentLocationType?: string;
  parentLocationBgColor?: string;
  parentLocationName?: string;
}

// Library BFF Response
export interface LibraryItemsBFFResponse {
  libraryItems: LibraryGameItemResponse[];
  recentlyAdded: LibraryGameItemResponse[];
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
  themeNames: string[] | null;
  gamesByPlatformAndLocation: Array<{
    gameId: number;
    platformId: number;
    platformName: string;
    type: 'physical' | 'digital';
    locationId: string;
    locationName: string;   // LEGACY FIELD: DO NOT USE. SEE LocationsBFFSublocationResponse type
    locationType?: string;  // LEGACY FIELD: DO NOT USE. SEE LocationsBFFSublocationResponse type
    sublocationId?: string;
    sublocationName?: string;
    sublocationType?: string;
    sublocationBgColor?: string; // LEGACY FIELD: DO NOT USE. SEE LocationsBFFSublocationResponse type
    parentLocationType?: string;
    parentLocationBgColor?: string;
    parentLocationName?: string;
    isActive?: boolean; // LEGACY FIELD: DO NOT USE. MARKED FOR DELETION.
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


// -- REFACTORED LIBRARY RESPONSE TYPES, TO LEGACY TYPES ABOVE WHEN COMPLETE --
export interface LibraryGameItemRefactoredResponse {
  id: number;
  name: string;
  coverUrl: string;
  isInWishlist: boolean;
  firstReleaseDate: number;
  genreNames: string[];
  gameType: {
    displayText: string;
    normalizedText: string;
  };
  favorite: boolean;
  totalPhysicalVersions: number;
  totalDigitalVersions: number;
  physicalLocations: PhysicalLocationResponse[];
  digitalLocations: DigitalLocationResponse[];
}

// Physical location response
export interface PhysicalLocationResponse {
  parentLocationName: string;
  parentLocationId: string;
  parentLocationType: string;
  parentLocationBgColor: string;
  sublocationName: string;
  sublocationId: string;
  sublocationType: string;
  gamePlatformVersions: PlatformVersionResponse[];
}

// Digital location response
export interface DigitalLocationResponse {
  digitalLocationName: string;
  digitalLocationId: string;
  gamePlatformVersions: PlatformVersionResponse[];
}

// Platform version response
export interface PlatformVersionResponse {
  platformName: string;
  platformId: number;
}

// Updated BFF response
export interface LibraryItemsBFFRefactoredResponse {
  libraryItems: LibraryGameItemRefactoredResponse[];
  recentlyAdded: LibraryGameItemRefactoredResponse[];
}