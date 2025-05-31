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
  { key: 'jcb', label: 'JCB' },
  { key: 'mir', label: 'MIR' },
  { key: 'alipay', label: 'Alipay' },
] as const;

export type FilterOptions = typeof PAYMENT_METHOD_OPTIONS;
