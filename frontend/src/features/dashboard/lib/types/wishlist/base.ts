import type { PlatformCategory } from '@/shared/types/platform';
import type { Price } from '@/shared/types/pricing';
import type { Rating } from '@/features/dashboard/lib/types/wishlist/ratings';

/**
 * Represents the core properties of a wishlist item
 * @interface BaseWishlistItem
 */
export interface BaseWishlistItem {
  /** Unique identifier for the wishlist item */
  readonly id: string;

  /** Title/name of the game or item */
  readonly title: string;

  /** URL to the item's thumbnail image */
  readonly thumbnailUrl: string;

  /** Array of descriptive tags for the item */
  readonly tags: ReadonlyArray<string>;

  /** Release date of the item
   * @remarks Consider using Date type in production environment
   */
  readonly releaseDate: string;

  /** Platform category the item belongs to */
  readonly platform: PlatformCategory;
}

/**
 * Extends BaseWishlistItem with additional properties for the full wishlist item
 * @interface WishlistItem
 * @extends {BaseWishlistItem}
 */
export interface WishlistItem extends BaseWishlistItem {
  /** Pricing information for the item */
  readonly price: Price;

  /** Optional platform compatibility information */
  readonly platformSupport?: PlatformSupport;

  /** Optional user rating information */
  readonly rating?: Rating;
}

/**
 * Defines platform compatibility options for a wishlist item
 * @interface PlatformSupport
 */
export interface PlatformSupport {
  /** Indicates if the item has macOS support */
  readonly hasMacOSVersion?: boolean;

  /** Indicates if the item has Android support */
  readonly hasAndroidVersion?: boolean;

  /** Indicates if the item has iOS support */
  readonly hasIOSVersion?: boolean;
}

/**
 * Groups wishlist items by platform category
 * @interface WishlistGroups
 */
export interface WishlistGroups {
  /** Collection of PC platform wishlist items */
  readonly pc: ReadonlyArray<WishlistItem>;

  /** Collection of console platform wishlist items */
  readonly console: ReadonlyArray<WishlistItem>;

  /** Collection of mobile platform wishlist items */
  readonly mobile: ReadonlyArray<WishlistItem>;
}
