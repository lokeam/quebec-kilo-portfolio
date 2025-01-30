/**
 * Represents the primary categorization of a spending transaction
 * This affects how the transaction is processed and displayed
 */
export type SpendTransaction = 'subscription' | 'one-time';

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
  NA: 'NA',
  MONTHLY: '1 month',
  QUARTERLY: '3 month',
  BIANNUAL: '6 month',
  ANNUAL: '1 year',
} as const;

export type BillingCycle = typeof BILLING_CYCLES[keyof typeof BILLING_CYCLES];

/**
 * Maps spend tracking concepts to their display names and internal identifiers
 * Used for consistent spend type identification across the spend tracking feature
 */
export const BASE_MEDIA_CATEGORIES = {
  HARDWARE: 'hardware',
  DLC: 'dlc',
  IN_GAME_PURCHASE: 'inGamePurchase',
  SUBSCRIPTION: 'subscription',
  PHYSICAL: 'physical',
  DISC: 'disc'
} as const;

export type BaseMediaCategory = typeof BASE_MEDIA_CATEGORIES[keyof typeof BASE_MEDIA_CATEGORIES];
