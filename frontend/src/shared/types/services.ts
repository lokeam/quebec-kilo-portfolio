
/**
 * All possible service provider keys
 */
export type ServiceProviderKey =
  | 'APPLE'
  | 'EA'
  | 'EPIC'
  | 'FANATICAL'
  | 'GOG'
  | 'GOOGLE'
  | 'GREENMAN'
  | 'HUMBLE'
  | 'META'
  | 'MICROSOFT'
  | 'NETFLIX'
  | 'NINTENDO'
  | 'NVIDIA'
  | 'PRIME'
  | 'SONY'
  | 'STEAM';

/**
 * All possible service provider IDs in the system
 * This type is derived from our constants to ensure they stay in sync
 */
export type ServiceProviderId =
  | 'apple'
  | 'ea'
  | 'epic'
  | 'fanatical'
  | 'gog'
  | 'google'
  | 'greenman'
  | 'humble'
  | 'meta'
  | 'microsoft'
  | 'netflix'
  | 'nintendo'
  | 'nvidia'
  | 'prime'
  | 'sony'
  | 'steam';

/**
 * Represents the structure of a service provider's metadata
 */
export interface ServiceProvider {
  readonly displayName: string;
  readonly id: string;
}

/**
 * A record of all service providers and their metadata
 * Used as a type for our constants to ensure type safety
 */
export type ServiceProviderRecord = {
  readonly [K in ServiceProviderKey]: ServiceProvider;
};
