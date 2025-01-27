/**
 * Individual service tier configuration.
 * Represents a specific subscription level within a service.
 */
export interface ServiceTier {
  /** Unique identifier for the tier */
  id: string;

  /** Display name of the tier */
  name: string;

  /** List of features included in this tier */
  features: string[];

  /** Indicates if this is the default tier for new subscriptions */
  isDefault: boolean;
}

/**
 * Available tier names for online services.
 * Standardized naming convention across different providers.
 */
export type ServiceTierName =
  | 'free'
  | 'basic'
  | 'standard'
  | 'premium'
  | 'ultimate';


/**
 * Complete tier configuration for a service.
 * Manages current tier and available upgrade options.
 */
export interface TierConfiguration {
  currentTier: ServiceTierName;
  availableTiers: ServiceTier[];
  maxDevices?: number;
  maxUsers?: number;
}