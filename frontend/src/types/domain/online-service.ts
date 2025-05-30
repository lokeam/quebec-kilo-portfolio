import type { GameItem } from '@/types/domain/game-item';

/**
 * Online Service Types
 *
 * Defines types for online gaming services and platforms.
 */

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
 * Complete tier configuration for a service.
 * Manages current tier and available upgrade options.
 */
export interface TierConfiguration {
  /** Current tier of the service */
  currentTier: ServiceTierName;

  /** Available tier options */
  availableTiers: ServiceTier[];

  /** Maximum number of devices allowed */
  maxDevices?: number;

  /** Maximum number of users allowed */
  maxUsers?: number;
}

/**
 * Represents an online service from the analytics endpoint
 */
export interface OnlineService {
  /** Unique identifier for the service */
  id: string;

  /** Name of the service (e.g. "Xbox") */
  name: string;

  /** Type of service - always "subscription" for online services */
  location_type: 'subscription';

  /** Number of items associated with this service */
  item_count: number;

  /** Whether this is a subscription service */
  is_subscription: boolean;

  /** Monthly cost of the subscription in the user's currency */
  monthly_cost: number;

  /** Whether the subscription is currently active */
  is_active: boolean;

  /** URL to the service's website */
  url: string;

  /** When the service was created */
  created_at: string;

  /** When the service was last updated */
  updated_at: string;
}

/**
 * Response shape for online services from the analytics endpoint
 */
export interface OnlineServicesResponse {
  success: boolean;
  user_id: string;
  data: {
    analytics: {
      storage: {
        total_digital_locations: number;
        digital_locations: OnlineService[];
      };
    };
  };
  metadata: {
    timestamp: string;
    request_id: string;
  };
}

/**
 * Billing configuration for a service
 */
export interface OnlineServiceBilling {
  /** Billing cycle (monthly, quarterly, annual) */
  cycle: string;

  /** Fee structure for different billing cycles */
  fees: {
    monthly: string;
    quarterly: string;
    annual: string;
  };

  /** Optional renewal date information */
  renewalDate?: {
    month: string;
    day: number;
  };

  /** Payment method information */
  paymentMethod: string | PaymentMethodData;
}

/**
 * Payment method details
 */
export interface PaymentMethodData {
  /** Unique identifier for the payment method */
  id: string;

  /** Display name for the payment method */
  displayName: string;

  /** Additional payment method properties */
  [key: string]: unknown;
}

/**
 * Represents a digital location (online service) from the analytics endpoint
 */
export interface DigitalLocation {
  /** Unique identifier for the digital location */
  id: string;

  /** Name of the service (e.g. "Xbox") */
  name: string;

  /** Type of location - always "subscription" for digital locations */
  locationType?: 'subscription';

  /** Number of items associated with this location */
  itemCount?: number;

  /** Whether this is a subscription service */
  isSubscription: boolean;

  /** Monthly cost of the subscription in the user's currency */
  monthlyCost: number;

  /** Whether the subscription is currently active */
  isActive: boolean;

  /** URL to the service's website */
  url: string;

  /** When the location was created */
  createdAt: string;

  /** When the location was last updated */
  updatedAt: string;

  /** List of items in this location */
  items: GameItem[];

  /** Payment method used for the location */
  paymentMethod: string;

  /** Date of the last payment */
  paymentDate: string;

  /** Billing cycle of the subscription */
  billingCycle: string;

  /** Cost per cycle of the subscription */
  costPerCycle: number;

  /** Next payment date */
  nextPaymentDate: string;
}

/**
 * Response shape from the analytics endpoint
 */
export interface AnalyticsResponse {
  success: boolean;
  user_id: string;
  data: {
    analytics: {
      storage: {
        total_physical_locations: number;
        total_digital_locations: number;
        digital_locations: DigitalLocation[];
        physical_locations: PhysicalLocation[];
      };
    };
  };
  metadata: {
    timestamp: string;
    request_id: string;
  };
}

/**
 * Represents a physical location from the analytics endpoint
 */
export interface PhysicalLocation {
  id: string;
  name: string;
  location_type: 'apartment' | 'house';
  item_count: number;
  map_coordinates: string;
  created_at: string;
  updated_at: string;
  sublocations: Sublocation[];
}

/**
 * Represents a sublocation within a physical location
 */
export interface Sublocation {
  id: string;
  name: string;
  location_type: 'box' | 'console';
  bg_color: string;
  stored_items: number;
  created_at: string;
  updated_at: string;
  items: GameItem[];
}
