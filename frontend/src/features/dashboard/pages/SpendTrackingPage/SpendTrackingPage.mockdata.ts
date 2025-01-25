

import type { SpendTrackingData } from '@/features/dashboard/lib/types/spend-tracking/aggregates';
import { ONLINE_SERVICE_PROVIDERS } from '@/shared/constants/service.constants';
import { PURCHASED_MEDIA_CATEGORIES } from '@/features/dashboard/lib/types/spend-tracking/media';
import { BILLING_CYCLES, PAYMENT_METHODS } from '@/features/dashboard/lib/types/spend-tracking/constants';

export const spendTrackingPageMockData: SpendTrackingData = {
  currentTotalThisMonth: [
    {
      id: 'sub-001',
      title: 'Playstation Plus',
      amount: '6.66',
      spendTransactionType: 'subscription',
      paymentMethod: PAYMENT_METHODS.VISA,
      mediaType: PURCHASED_MEDIA_CATEGORIES.SUBSCRIPTION,
      provider: ONLINE_SERVICE_PROVIDERS.SONY,
      createdAt: '2024-01-01',
      updatedAt: '2024-01-01',
      billingCycle: BILLING_CYCLES.QUARTERLY,
      nextBillingDate: '2025-04-01',
      isActive: true,
      yearlySpending: [
        { year: 2022, amount: '79.92' },  // $6.66 * 4 quarters
        { year: 2023, amount: '79.92' },
        { year: 2024, amount: '79.92' }
      ]
    },
    {
      id: 'sub-002',
      title: 'Xbox Game Pass Ultimate',
      amount: '14.99',
      spendTransactionType: 'subscription',
      paymentMethod: PAYMENT_METHODS.MASTERCARD,
      mediaType: PURCHASED_MEDIA_CATEGORIES.SUBSCRIPTION,
      provider: ONLINE_SERVICE_PROVIDERS.MICROSOFT,
      createdAt: '2024-01-01',
      updatedAt: '2024-01-01',
      billingCycle: BILLING_CYCLES.MONTHLY,
      nextBillingDate: '2025-04-01',
      isActive: true,
      yearlySpending: [
        { year: 2022, amount: '179.88' },  // $14.99 * 12 months
        { year: 2023, amount: '179.88' },
        { year: 2024, amount: '179.88' }
      ]
    },
    {
      id: 'sub-003',
      title: 'Nintendo Switch Online',
      amount: '3.99',
      spendTransactionType: 'subscription',
      paymentMethod: PAYMENT_METHODS.VISA,
      mediaType: PURCHASED_MEDIA_CATEGORIES.SUBSCRIPTION,
      provider: ONLINE_SERVICE_PROVIDERS.NINTENDO,
      createdAt: '2024-01-01',
      updatedAt: '2024-01-01',
      billingCycle: BILLING_CYCLES.ANNUAL,
      nextBillingDate: '2025-04-01',
      isActive: true,
      yearlySpending: [
        { year: 2022, amount: '47.88' },  // $3.99 * 12 months
        { year: 2023, amount: '47.88' },
        { year: 2024, amount: '47.88' }
      ]
    },
    {
      id: 'sub-004',
      title: 'Apple Arcade',
      amount: '6.99',
      spendTransactionType: 'subscription',
      paymentMethod: PAYMENT_METHODS.MASTERCARD,
      mediaType: PURCHASED_MEDIA_CATEGORIES.SUBSCRIPTION,
      provider: ONLINE_SERVICE_PROVIDERS.APPLE,
      createdAt: '2024-01-01',
      updatedAt: '2024-01-01',
      billingCycle: BILLING_CYCLES.MONTHLY,
      nextBillingDate: '2025-04-01',
      isActive: true,
      yearlySpending: [
        { year: 2022, amount: '83.88' },  // $6.99 * 12 months
        { year: 2023, amount: '83.88' },
        { year: 2024, amount: '83.88' }
      ]
    },
    {
      id: 'one-005',
      title: 'Helldivers 2',
      amount: '10.99',
      spendTransactionType: 'one-time',
      paymentMethod: PAYMENT_METHODS.PAYPAL,
      mediaType: PURCHASED_MEDIA_CATEGORIES.IN_GAME_PURCHASE,
      provider: ONLINE_SERVICE_PROVIDERS.STEAM,
      createdAt: '2025-01-05',
      updatedAt: '2025-01-05',
      isActive: true,
      isDigital: true,
      isWishlisted: false,
      purchaseDate: '2025-01-05',
    },
    {
      id: 'one-006',
      title: 'Path of Exile 2',
      amount: '20.00',
      spendTransactionType: 'one-time',
      createdAt: '2024-01-10',
      updatedAt: '2024-01-10',
      paymentMethod: PAYMENT_METHODS.PAYPAL,
      provider: ONLINE_SERVICE_PROVIDERS.STEAM,
      isActive: true,
      isDigital: true,
      mediaType: PURCHASED_MEDIA_CATEGORIES.IN_GAME_PURCHASE,
      isWishlisted: true,
      purchaseDate: '2025-01-10',
    },
    {
      id: 'one-007',
      title: 'Sid Meier\'s Civilization VI',
      amount: '19.99',
      spendTransactionType: 'one-time',
      paymentMethod: PAYMENT_METHODS.PAYPAL,
      provider: ONLINE_SERVICE_PROVIDERS.STEAM,
      createdAt: '2025-01-11',
      updatedAt: '2025-01-11',
      isActive: true,
      isDigital: true,
      mediaType: PURCHASED_MEDIA_CATEGORIES.DLC,
      isWishlisted: true,
      purchaseDate: '2025-01-11',
    },
    {
      id: 'one-008',
      title: 'ELDEN RING Shadow of the Erdtree',
      amount: '39.99',
      spendTransactionType: 'one-time',
      paymentMethod: PAYMENT_METHODS.PAYPAL,
      provider: ONLINE_SERVICE_PROVIDERS.STEAM,
      createdAt: '2025-01-02',
      updatedAt: '2025-01-02',
      isActive: true,
      isDigital: true,
      mediaType: PURCHASED_MEDIA_CATEGORIES.DLC,
      isWishlisted: true,
      purchaseDate: '2025-01-02',
    },
    {
      id: 'one-009',
      title: 'The Legend of Zelda: Breath of the Wild',
      amount: '29.99',
      spendTransactionType: 'one-time',
      paymentMethod: PAYMENT_METHODS.VISA,
      isActive: true,
      isDigital: false,
      isWishlisted: true,
      mediaType: PURCHASED_MEDIA_CATEGORIES.DISC,
      createdAt: '2025-01-12',
      updatedAt: '2025-01-12',
      purchaseDate: '2025-01-12',
    },
    {
      id: 'one-010',
      title: 'Gradius V',
      amount: '59.99',
      spendTransactionType: 'one-time',
      paymentMethod: PAYMENT_METHODS.VISA,
      isDigital: false,
      isWishlisted: true,
      mediaType: PURCHASED_MEDIA_CATEGORIES.DISC,
      createdAt: '2025-01-20',
      updatedAt: '2025-01-20',
      purchaseDate: '2025-01-20',
    },
    {
      id: 'one-011',
      title: 'G.Skill Trident Z5 Neo RGB DDR5-6000',
      amount: '219.98',
      spendTransactionType: 'one-time',
      paymentMethod: PAYMENT_METHODS.VISA,
      isDigital: false,
      isWishlisted: true,
      mediaType: PURCHASED_MEDIA_CATEGORIES.HARDWARE,
      createdAt: '2025-01-20',
      updatedAt: '2025-01-20',
      purchaseDate: '2025-01-20',
    }
  ],
  oneTimeThisMonth: [],
  recurringNextMonth: [
    {
      id: 'sub-012',
      title: 'Google Play Pass',
      provider: ONLINE_SERVICE_PROVIDERS.GOOGLE,
      amount: '5.99',
      billingCycle: BILLING_CYCLES.MONTHLY,
      spendTransactionType: 'subscription',
      paymentMethod: PAYMENT_METHODS.MASTERCARD,
      isActive: true,
      createdAt: '2025-01-01',
      updatedAt: '2025-01-01',
      nextBillingDate: '2025-02-01',
      mediaType: PURCHASED_MEDIA_CATEGORIES.SUBSCRIPTION,
    }
  ],
  totalSpendsThisMonth: '411.58',
  totalSpendsThisYear: '411.58',
  yearlyTotals: {
    subscriptionTotal: [
      { year: 2022, amount: '399.88' },
      { year: 2023, amount: '459.96' },
      { year: 2024, amount: '499.92' }
    ],
    oneTimeTotal: [
      { year: 2022, amount: '899.95' },
      { year: 2023, amount: '1299.97' },
      { year: 2024, amount: '799.96' }
    ],
    combinedTotal: [
      { year: 2022, amount: '1299.83' },
      { year: 2023, amount: '1759.93' },
      { year: 2024, amount: '1299.88' }
    ]
  }
}
