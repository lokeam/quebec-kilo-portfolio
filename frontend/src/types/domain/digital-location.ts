/**
 * Digital Location Types
 *
 * Defines types for digital media storage locations.
 */

import type { GamePlatform } from './game-platform';
import type { GameItem } from './game-item';

/**
 * Represents a digital media storage location
 */
export interface DigitalLocation {
  /** Unique identifier for the location */
  id: string;

  /** Name of the digital location */
  name: string;

  /** Type of digital location (e.g., 'steam', 'epic', 'gog') */
  type: GamePlatform;

  /** Optional description of the location */
  description?: string;

  /** Optional metadata for the location */
  metadata?: DigitalLocationMetadata;

  /** List of game items stored in this location */
  items: GameItem[];

  /** Timestamp when the location was created */
  createdAt: string;

  /** Timestamp when the location was last updated */
  updatedAt: string;

  /** Indicates whether the location is a subscription */
  isSubscription: boolean;

  /** Monthly cost of the subscription */
  monthlyCost: number;

  /** Indicates whether the location is active */
  isActive: boolean;

  /** URL of the location */
  url: string;

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
 * Metadata for a digital location
 */
export interface DigitalLocationMetadata {
  /** Platform-specific identifier */
  platformId?: string;

  /** Platform-specific username */
  username?: string;

  /** Platform-specific API key or token */
  apiKey?: string;

  /** Additional platform-specific metadata */
  [key: string]: unknown;
}

/**
 * Request type for creating a new digital location
 */
export interface CreateDigitalLocationRequest {
  /** Name of the digital location */
  name: string;

  /** Type of digital location */
  type: string;

  /** Optional description */
  description?: string;

  /** Optional metadata */
  metadata?: DigitalLocationMetadata;
}