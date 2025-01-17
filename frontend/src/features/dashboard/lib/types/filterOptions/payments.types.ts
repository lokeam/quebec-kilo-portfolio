import type { BILLING_CYCLE_OPTIONS } from '@/features/dashboard/lib/constants/filterOptions/payments/billingCycles.filterOptions';
import type { PAYMENT_METHOD_OPTIONS } from '@/features/dashboard/lib/constants/filterOptions/payments/paymentMethod.filterOptions';

export type BillingCycleOption = typeof BILLING_CYCLE_OPTIONS[number];
export type PaymentMethodOption = typeof PAYMENT_METHOD_OPTIONS[number];

export type PaymentFilterOptions =
  typeof BILLING_CYCLE_OPTIONS |
  typeof PAYMENT_METHOD_OPTIONS;