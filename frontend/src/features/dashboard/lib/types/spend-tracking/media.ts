import { BASE_MEDIA_CATEGORIES } from '@/features/dashboard/lib/types/spend-tracking/constants';

/**
 * Represents the types of a purchased item/piece of content
 * This helps with reporting and analytics
 */
export const PURCHASED_MEDIA_CATEGORIES = BASE_MEDIA_CATEGORIES;

export type PurchasedMediaCategory = typeof PURCHASED_MEDIA_CATEGORIES[keyof typeof PURCHASED_MEDIA_CATEGORIES];
