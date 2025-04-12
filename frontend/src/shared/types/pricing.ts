/**
 * Represents pricing information for a game or item
 */
export interface Price {
  /** Original price before any discounts */
  original: number;
  /** Current discounted price */
  discounted: number;
  /** Discount percentage (0-100) */
  discountPercentage: number;
  /** Name of the vendor/store offering the price */
  vendor: string;
}

/**
 * Defines a price range with minimum and maximum values
 * @interface PriceRange
 */
export interface PriceRange {
  /** Minimum price in the range */
  readonly min: number;

  /** Maximum price in the range */
  readonly max: number;
}
