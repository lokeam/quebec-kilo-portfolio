import type { BaseSpendTracking } from '@/features/dashboard/lib/types/spend-tracking/base';
import type { BillingCycle, SpendTransaction } from '@/features/dashboard/lib/types/spend-tracking/constants';
import type { ISO8601Date } from '@/shared/types/types';
import type { OnlineServiceProviderId, OnlineServiceProviderDisplay } from '@/shared/constants/service.constants';
import type { YearlySpending } from '@/features/dashboard/lib/types/spend-tracking/base';

/**
 * Represents a subscription-based spending service
 * Examples include gaming subscriptions like PlayStation Plus or Xbox Game Pass
 */
export interface SubscriptionSpend extends BaseSpendTracking {
  spendTransactionType: SpendTransaction;
  billingCycle: BillingCycle;
  nextBillingDate: ISO8601Date;
  isActive: boolean;
  provider: {
    id: OnlineServiceProviderId;
    displayName: OnlineServiceProviderDisplay;
  }

  /* Subscription-specific fields */
  tier?: string;
  features?: string[];
  autoRenew?: boolean;
  yearlySpending?: YearlySpending[];
}
