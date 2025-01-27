import type { BaseLocation } from '@/features/dashboard/lib/types/media-storage/base';
import type { GamePlatform } from '@/features/dashboard/lib/types/media-storage/constants';

/**
 * Configuration for digital platform subscriptions.
 * Tracks subscription status, costs, + associated benefits.
 *
 * @interface PlatformSubscription
 */
export interface PlatformSubscription {
  /** Whether or not subscription is currently active */
  isActive: boolean;

  /** Whether or not platform is free to use */
  isFree: boolean;

  /** Monthly subscription cost (if applicable) */
  monthlyFee?: string;

  /** Next renewal date for the subscription */
  renewalDate?: Date;

  /** List of subscription benefits + features */
  benefits?: string[];
}

/**
 * Digital storage location configuration.
 * Extends BaseLocation with digital platform-specific properties.
 *
 * @interface DigitalLocation
 * @extends {BaseLocation}
 */
export interface DigitalLocation extends BaseLocation {
  /** Digital platform identifier */
  platform: GamePlatform;

  /** URL to platform storefront */
  url: string;

  /** Platform subscription details */
  subscription?: PlatformSubscription;

  /** Last synchronization timestamp with the platform */
  lastSync?: Date;

  /** Storage capacity information */
  totalStorage?: {
    /** Amount of storage used */
    used: number;
    /** Total available storage */
    total: number;
    /** Storage unit (GB or TB) */
    unit: 'GB' | 'TB';
  };
}

/**
 * Input type for creating new digital locations.
 * Omits server-generated fields from DigitalLocation.
 *
 * @type CreateDigitalLocationInput
 */
export type CreateDigitalLocationInput = Omit<
  DigitalLocation,
  'id' | 'createdAt' | 'updatedAt'
>;

/**
 * Input type for updating existing digital locations.
 * Makes all fields optional except id.
 *
 * @type UpdateDigitalLocationInput
 */
export type UpdateDigitalLocationInput = Partial<
  Omit<DigitalLocation, 'id' | 'createdAt' | 'updatedAt'>
>;
