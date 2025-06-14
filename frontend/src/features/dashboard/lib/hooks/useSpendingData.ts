import { useMemo } from 'react';

// Local Type Definitions
interface BaseSpendItem {
  id: string;
  title: string;
  amount: number;
  spendTransactionType: 'subscription' | 'one-time';
  paymentMethod: string;
  mediaType: string;
  serviceName?: {
    id: string;
    displayName: string;
  };
  createdAt: number;
  updatedAt: number;
  isActive: boolean;
}

interface SubscriptionSpend extends BaseSpendItem {
  spendTransactionType: 'subscription';
  billingCycle: string;
  nextBillingDate: number;
  yearlySpending: Array<{
    year: number;
    amount: number;
  }>;
}

interface OneTimeSpend extends BaseSpendItem {
  spendTransactionType: 'one-time';
  isDigital: boolean;
  isWishlisted: boolean;
  purchaseDate: number;
}

interface YearlySpending {
  year: number;
  amount: number;
}

interface SpendingData {
  spendingData: YearlySpending[];
  title: string;
  isSubscription: boolean;
}

// Type Guard
const isSubscriptionSpend = (item: SubscriptionSpend | OneTimeSpend): item is SubscriptionSpend => {
  return item.spendTransactionType === 'subscription';
};

export function useSpendingData(
  item: SubscriptionSpend | OneTimeSpend,
  oneTimeTotal: YearlySpending[]
): SpendingData {
  const isSubscription = isSubscriptionSpend(item);
  const yearlySpending = isSubscription ? item.yearlySpending : oneTimeTotal;

  return useMemo(() => ({
    spendingData: yearlySpending?.sort((a, b) => b.year - a.year) ?? [],
    title: isSubscription
      ? `Total spent per year on ${item.serviceName?.id ?? 'subscription'}`
      : 'Total spent per year on one-time purchases',
    isSubscription
  }), [yearlySpending, isSubscription, item.serviceName?.id]);
}
