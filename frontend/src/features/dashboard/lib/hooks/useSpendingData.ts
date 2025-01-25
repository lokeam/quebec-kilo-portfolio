import { useMemo } from 'react';
import type { SubscriptionSpend } from '@/features/dashboard/lib/types/spend-tracking/subscription';
import type { OneTimeSpend } from '@/features/dashboard/lib/types/spend-tracking/purchases';
import type { YearlySpending } from '@/features/dashboard/lib/types/spend-tracking/base';
import { isSubscriptionSpend } from '@/features/dashboard/lib/types/spend-tracking/guards';

interface SpendingData {
  spendingData: YearlySpending[];
  title: string;
  isSubscription: boolean;
}

export function useSpendingData(
  item: SubscriptionSpend | OneTimeSpend,
  oneTimeTotal: YearlySpending[]
): SpendingData {
  const isSubscription = isSubscriptionSpend(item);
  const yearlySpending = isSubscription ? item.yearlySpending : oneTimeTotal;

  return useMemo(() => ({
    spendingData: yearlySpending?.sort((a, b) => b.year - a.year) ?? [],
    title: isSubscription
      ? `Total spent per year on ${item.provider?.id ?? 'subscription'}`
      : 'Total spent per year on one-time purchases',
    isSubscription
  }), [yearlySpending, isSubscription, item.provider?.id]);
}
