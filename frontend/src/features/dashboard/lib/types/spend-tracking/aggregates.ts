import type { SubscriptionSpend } from '@/features/dashboard/lib/types/spend-tracking/subscription';
import type { OneTimeSpend } from '@/features/dashboard/lib/types/spend-tracking/purchases';
import type { Currency } from '@/shared/types/types';
import type { YearlySpending } from '@/features/dashboard/lib/types/spend-tracking/base';

interface SpendingTotals {
  subscriptionTotal: YearlySpending[];
  oneTimeTotal: YearlySpending[];
  combinedTotal: YearlySpending[];
}

/**
 * Represents grouped spending data for reporting and analysis
 * Used throughout the dashboard for various spending views
 */
export interface SpendTrackingData {
  currentTotalThisMonth: (SubscriptionSpend | OneTimeSpend)[];
  recurringNextMonth: (SubscriptionSpend | OneTimeSpend)[];
  oneTimeThisMonth: OneTimeSpend[];
  totalSpendsThisMonth: Currency;
  totalSpendsThisYear: Currency;
  yearlyTotals: SpendingTotals;
}
