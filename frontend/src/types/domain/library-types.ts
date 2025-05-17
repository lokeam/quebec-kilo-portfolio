export interface CreateLibraryGameRequest {
  /** Name of the digital location */
  name: string;

  /** Type of digital location */
  type: string;

  /** Optional description */
  description?: string;

  /** Optional metadata */
  metadata?: CreateLibraryGameMetadata;
}

/**
 * Metadata for a digital location
 */
export interface CreateLibraryGameMetadata {
  /** Platform-specific identifier */
  platformId?: string;

  /** Platform-specific username */
  username?: string;

  /** Platform-specific API key or token */
  apiKey?: string;

  /** Additional platform-specific metadata */
  [key: string]: unknown;
}