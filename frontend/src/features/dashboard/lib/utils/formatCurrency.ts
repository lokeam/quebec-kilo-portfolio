
export type SupportedCurrency = 'USD' | 'EUR' | 'GBP' | 'JPY' | 'KRW' | 'CNY';
export type CurrencyAmount = number | null | undefined;
const CURRENCY_CONFIGS = {
  USD: { minDigits: 2, maxDigits: 2 },
  EUR: { minDigits: 2, maxDigits: 2 },
  GBP: { minDigits: 2, maxDigits: 2 },
  JPY: { minDigits: 0, maxDigits: 0 },
  KRW: { minDigits: 0, maxDigits: 0 },
  CNY: { minDigits: 2, maxDigits: 2 },
} as const;

// PERF: Pre-cache common formatters
const COMMON_FORMATTERS = {
  'en-US-USD': new Intl.NumberFormat('en-US', {
    style: 'currency', // NOTE: This tells Intl.NumberFormat the kind of formatting we need
    currency: 'USD',
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  }),
  'ja-JP-JPY': new Intl.NumberFormat('ja-JP', {
    style: 'currency',
    currency: 'JPY',
    minimumFractionDigits: 0,
    maximumFractionDigits: 0,
  }),
} as const;

const currencyFormatters = new Map<string, Intl.NumberFormat>();

const getFormatter = (currency: string, locale: string): Intl.NumberFormat => {
  const key = `${locale}-${currency}`;

  // Check pre-cached formatters first
  if (key in COMMON_FORMATTERS) {
    return COMMON_FORMATTERS[key as keyof typeof COMMON_FORMATTERS];
  }

  // Fallback to dynamic creation
  if (!currencyFormatters.has(key)) {
    const config = CURRENCY_CONFIGS[currency as keyof typeof CURRENCY_CONFIGS] || { minDigits: 2, maxDigits: 2 };

    currencyFormatters.set(key, new Intl.NumberFormat(locale, {
      style: 'currency',
      currency,
      minimumFractionDigits: config.minDigits,
      maximumFractionDigits: config.maxDigits,
    }));
  }

  return currencyFormatters.get(key)!;
};

/**
 * Format a number as USD currency
 * @param amount - Amount to format (number)
 * @returns Formatted currency string (e.g., "$52.80")
 */
export const formatUSD = (amount: CurrencyAmount): string => {
  if (amount == null) return '$0.00';
  if (amount <  0) return '$0.00';
  if (!Number.isFinite(amount)) return '$0.00';

  try {
    // Fallback formatting if for some reason Int.NumberFormat fails
    const formatter = getFormatter('USD', 'en-US');
    return formatter.format(amount);
  } catch (error) {
    console.warn('Currency formatting failed, using fallback: ', error);
    return `$${amount.toFixed(2)}`;
  }
}


/**
 * EXPANSION: Format a number as currency with custom currency and locale
 * @param amount - Amount to format (number)
 * @param currency - Currency code (default: 'USD')
 * @param locale - Locale string (default: 'en-US')
 * @returns Formatted currency string
 */
export const formatCurrency = (
  amount: CurrencyAmount,
  currency: string = 'USD',
  locale: string = 'en-US'
): string => {
  if (amount == null) return formatUSD(0);
  if (amount < 0) return formatUSD(0);
  if (!Number.isFinite(amount)) return formatUSD(0);

  try {
    const formatter = getFormatter(currency, locale);
    return formatter.format(amount);
  } catch (error) {

    // Fallback to USD formatting if custom formatting fails
    console.warn(`Currency formatting failed for ${currency}-${locale}, using USD fallback:`, error);
    return formatUSD(amount);
  }
};

// Clear cached formatters for testing
export const clearCurrencyFormatters = (): void => currencyFormatters.clear();

// Get number of the above formatters in case we need to debug
export const getCurrencyFormatterCount = (): number => currencyFormatters.size;
