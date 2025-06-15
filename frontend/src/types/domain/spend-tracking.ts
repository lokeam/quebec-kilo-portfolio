/**
 * Core types for spend tracking functionality
 *
 * Used in:
 * - MonthlySpendingAccordion.tsx
 * - MonthlySpendingAccordionItem.tsx
 * - MonthlySpendingItemDetails.tsx
 * - useSpendingData.ts
 */

// Enums
export enum TransactionType {
  SUBSCRIPTION = 'subscription',
  ONE_TIME = 'one-time'
}

export enum MediaCategory {
  HARDWARE = 'hardware',
  DLC = 'dlc',
  IN_GAME_PURCHASE = 'inGamePurchase',
  SUBSCRIPTION = 'subscription',
  PHYSICAL = 'physical',
  DISC = 'disc'
}

export enum BillingCycle {
  MONTHLY = 'monthly',
  QUARTERLY = 'quarterly',
  ANNUAL = 'annual'
}

// Core Types
export interface SpendingItem {
  id: string;
  title: string;
  amount: number;
  spendTransactionType: TransactionType;
  mediaType: MediaCategory;
  paymentMethod: string;
  serviceName?: {
    id: string;
    displayName: string;
  };
  provider?: string;
  isActive: boolean;
  createdAt: number;
  updatedAt: number;

  // Optional fields based on transaction type
  billingCycle?: BillingCycle;
  nextBillingDate?: number;
  purchaseDate?: number;
  isDigital?: boolean;
  isWishlisted?: boolean;
  yearlySpending?: YearlySpending[];
}

export interface YearlySpending {
  year: number;
  amount: number;
}

export interface SpendingCategory {
  name: string;
  value: number;
}

export interface MonthlyExpenditure {
  month: string;
  expenditure: number;
}

// Type Guards
export function isSubscriptionSpend(item: SpendingItem): boolean {
  return item.spendTransactionType === TransactionType.SUBSCRIPTION;
}

export function isOneTimeSpend(item: SpendingItem): boolean {
  return item.spendTransactionType === TransactionType.ONE_TIME;
}

// Utility Types
export type Money = number;
export type ISO8601Date = string;

// Type aliases for backward compatibility
export type SpendItem = SpendingItem;
export type MediaType = MediaCategory;
export type SpendTransactionType = TransactionType;
