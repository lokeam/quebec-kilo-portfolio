/**
 * Digital Location Types
 *
 * Type definitions for digital locations and related entities.
 */

export interface DigitalLocation {
  id: string;
  userId: string;
  name: string;
  serviceType: 'basic' | 'subscription';
  isActive: boolean;
  url?: string;
  createdAt: string;
  updatedAt: string;
}

export interface CreateDigitalLocationInput {
  name: string;
  serviceType: 'basic' | 'subscription';
  isActive: boolean;
  url?: string;
}

export interface UpdateDigitalLocationInput {
  name?: string;
  serviceType?: 'basic' | 'subscription';
  isActive?: boolean;
  url?: string;
}

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