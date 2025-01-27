import type { BillingCycle, PaymentMethod } from "@/shared/constants/payment";
import type { Currency } from "@/shared/types/types";

/**
 * Represents complete billing information for a subscription
 * @interface BillingDetails
 */
export interface BillingDetails {
  /** Current billing cycle selection (e.g., monthly, quarterly, annual) */
  cycle: BillingCycle;

  /**
   * Fee structure for different billing cycles
   * @property {Currency} monthly - Monthly subscription fee
   * @property {Currency} quarterly - Quarterly subscription fee
   * @property {Currency} annual - Annual subscription fee
   */
  fees: {
    monthly: Currency;
    quarterly: Currency;
    annual: Currency;
  };

  /**
   * Subscription renewal date information
   * @property {string} day - Day of renewal
   * @property {string} month - Month of renewal
   */
  renewalDate: {
    day: string;
    month: string;
  };

  /**
   * Current payment method used for billing
   * @optional
   */
  paymentMethod?: PaymentMethod;

  /**
   * ISO date string of the last successful billing
   * @optional
   */
  lastBilledAt?: string;

  /**
   * ISO date string of the next scheduled billing
   * @optional
   */
  nextBillingDate?: string;
}
