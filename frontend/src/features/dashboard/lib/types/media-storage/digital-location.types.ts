import type { BaseLocation } from '@/features/dashboard/lib/types/media-storage/base';
import type { GamePlatform } from '@/features/dashboard/lib/types/media-storage/constants';

/**
 * Types for the digital locations feature
 * Aligns with the backend models for digital locations
 */

/**
 * Subscription model for digital locations
 */
export interface Subscription {
  id: number;
  locationId: string;
  billingCycle: string;
  costPerCycle: number;
  nextPaymentDate: string;
  paymentMethod: string;
  createdAt: Date;
  updatedAt: Date;
}

/**
 * Digital Location model - represents an online service that stores games
 */
export interface DigitalLocation extends BaseLocation {
  /** Digital platform identifier */
  platform: GamePlatform;

  /** URL to platform storefront */
  url: string;

  /** Service type (basic or subscription) */
  serviceType: string;

  /** Whether the service is active */
  isActive: boolean;

  /** Whether the service is a subscription service */
  isSubscriptionService: boolean;

  /** Service logo identifier */
  logo: string;

  /** Service display label */
  label: string;

  /** Service name */
  name: string;

  /** Billing information */
  billing?: {
    cycle: string;
    fees: {
      monthly: string;
      quarterly: string;
      annual: string;
    };
    paymentMethod: string;
    renewalDate?: {
      day: number;
      month: string;
    };
  };

  /** Platform subscription details */
  subscription?: {
    isActive: boolean;
    isFree: boolean;
    monthlyFee?: string;
    renewalDate?: Date;
    benefits?: string[];
  };

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
 * Interface for creating a new digital location
 */
export interface CreateDigitalLocationRequest {
  name: string;
  platform: GamePlatform;
  url: string;
  subscription?: {
    isActive: boolean;
    isFree: boolean;
    monthlyFee?: string;
    renewalDate?: Date;
    benefits?: string[];
  };
}