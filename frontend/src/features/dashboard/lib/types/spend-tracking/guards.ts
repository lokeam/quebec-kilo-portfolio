import { PURCHASED_MEDIA_CATEGORIES } from '@/features/dashboard/lib/types/spend-tracking/media';
import type { SubscriptionSpend } from '@/features/dashboard/lib/types/spend-tracking/subscription';
import type { OneTimeSpend } from '@/features/dashboard/lib/types/spend-tracking/purchases';

export function isSubscriptionSpend(
  item: SubscriptionSpend | OneTimeSpend
): item is SubscriptionSpend {
  return item.spendTransactionType === PURCHASED_MEDIA_CATEGORIES.SUBSCRIPTION;
}
