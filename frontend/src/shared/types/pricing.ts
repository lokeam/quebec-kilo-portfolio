/**
 * Represents pricing information for a purchasable item
 * @interface Price
 */
export interface Price {
  /** Original/base price of the item in the store's currency */
  readonly original: number;

  /**
   * Current discounted price, if available
   * @remarks Will be null or undefined if no discount is active
   */
  readonly discounted?: number | null;

  /**
   * Percentage of the current discount
   * @remarks Will be null or undefined if no discount is active
   */
  readonly discountPercentage?: number | null;

  /** Name of the vendor/store offering the item */
  readonly vendor: string;
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
