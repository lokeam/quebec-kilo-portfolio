export const PAYMENT_METHOD_OPTIONS = [
  { key: 'Visa', label: 'Visa' },
  { key: 'Mastercard', label: 'Mastercard' },
  { key: 'Amex', label: 'Amex' },
  { key: 'Discover', label: 'Discover' },
  { key: 'Paypal', label: 'Paypal' },
  { key: 'Apple_pay', label: 'Apple Pay' },
  { key: 'Google_pay', label: 'Google Pay' },
  { key: 'Amazon_pay', label: 'Amazon Pay' },
  { key: 'Samsung_pay', label: 'Samsung Pay' },
  { key: 'Jcb', label: 'JCB' },
  { key: 'Mir', label: 'MIR' },
  { key: 'Alipay', label: 'Alipay' },
] as const;

export type FilterOptions = typeof PAYMENT_METHOD_OPTIONS;
