import { type LibraryItem, type DigitalLibraryItem, type PhysicalLibraryItem } from '@/features/dashboard/lib/types/library/items';
import { MEDIA_ITEM_TYPES } from '@/features/dashboard/lib/types/library/constants';

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
