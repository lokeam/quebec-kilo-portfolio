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

export const spendTrackingPageMockData = {
  totalMonthlySpending: {
    currentMonthTotal: 1784.04,
    lastMonthTotal: 2255.92,
    percentageChange: -20.91,
    comparisonDateRange: "Dec 1 - Dec 22, 2024",
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
    dateRange: "January 2024 - January 2025",
    monthlyExpenditures: [
      { month: "Jan", expenditure: 450 },
      { month: "Feb", expenditure: 380 },
      { month: "Mar", expenditure: 420 },
      { month: "Apr", expenditure: 390 },
      { month: "May", expenditure: 410 },
      { month: "Jun", expenditure: 430 },
      { month: "Jul", expenditure: 400 },
      { month: "Aug", expenditure: 440 },
      { month: "Sep", expenditure: 370 },
      { month: "Oct", expenditure: 460 },
      { month: "Nov", expenditure: 420 },
      { month: "Dec", expenditure: 390 },
    ],
    medianMonthlyCost: 410
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
      createdAt: 1704067200000,
      updatedAt: 1704067200000,
      billingCycle: 'quarterly',
      nextBillingDate: 1743638400000,
      isActive: true,
      yearlySpending: [
        { year: 2022, amount: 79.92 },
        { year: 2023, amount: 79.92 },
        { year: 2024, amount: 79.92 }
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
      createdAt: 1704067200000,
      updatedAt: 1704067200000,
      billingCycle: 'monthly',
      nextBillingDate: 1743638400000,
      isActive: true,
      yearlySpending: [
        { year: 2022, amount: 179.88 },
        { year: 2023, amount: 179.88 },
        { year: 2024, amount: 179.88 }
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
      createdAt: 1704067200000,
      updatedAt: 1704067200000,
      billingCycle: 'annual',
      nextBillingDate: 1743638400000,
      isActive: true,
      yearlySpending: [
        { year: 2022, amount: 47.88 },
        { year: 2023, amount: 47.88 },
        { year: 2024, amount: 47.88 }
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
      createdAt: 1704067200000,
      updatedAt: 1704067200000,
      billingCycle: 'monthly',
      nextBillingDate: 1743638400000,
      isActive: true,
      yearlySpending: [
        { year: 2022, amount: 83.88 },
        { year: 2023, amount: 83.88 },
        { year: 2024, amount: 83.88 }
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
      createdAt: 1733529600000,
      updatedAt: 1733529600000,
      isActive: true,
      isDigital: true,
      isWishlisted: false,
      purchaseDate: 1733529600000,
    } as OneTimeSpend,
    {
      id: 'one-006',
      title: 'Path of Exile 2',
      amount: 20.00,
      spendTransactionType: TransactionType.ONE_TIME,
      createdAt: 1704844800000,
      updatedAt: 1704844800000,
      paymentMethod: 'PayPal',
      provider: 'steam',
      isActive: true,
      isDigital: true,
      mediaType: MediaCategory.IN_GAME_PURCHASE,
      isWishlisted: true,
      purchaseDate: 1733616000000,
    } as OneTimeSpend,
    {
      id: 'one-007',
      title: 'Sid Meier\'s Civilization VI',
      amount: 19.99,
      spendTransactionType: TransactionType.ONE_TIME,
      paymentMethod: 'PayPal',
      provider: 'steam',
      createdAt: 1733702400000,
      updatedAt: 1733702400000,
      isActive: true,
      isDigital: true,
      mediaType: MediaCategory.DLC,
      isWishlisted: true,
      purchaseDate: 1733702400000,
    } as OneTimeSpend,
    {
      id: 'one-008',
      title: 'ELDEN RING Shadow of the Erdtree',
      amount: 39.99,
      spendTransactionType: TransactionType.ONE_TIME,
      paymentMethod: 'PayPal',
      provider: 'steam',
      createdAt: 1733270400000,
      updatedAt: 1733270400000,
      isActive: true,
      isDigital: true,
      mediaType: MediaCategory.DLC,
      isWishlisted: true,
      purchaseDate: 1733270400000,
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
      createdAt: 1733788800000,
      updatedAt: 1733788800000,
      purchaseDate: 1733788800000,
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
      createdAt: 1734393600000,
      updatedAt: 1734393600000,
      purchaseDate: 1734393600000,
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
      createdAt: 1734393600000,
      updatedAt: 1734393600000,
      purchaseDate: 1734393600000,
    } as OneTimeSpend
  ] as SpendingItemBFFResponse[],
  oneTimeThisMonth: [] as SpendingItemBFFResponse[],
  recurringNextMonth: [
    {
      id: 'sub-012',
      title: 'Google Play Pass',
      provider: 'google',
      amount: 5.99,
      billingCycle: 'monthly',
      spendTransactionType: TransactionType.SUBSCRIPTION,
      paymentMethod: 'Mastercard',
      isActive: true,
      createdAt: 1733270400000,
      updatedAt: 1733270400000,
      nextBillingDate: 1735948800000,
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
      { year: 2022, amount: 399.88 },
      { year: 2023, amount: 459.96 },
      { year: 2024, amount: 499.92 }
    ],
    oneTimeTotal: [
      { year: 2022, amount: 899.95 },
      { year: 2023, amount: 1299.97 },
      { year: 2024, amount: 799.96 }
    ],
    combinedTotal: [
      { year: 2022, amount: 1299.83 },
      { year: 2023, amount: 1759.93 },
      { year: 2024, amount: 1299.88 }
    ]
  }
};
