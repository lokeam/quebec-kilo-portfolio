export interface CreateLibraryGameRequest {
  /** Name of the game */
  gameId: number;

  /** Type of digital location */
  gamesByPlatformAndLocation: CreateLibraryGameRequestLocationEntry[];
}

/**
 * Metadata for a digital location
 */
export interface CreateLibraryGameRequestLocationEntry {
  /** Platform-specific name */
  platformName?: string;

  /** Platform-specific type. May be either digital or physical */
  type: 'digital' | 'physical';

  /** Platform-specific location ID */
  location: {
    sublocationId?: string;
    digitalLocationId?: string;
  };
}