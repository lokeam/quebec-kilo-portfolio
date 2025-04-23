import type { BaseOnlineService } from '@/features/dashboard/lib/types/online-services/base';
import type { ServiceType } from '@/shared/constants/service.constants';
import type { TierConfiguration } from '@/features/dashboard/lib/types/online-services/tiers';
import type { ViewMode } from '@/features/dashboard/pages/OnlineServices/useViewMode';

/**
 * Complete online service configuration including billing and features.
 * Extends BaseOnlineService with subscription-specific properties.
 */
export interface OnlineService extends BaseOnlineService {
  /** Category of the service */
  type: ServiceType;

  /** Billing configuration and payment details */
  billing?: OnlineServiceBilling;

  /** Service tier information and available upgrades */
  tier: TierConfiguration;

  /** List of features provided by the service */
  features: string[];
}

interface OnlineServiceBilling {
  cycle: string;
  fees: {
    monthly: string;
    quarterly: string;
    annual: string;
  };
  renewalDate?: {
    month: string;
    day: number;
  };
  // Make paymentMethod accept either string or PaymentMethod object
  paymentMethod: string | PaymentMethodData;
}

interface PaymentMethodData {
  id: string;
  displayName: string;
  // other properties...
}

/**
 * API response state for online services data.
 * Used to track loading states and error handling.
 */
export interface OnlineServicesAPIState {
  /** List of available online services */
  services: OnlineService[];

  /** Total count of services in the system */
  totalServices: number;

  /** Error message if API request fails */
  error: string;

  /** Loading state indicator */
  isLoading: boolean;
}

/**
 * UI state for online services view.
 * Manages display preferences and active filters.
 */
export interface OnlineServicesUIState {
  viewMode: ViewMode;
  searchQuery: string;
  billingCycleFilters: string[];
  paymentMethodFilters: string[];
}
