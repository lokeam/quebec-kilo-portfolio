import { type LibraryItem, type DigitalLibraryItem, type PhysicalLibraryItem } from '@/features/dashboard/lib/types/library/items';
import { MEDIA_ITEM_TYPES } from '@/features/dashboard/lib/types/library/constants';
import { NOTIFICATION_ICONS, type NotificationIcon } from '@/features/dashboard/lib/types/notifications/constants';

/**
 * Type guard to check if a library item is a digital item
 * @param item - The library item to check
 * @returns True if the item is a digital library item
 *
 * @example
 * ```ts
 * const item: LibraryItem = getLibraryItem();
 * if (isDigitalLibraryItem(item)) {
 *   // item is now typed as DigitalLibraryItem
 *   console.log(item.location.service);
 * }
 * ```
 */
export function isDigitalLibraryItem(item: LibraryItem): item is DigitalLibraryItem {
  return item.type === MEDIA_ITEM_TYPES.DIGITAL;
}

/**
 * Type guard to check if a library item is a physical item
 * @param item - The library item to check
 * @returns True if the item is a physical library item
 *
 * @example
 * ```ts
 * const item: LibraryItem = getLibraryItem();
 * if (isPhysicalLibraryItem(item)) {
 *   // item is now typed as PhysicalLibraryItem
 *   console.log(item.location.sublocation);
 * }
 * ```
 */
export function isPhysicalLibraryItem(item: LibraryItem): item is PhysicalLibraryItem {
  return item.type === MEDIA_ITEM_TYPES.PHYSICAL;
}

/**
 * Transforms a notification icon string into its corresponding icon component identifier
 * @param icon - The notification icon type from the API
 * @returns The icon identifier used by the UI component system
 *
 * @example
 * ```ts
 * const iconName = transformNotificationIcon(NOTIFICATION_ICONS.CHECK);
 * // Returns 'CheckIcon' or similar component identifier
 * ```
 */
export function transformNotificationIcon(icon: NotificationIcon): string {
  const iconMap: Record<NotificationIcon, string> = {
    [NOTIFICATION_ICONS.CHECK]: 'check',
    [NOTIFICATION_ICONS.TAG]: 'tag',
    [NOTIFICATION_ICONS.BAR_CHART]: 'barChart',
    [NOTIFICATION_ICONS.ALERT_TRIANGLE]: 'alertTriangle',
  };

  return iconMap[icon] || 'QuestionMarkIcon'; // Fallback icon for unexpected values
}
