/**
 * Digital Location Types
 *
 * All types for digital media storage locations and related operations.
 */

import type { GamePlatform } from './game-platform';
import type { GameItem } from './game-item';

/**
 * Core type representing a digital media storage location
 */
export interface DigitalLocation {
  id: string;
  name: string;
  type: GamePlatform;
  items: GameItem[];
  createdAt: string;
  updatedAt: string;
  isSubscription: boolean;
  monthlyCost: number;
  isActive: boolean;
  url: string;
  paymentMethod: string;
  paymentDate: string;
  billingCycle: string;
  costPerCycle: number;
  nextPaymentDate: string;
}


/**
 * Types used for write operations originating from the Digital Location Page
 *
 * Request type for creating a new digital location
 */
export interface CreateDigitalLocationRequest {
  name: string;
  isActive: boolean;
  url: string;
  isSubscription: boolean;
  subscription?: {
    billing_cycle: string;
    cost_per_cycle: number;
    anchor_date: string;
    payment_method: string;
  };
}


export type BillingCycle = '1 month' | '3 month' | '6 month' | '12 month';


/**
 * Types used for Digital Location BFF Response
 *
 * Represents response object from Digital BFF endpoint
 */
export interface DigitalLocationBFFResponse {
  digitalItems: DigitalLocationBFFResponseItem[];
}

export interface DigitalLocationBFFResponseItem {
  id: string;
  name: string;
  logo?: string;
  locationType?: 'digital';
  itemCount?: number;
  isSubscription: boolean;
  monthlyCost: number;
  isActive: boolean;
  url: string;
  createdAt: string;
  updatedAt: string;
  items: GameItem[];
  paymentMethod: string;
  paymentDate: string;
  billingCycle: string;
  costPerCycle: number;
  nextPaymentDate: string;
}


/**
 * Types used for Catalog Response
 *
 * Represents a digital service item for use in the Combobox components
 */
export interface DigitalServiceItem {
  /** Unique identifier for the service */
  id: string;

  /** Name of the service (e.g. "Xbox") */
  name: string;

  /** Logo identifier for the service */
  logo: string;

  /** Whether this is a subscription service */
  isSubscriptionService: boolean;

  /** URL to the service's website */
  url: string;
}
