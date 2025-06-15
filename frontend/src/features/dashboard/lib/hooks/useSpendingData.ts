import { useMemo } from 'react';

// Types
import type { SpendingItem, YearlySpending } from '@/types/domain/spend-tracking';

interface SpendingData {
  spendingData: YearlySpending[];
  title: string;
  isSubscription: boolean;
}

export function useSpendingData(
  item: SpendingItem,
  oneTimeTotal: YearlySpending[]
): SpendingData {
  const isSubscription = item.spendTransactionType === 'subscription';
  const yearlySpending = isSubscription ? item.yearlySpending : oneTimeTotal;

  return useMemo(() => ({
    spendingData: yearlySpending?.sort((a, b) => b.year - a.year) ?? [],
    title: isSubscription
      ? `Total spent per year on ${item.serviceName?.id ?? 'subscription'}`
      : 'Total spent per year on one-time purchases',
    isSubscription
  }), [yearlySpending, isSubscription, item.serviceName?.id]);
}
