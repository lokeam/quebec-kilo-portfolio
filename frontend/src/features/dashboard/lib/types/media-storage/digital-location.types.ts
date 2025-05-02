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
  location_id: string;
  billing_cycle: string;
  cost_per_cycle: number;
  next_payment_date: string;
  payment_method: string;
  created_at: string;
  updated_at: string;
}

/**
 * Digital Location model - represents an online service that stores games
 */
export interface DigitalLocation extends BaseLocation {
  /** Digital platform identifier */
  platform: GamePlatform;

  /** URL to platform storefront */
  url: string;

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