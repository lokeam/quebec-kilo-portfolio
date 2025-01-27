/**
 * Supported payment methods in our system
 * Based on major payment providers' standards
 */
export const PAYMENT_METHODS = {
  ALIPAY: 'Alipay',
  AMEX: 'Amex',
  CODE: 'Code',
  CODE_FRONT: 'CodeFront',
  DINERS: 'Diners',
  DISCOVER: 'Discover',
  ELO: 'Elo',
  GENERIC: 'Generic',
  HIPER: 'Hiper',
  HIPERCARD: 'Hipercard',
  JCB: 'Jcb',
  MAESTRO: 'Maestro',
  MASTERCARD: 'Mastercard',
  MIR: 'Mir',
  PAYPAL: 'Paypal',
  UNIONPAY: 'Unionpay',
  VISA: 'Visa',
} as const;

export type PaymentMethod = typeof PAYMENT_METHODS[keyof typeof PAYMENT_METHODS];

/**
 * Defines all possible billing cycles for subscription services
 */
export const BILLING_CYCLES = {
  MONTHLY: '1 month',
  QUARTERLY: '3 month',
  BIANNUAL: '6 month',
  ANNUAL: '1 year',
} as const;

export type BillingCycle = typeof BILLING_CYCLES[keyof typeof BILLING_CYCLES];
