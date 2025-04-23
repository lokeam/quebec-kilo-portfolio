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
export interface DigitalLocation {
  id: string;
  user_id: string;
  name: string;
  service_type: string;
  is_active: boolean;
  url: string;
  created_at: string;
  updated_at: string;
  logo?: string;
  label?: string;
  isSubscriptionService?: boolean;
  billing?: {
    cycle: string;
    fees: {
      monthly: string;
      quarterly: string;
      annual: string;
    };
    paymentMethod: string;
    renewalDate: {
      month: string;
      day: number;
    };
  };
}

/**
 * Interface for creating a new digital location
 */
export interface CreateDigitalLocationRequest {
  name: string;
  service_type: string;
  is_active: boolean;
  url: string;
  subscription?: {
    billing_cycle: string;
    cost_per_cycle: number;
    next_payment_date: string;
    payment_method: string;
  };
}