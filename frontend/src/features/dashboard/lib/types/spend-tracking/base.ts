import type { Currency, ISO8601Date } from '@/shared/types/types';
import type { PaymentMethod, SpendTransaction } from '@/features/dashboard/lib/types/spend-tracking/constants';
import type { PurchasedMediaCategory } from '@/features/dashboard/lib/types/spend-tracking/media';

/**
 * Yearly spending record interface
 */
export interface YearlySpending {
  year: number;
  amount: Currency;
}

/**
 * Base interface for all spending transactions
 * Contains common fields required for any type of purchase
 */
export interface BaseSpendTracking {
  id: string;
  amount: Currency;
  title: string;
  spendTransactionType: SpendTransaction;
  paymentMethod: PaymentMethod;
  mediaType: PurchasedMediaCategory;
  createdAt: ISO8601Date;
  updatedAt: ISO8601Date;
};
