import type { BaseSpendTracking } from '@/features/dashboard/lib/types/spend-tracking/base';
import type { ItemCondition } from '@/shared/constants/service.constants';
import type { ISO8601Date } from '@/shared/types/types';
import type { OnlineServiceProviderId, OnlineServiceProviderDisplay } from '@/shared/constants/service.constants';

/**
 * Represents different types of in-game purchase content
 * Used to categorize microtransactions and virtual goods
 */
export const IN_GAME_PURCHASE_TYPES = {
  CURRENCY: 'currency',
  ITEM: 'item',
  COSMETIC: 'cosmetic',
  FEATURE: 'feature'
} as const;

export type InGamePurchaseType = typeof IN_GAME_PURCHASE_TYPES[keyof typeof IN_GAME_PURCHASE_TYPES];

/**
 * Base interface for one-time purchases
 * Extends the base tracking with fields common to all one-time purchases
 */
export interface OneTimeSpend extends BaseSpendTracking {
  spendTransactionType: 'one-time';
  inGamePurchaseType?: InGamePurchaseType;
  isDigital: boolean;
  isWishlisted: boolean;
  itemCondition?: ItemCondition;
  provider?: {
    id: OnlineServiceProviderId;
    displayName: OnlineServiceProviderDisplay;
  };
  purchaseDate: ISO8601Date;
}
