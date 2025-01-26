import type { ReportCategory } from '@/features/dashboard/lib/types/notifications/constants';
import type { BaseNotification } from '@/features/dashboard/lib/types/notifications/base';
import type { NOTIFICATION_CATEGORIES } from '@/features/dashboard/lib/types/notifications/constants';

/**
 * Payload for application update notifications
 */
export interface AppUpdatePayload {
  /** Version number of the update (optional) */
  version?: string;
  /** URL to the update info page / changelog */
  infoUrl: string;
  /** List of changes introduced in the update */
  changes: string[];
}

/**
 * Payload for report generation notifications
 */
export interface ReportPayload {
  type: ReportCategory;
  period: {
    month?: string;
    year: string;
  };
  downloadUrl: string;
  fileSize: string;
};

/**
 * Payload for wishlist item notifications (e.g., price drops, sales)
 */
export interface WishlistPayload {
  name: string;
  salePrice: string;
  originalPrice: string;
  discountPercentage: number;
  coverImageUrl?: string;
  storeName: string;
  storeUrl: string;
  /** Promotional code for additional savings (optional) */
  saleCode?: string;
  /** Date of the sale (optional) */
  saleDate?: string;
}

/**
 * Notification for application updates
 * Extends base notification with app update specific data
 */
export interface AppUpdateNotification extends BaseNotification {
  /** Discriminator for app update notifications */
  type: typeof NOTIFICATION_CATEGORIES.APP_UPDATE;
  /** App update specific information */
  update: AppUpdatePayload;
}

/**
 * Notification for generated reports
 * Extends base notification with report specific data
 */
export interface ReportNotification extends BaseNotification {
  /** Discriminator for report notifications */
  type: typeof NOTIFICATION_CATEGORIES.REPORT;
  /** Report specific information */
  report: ReportPayload;
}

/**
 * Notification for wishlist items
 * Extends base notification with wishlist specific data
 */
export interface WishlistNotification extends BaseNotification {
  /** Discriminator for wishlist notifications */
  type: typeof NOTIFICATION_CATEGORIES.WISHLIST;
  /** Wishlist item specific information */
  item: WishlistPayload;
}

/**
 * Union type of all possible notification variants
 * Used for type discrimination based on the 'type' property
 */
export type Notification =
  | AppUpdateNotification
  | ReportNotification
  | WishlistNotification;
