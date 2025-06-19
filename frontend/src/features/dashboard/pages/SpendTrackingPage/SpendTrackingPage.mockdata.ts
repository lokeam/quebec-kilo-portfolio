import {
  MediaCategory,
  TransactionType,
} from '@/types/domain/spend-tracking';

import type {
  SpendingItemBFFResponse,
  SingleYearlyTotalBFFResponse
} from '@/types/domain/spend-tracking';

interface BaseSpendItem {
  id: string;
  title: string;
  amount: number;
  spendTransactionType: TransactionType;
  paymentMethod: string;
  mediaType: MediaCategory;
  provider: string;
  createdAt: number;
  updatedAt: number;
  isActive: boolean;
}

interface SubscriptionSpend extends BaseSpendItem {
  spendTransactionType: TransactionType.SUBSCRIPTION;
  billingCycle: string;
  nextBillingDate: number;
  yearlySpending: SingleYearlyTotalBFFResponse[];
}

interface OneTimeSpend extends BaseSpendItem {
  spendTransactionType: TransactionType.ONE_TIME;
  isDigital: boolean;
  isWishlisted: boolean;
  purchaseDate: number;
}

export const BASE_MEDIA_CATEGORIES = {
  HARDWARE: MediaCategory.HARDWARE,
  DLC: MediaCategory.DLC,
  IN_GAME_PURCHASE: MediaCategory.IN_GAME_PURCHASE,
  SUBSCRIPTION: MediaCategory.SUBSCRIPTION,
  PHYSICAL: MediaCategory.PHYSICAL,
  DISC: MediaCategory.DISC
} as const;

/**
 * NOTE:
 * Subscription dates follow the Anchor Date Pattern:
 * - createdAt: Represents the anchor_date (first payment date)
 * - nextBillingDate: Computed date based on anchor_date + billing_cycle
 * - Example: anchor_date = Jan 5, 2025, billing_cycle = '3 month' → nextBillingDate = Apr 5, 2025
 */


export const spendTrackingPageMockData = {
  totalMonthlySpending: {
    currentMonthTotal: 1784.04,
    lastMonthTotal: 879.37,
    percentageChange: -50.71,
    comparisonDateRange: "Jun 1 - Jun 22, 2025",
    spendingCategories: [
      { name: MediaCategory.HARDWARE, value: 534.04 },
      { name: MediaCategory.DLC, value: 267.02 },
      { name: MediaCategory.IN_GAME_PURCHASE, value: 178.01 },
      { name: MediaCategory.SUBSCRIPTION, value: 356.02 },
      { name: MediaCategory.PHYSICAL, value: 267.02 },
      { name: MediaCategory.DISC, value: 178.01 },
    ]
  },
  totalAnnualSpending: {
    dateRange: "January 2025 - January 2026",
    monthlyExpenditures: [
      { month: "Jan", expenditure: 201.65 },
      { month: "Feb", expenditure: 10.99 },
      { month: "Mar", expenditure: 39.55 },
      { month: "Apr", expenditure: 25.99 },
      { month: "May", expenditure: 879.37 },
      { month: "Jun", expenditure: 1784.0 },
      { month: "Jul", expenditure: 0 },
      { month: "Aug", expenditure: 0 },
      { month: "Sep", expenditure: 0 },
      { month: "Oct", expenditure: 0 },
      { month: "Nov", expenditure: 0 },
      { month: "Dec", expenditure: 0 },
    ],
    medianMonthlyCost: 490.25 // we need to calculate this
  },
  currentTotalThisMonth: [
    {
      id: 'sub-001',
      title: 'Playstation Plus',
      amount: 6.66,
      spendTransactionType: TransactionType.SUBSCRIPTION,
      paymentMethod: 'Visa',
      mediaType: MediaCategory.SUBSCRIPTION,
      provider: 'playstation',
      createdAt: 1748836800000,
      updatedAt: 1748836800000,
      billingCycle: '3 month',
      nextBillingDate: 1751428800000,
      isActive: true,
      yearlySpending: [
        { year: 2023, amount: 79.92 },
        { year: 2024, amount: 79.92 },
        { year: 2025, amount: 79.92 }
      ]
    } as SubscriptionSpend,
    {
      id: 'sub-002',
      title: 'Xbox Game Pass Ultimate',
      amount: 14.99,
      spendTransactionType: TransactionType.SUBSCRIPTION,
      paymentMethod: 'Mastercard',
      mediaType: MediaCategory.SUBSCRIPTION,
      provider: 'xbox',
      createdAt: 1749096000000,
      updatedAt: 1751688000000,
      billingCycle: '1 month',
      nextBillingDate: 1751688000000,
      isActive: true,
      yearlySpending: [
        { year: 2023, amount: 179.88 },
        { year: 2024, amount: 179.88 },
        { year: 2025, amount: 179.88 }
      ]
    } as SubscriptionSpend,
    {
      id: 'sub-003',
      title: 'Nintendo Switch Online',
      amount: 3.99,
      spendTransactionType: TransactionType.SUBSCRIPTION,
      paymentMethod: 'Visa',
      mediaType: MediaCategory.SUBSCRIPTION,
      provider: 'nintendo',
      createdAt: 1749355200000,
      updatedAt: 1749355200000,
      billingCycle: '12 month',
      nextBillingDate: 1751947200000,
      isActive: true,
      yearlySpending: [
        { year: 2023, amount: 47.88 },
        { year: 2024, amount: 47.88 },
        { year: 2025, amount: 47.88 }
      ]
    } as SubscriptionSpend,
    {
      id: 'sub-004',
      title: 'Apple Arcade',
      amount: 6.99,
      spendTransactionType: TransactionType.SUBSCRIPTION,
      paymentMethod: 'Mastercard',
      mediaType: MediaCategory.SUBSCRIPTION,
      provider: 'apple',
      createdAt: 1749441600000,
      updatedAt: 1749441600000,
      billingCycle: '1 month',
      nextBillingDate: 1752033600000,
      isActive: true,
      yearlySpending: [
        { year: 2023, amount: 83.88 },
        { year: 2024, amount: 83.88 },
        { year: 2025, amount: 83.88 }
      ]
    } as SubscriptionSpend,
    {
      id: 'one-005',
      title: 'Helldivers 2',
      amount: 10.99,
      spendTransactionType: TransactionType.ONE_TIME,
      paymentMethod: 'PayPal',
      mediaType: MediaCategory.IN_GAME_PURCHASE,
      provider: 'steam',
      createdAt: 1752120000000,
      updatedAt: 1752120000000,
      isActive: true,
      isDigital: true,
      isWishlisted: false,
      purchaseDate: 1752120000000,
    } as OneTimeSpend,
    {
      id: 'one-006',
      title: 'Path of Exile 2',
      amount: 20.00,
      spendTransactionType: TransactionType.ONE_TIME,
      createdAt: 1752120000000,
      updatedAt: 1752120000000,
      paymentMethod: 'PayPal',
      provider: 'steam',
      isActive: true,
      isDigital: true,
      mediaType: MediaCategory.IN_GAME_PURCHASE,
      isWishlisted: true,
      purchaseDate: 1752120000000,
    } as OneTimeSpend,
    {
      id: 'one-007',
      title: 'Sid Meier\'s Civilization VI',
      amount: 19.99,
      spendTransactionType: TransactionType.ONE_TIME,
      paymentMethod: 'PayPal',
      provider: 'steam',
      createdAt: 1752292800000,
      updatedAt: 1752292800000,
      isActive: true,
      isDigital: true,
      mediaType: MediaCategory.DLC,
      isWishlisted: true,
      purchaseDate: 1752292800000,
    } as OneTimeSpend,
    {
      id: 'one-008',
      title: 'ELDEN RING Shadow of the Erdtree',
      amount: 39.99,
      spendTransactionType: TransactionType.ONE_TIME,
      paymentMethod: 'PayPal',
      provider: 'steam',
      createdAt: 1752379200000,
      updatedAt: 1752379200000,
      isActive: true,
      isDigital: true,
      mediaType: MediaCategory.DLC,
      isWishlisted: true,
      purchaseDate: 1752379200000,
    } as OneTimeSpend,
    {
      id: 'one-009',
      title: 'The Legend of Zelda: Breath of the Wild',
      amount: 29.99,
      spendTransactionType: TransactionType.ONE_TIME,
      paymentMethod: 'Visa',
      isActive: true,
      isDigital: false,
      isWishlisted: true,
      mediaType: MediaCategory.DISC,
      createdAt: 1752465600000,
      updatedAt: 1752465600000,
      purchaseDate: 1752465600000,
    } as OneTimeSpend,
    {
      id: 'one-010',
      title: 'Gradius V',
      amount: 59.99,
      spendTransactionType: TransactionType.ONE_TIME,
      paymentMethod: 'Visa',
      isDigital: false,
      isWishlisted: true,
      mediaType: MediaCategory.DISC,
      createdAt: 1752552000000,
      updatedAt: 1752552000000,
      purchaseDate: 1752552000000,
    } as OneTimeSpend,
    {
      id: 'one-011',
      title: 'G.Skill Trident Z5 Neo RGB DDR5-6000',
      amount: 219.98,
      spendTransactionType: TransactionType.ONE_TIME,
      paymentMethod: 'Visa',
      isDigital: false,
      isWishlisted: true,
      mediaType: MediaCategory.HARDWARE,
      createdAt: 1752552000000,
      updatedAt: 1752552000000,
      purchaseDate: 1752552000000,
    } as OneTimeSpend
  ] as SpendingItemBFFResponse[],
  oneTimeThisMonth: [
    {
      id: 'one-005',
      title: 'Helldivers 2',
      amount: 10.99,
      spendTransactionType: TransactionType.ONE_TIME,
      paymentMethod: 'PayPal',
      mediaType: MediaCategory.IN_GAME_PURCHASE,
      provider: 'steam',
      createdAt: 1752120000000,
      updatedAt: 1752120000000,
      isActive: true,
      isDigital: true,
      isWishlisted: false,
      purchaseDate: 1752120000000,
    } as OneTimeSpend,
    {
      id: 'one-006',
      title: 'Path of Exile 2',
      amount: 20.00,
      spendTransactionType: TransactionType.ONE_TIME,
      createdAt: 1752120000000,
      updatedAt: 1752120000000,
      paymentMethod: 'PayPal',
      provider: 'steam',
      isActive: true,
      isDigital: true,
      mediaType: MediaCategory.IN_GAME_PURCHASE,
      isWishlisted: true,
      purchaseDate: 1752120000000,
    } as OneTimeSpend,
    {
      id: 'one-007',
      title: 'Sid Meier\'s Civilization VI',
      amount: 19.99,
      spendTransactionType: TransactionType.ONE_TIME,
      paymentMethod: 'PayPal',
      provider: 'steam',
      createdAt: 1752292800000,
      updatedAt: 1752292800000,
      isActive: true,
      isDigital: true,
      mediaType: MediaCategory.DLC,
      isWishlisted: true,
      purchaseDate: 1752292800000,
    } as OneTimeSpend,
    {
      id: 'one-008',
      title: 'ELDEN RING Shadow of the Erdtree',
      amount: 39.99,
      spendTransactionType: TransactionType.ONE_TIME,
      paymentMethod: 'PayPal',
      provider: 'steam',
      createdAt: 1752379200000,
      updatedAt: 1752379200000,
      isActive: true,
      isDigital: true,
      mediaType: MediaCategory.DLC,
      isWishlisted: true,
      purchaseDate: 1752379200000,
    } as OneTimeSpend,
    {
      id: 'one-009',
      title: 'The Legend of Zelda: Breath of the Wild',
      amount: 29.99,
      spendTransactionType: TransactionType.ONE_TIME,
      paymentMethod: 'Visa',
      isActive: true,
      isDigital: false,
      isWishlisted: true,
      mediaType: MediaCategory.DISC,
      createdAt: 1752465600000,
      updatedAt: 1752465600000,
      purchaseDate: 1752465600000,
    } as OneTimeSpend,
    {
      id: 'one-010',
      title: 'Gradius V',
      amount: 59.99,
      spendTransactionType: TransactionType.ONE_TIME,
      paymentMethod: 'Visa',
      isDigital: false,
      isWishlisted: true,
      mediaType: MediaCategory.DISC,
      createdAt: 1752552000000,
      updatedAt: 1752552000000,
      purchaseDate: 1752552000000,
    } as OneTimeSpend,
    {
      id: 'one-011',
      title: 'G.Skill Trident Z5 Neo RGB DDR5-6000',
      amount: 219.98,
      spendTransactionType: TransactionType.ONE_TIME,
      paymentMethod: 'Visa',
      isDigital: false,
      isWishlisted: true,
      mediaType: MediaCategory.HARDWARE,
      createdAt: 1752552000000,
      updatedAt: 1752552000000,
      purchaseDate: 1752552000000,
    } as OneTimeSpend
  ] as SpendingItemBFFResponse[],
  recurringNextMonth: [
    {
      id: 'sub-012',
      title: 'Google Play Pass',
      provider: 'google',
      amount: 5.99,
      billingCycle: '1 month',
      spendTransactionType: TransactionType.SUBSCRIPTION,
      paymentMethod: 'Mastercard',
      isActive: true,
      createdAt: 1740027600000,
      updatedAt: 1750392000000,
      nextBillingDate: 1752984000000,
      mediaType: MediaCategory.SUBSCRIPTION,
      yearlySpending: [
        { year: 2022, amount: 71.88 },
        { year: 2023, amount: 71.88 },
        { year: 2024, amount: 71.88 }
      ]
    } as SubscriptionSpend
  ] as SpendingItemBFFResponse[],
  yearlyTotals: {
    subscriptionTotal: [
      { year: 2023, amount: 399.88 },
      { year: 2024, amount: 459.96 },
      { year: 2025, amount: 38.62 }
    ],
    oneTimeTotal: [
      { year: 2023, amount: 899.95 },
      { year: 2024, amount: 1299.97 },
      { year: 2025, amount: 601.25 }
    ],
    combinedTotal: [
      { year: 2023, amount: 1299.83 },
      { year: 2024, amount: 1759.93 },
      { year: 2025, amount: 639.87 }
    ]
  }
};
