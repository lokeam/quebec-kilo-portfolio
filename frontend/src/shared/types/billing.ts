import type { BillingCycle, PaymentMethod } from "@/shared/constants/payment";
import type { Currency } from "@/shared/types/types";

export interface BillingDetails {
  cycle: BillingCycle;
  fees: {
    monthly: Currency;
    quarterly: Currency;
    annual: Currency;
  };
  renewalDate: {
    day: string;
    month: string;
  };
  paymentMethod?: PaymentMethod;
  lastBilledAt?: string;
  nextBillingDate?: string;
}