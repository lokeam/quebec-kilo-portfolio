import type { Game } from './game';

/**
 * Domain representation of a search result
 * Contains the game data and any additional search-specific metadata
 */
export interface SearchResult {
  game: Game;
  relevance: number;      // Search relevance score (0-1)
  matchType: 'exact' | 'partial' | 'fuzzy';  // How the result matched the search
  matchedFields: string[];  // Which fields matched the search (e.g., ['name', 'summary'])
}

/**
 * Domain representation of search parameters
 * Defines how a search should be performed
 */
export interface SearchCriteria {
  query: string;          // The search query string
  filters?: {
    platforms?: string[];  // Filter by platform
    genres?: string[];     // Filter by genre
    themes?: string[];     // Filter by theme
    rating?: number;       // Minimum rating
    releaseYear?: number;  // Release year
  };
  sortBy?: 'relevance' | 'name' | 'rating' | 'releaseDate';
  sortOrder?: 'asc' | 'desc';
  page?: number;          // Page number (0-based)
  limit?: number;         // Number of results per page
}

/**
 * Domain representation of search metadata
 * Contains information about the search results and pagination
 */
export interface SearchMetadata {
  totalResults: number;    // Total number of results
  pageSize: number;        // Number of results per page
  currentPage: number;     // Current page number (0-based)
  totalPages: number;      // Total number of pages
  executionTime: number;   // Search execution time in milliseconds
  timestamp: string;       // When the search was performed (ISO string)
}

/**
 * Storage location types for the add game form
 */
export interface AddGameFormStorageLocationsResponse {
  physicalLocations: AddGameFormPhysicalLocationsResponse[];
  digitalLocations: AddGameFormDigitalLocationsResponse[];
}

export interface AddGameFormPhysicalLocationsResponse {
  parentLocationId: string;
  parentLocationName: string;
  parentLocationType: string;
  parentLocationBgColor: string;
  sublocationId: string;
  sublocationName: string;
  sublocationType: string;
}

export interface AddGameFormDigitalLocationsResponse {
  digitalLocationId: string;
  digitalLocationName: string;
  isSubscription: boolean;
  isActive: boolean;
}