/**
 * Supported payment methods in our system
 * Based on major payment providers' standards
 */
// export const PAYMENT_METHODS = {
//   ALIPAY: 'Alipay',
//   AMEX: 'Amex',
//   CODE: 'Code',
//   CODE_FRONT: 'CodeFront',
//   DINERS: 'Diners',
//   DISCOVER: 'Discover',
//   ELO: 'Elo',
//   GENERIC: 'Generic',
//   HIPER: 'Hiper',
//   HIPERCARD: 'Hipercard',
//   JCB: 'Jcb',
//   MAESTRO: 'Maestro',
//   MASTERCARD: 'Mastercard',
//   MIR: 'Mir',
//   PAYPAL: 'Paypal',
//   UNIONPAY: 'Unionpay',
//   VISA: 'Visa',
// } as const;

// export type PaymentMethod = typeof PAYMENT_METHODS[keyof typeof PAYMENT_METHODS];

// --- Payment method refactor
export type PaymentMethodKey =
  | 'ALIPAY'
  | 'AMEX'
  | 'DINERS'
  | 'DISCOVER'
  | 'ELO'
  | 'GENERIC'
  | 'HIPER'
  | 'HIPERCARD'
  | 'JCB'
  | 'MAESTRO'
  | 'MASTERCARD'
  | 'MIR'
  | 'PAYPAL'
  | 'UNIONPAY'
  | 'VISA';

  export type PaymentMethodId =
  | 'alipay'
  | 'amex'
  | 'diners'
  | 'discover'
  | 'elo'
  | 'generic'
  | 'hiper'
  | 'hipercard'
  | 'jcb'
  | 'maestro'
  | 'mastercard'
  | 'mir'
  | 'paypal'
  | 'unionpay'
  | 'visa';

export interface PaymentMethod {
  readonly id: string;
  readonly displayName: string;
  readonly type: 'credit' | 'debit' | 'digital' | 'other';
  [key: string]: string | number | boolean | undefined;
}

export type PaymentMethodRecord = {
  readonly [K in PaymentMethodKey]: PaymentMethod;
}

export const PAYMENT_METHODS: PaymentMethodRecord = {
  ALIPAY: {
    displayName: 'Alipay',
    id: 'alipay',
    type: 'digital'
  },
  AMEX: {
    displayName: 'Amex',
    id: 'amex',
    type: 'credit'
  },
  DINERS: {
    displayName: 'Diners',
    id: 'diners',
    type: 'credit'
  },
  DISCOVER: {
    displayName: 'Discover',
    id: 'discover',
    type: 'credit'
  },
  ELO: {
    displayName: 'Elo',
    id: 'elo',
    type: 'debit'
  },
  GENERIC: {
    displayName: 'Generic',
    id: 'generic',
    type: 'other'
  },
  HIPER: {
    displayName: 'Hiper',
    id: 'hiper',
    type: 'debit'
  },
  HIPERCARD: {
    displayName: 'Hipercard',
    id: 'hipercard',
    type: 'debit'
  },
  JCB: {
    displayName: 'Jcb',
    id: 'jcb',
    type: 'credit'
  },
  MAESTRO: {
    displayName: 'Maestro',
    id: 'maestro',
    type: 'debit'
  },
  MASTERCARD: {
    displayName: 'Mastercard',
    id: 'mastercard',
    type: 'credit'
  },
  MIR: {
    displayName: 'Mir',
    id: 'mir',
    type: 'credit'
  },
  PAYPAL: {
    displayName: 'Paypal',
    id: 'paypal',
    type: 'digital'
  },
  UNIONPAY: {
    displayName: 'Unionpay',
    id: 'unionpay',
    type: 'digital'
  },
  VISA: {
    displayName: 'Visa',
    id: 'visa',
    type: 'credit'
  }
} as const;

export type PaymentMethodDisplay = typeof PAYMENT_METHODS[keyof typeof PAYMENT_METHODS]['displayName'];
// --- Payment method refactor

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

// Define a proper type for payment methods to use in return type
type ValidPaymentMethod = "alipay" | "amex" | "diners" | "discover" | "elo" | "generic" |
  "hiper" | "hipercard" | "jcb" | "maestro" | "mastercard" | "mir" | "paypal" | "unionpay" | "visa";

// Fixed function with proper return type
export function validatePaymentMethod(method: string | undefined): ValidPaymentMethod {
  const validMethods = ["alipay", "amex", "diners", "discover", "elo", "generic",
    "hiper", "hipercard", "jcb", "maestro", "mastercard", "mir", "paypal", "unionpay", "visa"];

  return validMethods.includes(method?.toLowerCase() || '')
    ? (method?.toLowerCase() as ValidPaymentMethod)
    : 'generic';
}
