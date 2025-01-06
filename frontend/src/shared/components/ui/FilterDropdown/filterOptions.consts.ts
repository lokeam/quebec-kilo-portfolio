export const BILLING_CYCLE_OPTIONS = [
  { key: 'free', label: 'Free' },
  { key: 'monthly', label: 'Monthly' },
  { key: 'quarterly', label: 'Quarterly' },
  { key: 'yearly', label: 'Yearly' },
] as const;

export const PAYMENT_METHOD_OPTIONS = [
  { key: 'visa', label: 'Visa' },
  { key: 'mastercard', label: 'Mastercard' },
  { key: 'amex', label: 'Amex' },
  { key: 'discover', label: 'Discover' },
  { key: 'paypal', label: 'Paypal' },
  { key: 'apple_pay', label: 'Apple Pay' },
  { key: 'google_pay', label: 'Google Pay' },
  { key: 'amazon_pay', label: 'Amazon Pay' },
  { key: 'samsung_pay', label: 'Samsung Pay' },
] as const;

export type FilterOptions = typeof BILLING_CYCLE_OPTIONS | typeof PAYMENT_METHOD_OPTIONS;
