/**
 * Types used across the spend tracking feature
 *
 * Used in:
 * - MonthlySpendingAccordion.tsx
 * - MonthlySpendingAccordionItem.tsx
 * - MonthlySpendingItemDetails.tsx
 * - useSpendingData.ts
 *
 *
 *
 * IMPORTANT: LEGACY SPEND-TRACKING TYPES. DO NOT USE. THIS FILE IS MARKED FOR DELETION.
 *
 *
 */

import type { SingleYearlyTotalBFFResponse } from "@/types/domain/spend-tracking";

/**
 * Represents a spending item, either a subscription or one-time purchase
 */
export type SpendItem = {
  id: string;
  title: string;
  amount: number;
  spendTransactionType: SpendTransactionType;
  paymentMethod: string;
  mediaType: MediaType;
  serviceName?: {
    id: string;
    displayName: string;
  };
  createdAt: number;
  updatedAt: number;
  isActive: boolean;

  // Optional fields that may or may not exist
  billingCycle?: string;
  nextBillingDate?: number;
  yearlySpending?: Array<SingleYearlyTotalBFFResponse>;
  isDigital?: boolean;
  isWishlisted?: boolean;
  purchaseDate?: number;
}

/**
 * Represents yearly spending data
 * Used in:
 * - MonthlySpendingAccordion.tsx
 * - MonthlySpendingItemDetails.tsx
 * - useSpendingData.ts
 */
export type SingleYearlyTotalBFFResponse = {
  year: number;
  amount: number;
}

/**
 * Valid transaction types for spending items
 */
export type SpendTransactionType = 'subscription' | 'one-time';

/**
 * Valid media types for spending items
 */
export type MediaType = 'subscription' | 'dlc' | 'inGamePurchase' | 'disc' | 'hardware';