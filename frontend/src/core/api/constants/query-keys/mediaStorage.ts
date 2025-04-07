/**
 * Parameters for media storage location operations
 */
interface LocationParams {
  /**
   * ID of the location
   * @optional - Only required for specific location queries
   */
  id?: string;

  /**
   * Storage type (physical, digital, sublocation)
   */
  type?: 'physical' | 'digital' | 'sublocation';

  /**
   * Type of location
   * @optional - Used for filtering
   */
  locationType?: string;

  /**
   * Whether to include sublocations in the response
   * @default false
   */
  includeSublocations?: boolean;

  /**
   * Parent location ID for sublocations
   */
  parentId?: string;
}

/**
 * Type representing a media storage query key
 * Used for React Query cache management
 */
type MediaStorageKey = readonly ['mediaStorage', LocationParams?];

/**
 * Query key factory for media storage related operations
 * Follows the project pattern for consistent query key management
 */
export const mediaStorageKeys = Object.freeze({
  /** Base key for all media storage related queries */
  all: ['mediaStorage'] as const,

  /** Key for all physical locations */
  locations: {
    /** All physical locations */
    all: ['mediaStorage', { type: 'physical' }] as const,

    /**
     * Generates a query key for a specific location by ID
     * @param id - The location ID
     */
    byId: (id: string): MediaStorageKey =>
      ['mediaStorage', { type: 'physical', id }],

    /**
     * Generates a query key for locations filtered by type
     * @param locationType - The location type to filter by
     */
    byType: (locationType: string): MediaStorageKey =>
      ['mediaStorage', { type: 'physical', locationType }]
  },

  /** Key for all digital locations */
  digitalLocations: {
    /** All digital locations */
    all: ['mediaStorage', { type: 'digital' }] as const,

    /**
     * Generates a query key for a specific digital location
     * @param id - The digital location ID
     */
    byId: (id: string): MediaStorageKey =>
      ['mediaStorage', { type: 'digital', id }]
  },

  /** Key for all sublocations */
  sublocations: {
    /** All sublocations */
    all: ['mediaStorage', { type: 'sublocation' }] as const,

    /**
     * Generates a query key for sublocations by parent location ID
     * @param locationId - The parent location ID
     */
    byLocation: (locationId: string): MediaStorageKey =>
      ['mediaStorage', { type: 'sublocation', parentId: locationId }]
  }
});
