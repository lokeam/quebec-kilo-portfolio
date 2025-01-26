/**
 * Defines the fundamental types of media items in our library system
 * This discriminator affects how items are processed, stored, and displayed
 */
export const MEDIA_ITEM_TYPES = {
  PHYSICAL: 'physical',
  DIGITAL: 'digital',
} as const;

/**
 * Type derived from our media item type constants
 * Used to ensure type safety when discriminating between physical and digital items
 */
export type MediaItemType = typeof MEDIA_ITEM_TYPES[keyof typeof MEDIA_ITEM_TYPES];
