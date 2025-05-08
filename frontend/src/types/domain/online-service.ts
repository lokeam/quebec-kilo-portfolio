/**
 * Online Service Types
 *
 * Defines types for online gaming services and platforms.
 */

import type { GameItem } from './game-item';
import type { ServiceStatusCode, ServiceType } from '@/shared/constants/service.constants';

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
 * Represents an online gaming service or platform
 */
export interface OnlineService {
  /** Unique identifier for the service */
  id: string;

  /** Name of the service */
  name: string;

  /** Display label for the service */
  label: string;

  /** URL to the service's logo */
  logo: string;

  /** Service's website URL */
  url: string;

  /** Current status of the service */
  status: ServiceStatusCode;

  /** Timestamp when the service was added */
  createdAt: Date;

  /** Timestamp when the service was last updated */
  updatedAt: Date;

  /** Whether this is a subscription-based service */
  isSubscriptionService: boolean;

  /** Type of service */
  serviceType: ServiceType;

  /** Whether the service is currently active */
  isActive: boolean;

  /** Category of the service */
  type: ServiceType;

  /** Billing configuration and payment details */
  billing?: OnlineServiceBilling;

  /** Service tier information and available upgrades */
  tier: TierConfiguration;

  /** List of features provided by the service */
  features: string[];

  /** List of games available through this service */
  games?: GameItem[];
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