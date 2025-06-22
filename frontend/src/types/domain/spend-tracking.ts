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
  PHYSICAL = 'physical_game',
  DIGITAL = 'digital_game',
  MISC = 'misc'
}

export enum BillingCycle {
  MONTHLY = 'monthly',
  QUARTERLY = 'quarterly',
  ANNUAL = 'annual'
}

// Core Types
// original name: SpendingItem
export interface SpendingItemBFFResponse {
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
  yearlySpending?: SingleYearlyTotalBFFResponse[];
}

// original name: YearlySpending
// Linked to SingleYearlyTotalBFFResponse in spend_tracking_response_types.go
export interface SingleYearlyTotalBFFResponse {
  year: number;
  amount: number;
}

// original name: SpendingCategory
// Linked to SpendingCategoryBFFResponseFINAL in spend_tracking_response_types.go
export interface SpendingCategoryBFFResponse {
  name: string;
  value: number;
}

// original name: MonthlyExpenditure
// Linked to MonthlyExpenditureBFFResponseFINAL in spend_tracking_response_types.go
export interface SingleMonthlyExpenditureBFFResponse {
  month: string;
  expenditure: number;
}

// Linked to MonthlySpendingBFFResponseFINAL in spend_tracking_response_types.go
export interface MonthlySpendingBFFResponse {
  currentMonthTotal: number;
  lastMonthTotal: number;
  percentageChange: number;
  comparisonDateRange: string;
  spendingCategories: SpendingCategoryBFFResponse[];
}

// Linked to AnnualSpendingBFFResponseFINAL in spend_tracking_response_types.go
export interface AnnualSpendingBFFResponse {
  dateRange: string;
  monthlyExpenditures: SingleMonthlyExpenditureBFFResponse[];
  medianMonthlyCost: number;
}

// Linked to AllYearlyTotalsBFFResponseFINAL in spend_tracking_response_types.go
export interface AllYearlyTotalsBFFResponse {
  subscriptionTotal: SingleYearlyTotalBFFResponse[];
  oneTimeTotal: SingleYearlyTotalBFFResponse[];
  combinedTotal: SingleYearlyTotalBFFResponse[];
}

export interface SpendTrackingBFFResponse {
  totalMonthlySpending: MonthlySpendingBFFResponse;
  totalAnnualSpending: AnnualSpendingBFFResponse;
  currentTotalThisMonth: SpendingItemBFFResponse[];
  oneTimeThisMonth: SpendingItemBFFResponse[];
  recurringNextMonth: SpendingItemBFFResponse[];
  yearlyTotals: AllYearlyTotalsBFFResponse;
};


// Type Guards
export function isSubscriptionSpend(item: SpendingItemBFFResponse): boolean {
  return item.spendTransactionType === TransactionType.SUBSCRIPTION;
}

export function isOneTimeSpend(item: SpendingItemBFFResponse): boolean {
  return item.spendTransactionType === TransactionType.ONE_TIME;
}

// Utility Types
export type Money = number;
export type ISO8601Date = string;

// Type aliases for backward compatibility
export type SpendItem = SpendingItemBFFResponse;
export type MediaType = MediaCategory;
export type SpendTransactionType = TransactionType;
