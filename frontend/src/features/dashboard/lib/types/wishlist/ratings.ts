/**
 * Represents user ratings for a game or item
 * @interface Rating
 */
export interface Rating {
  /** Number of positive ratings/reviews */
  readonly positive: number;

  /** Number of negative ratings/reviews */
  readonly negative: number;

  /** Total number of ratings/reviews (positive + negative) */
  readonly totalReviews: number;
}

/**
 * Calculates the positive rating percentage from a Rating object
 * @param {Rating} rating - The rating object containing positive and total review counts
 * @returns {number} The percentage of positive ratings (0-100)
 * @example
 * const rating = { positive: 80, negative: 20, totalReviews: 100 };
 * const percentage = calculateRatingPercentage(rating); // Returns 80
 */
export const calculateRatingPercentage = (rating: Rating): number => {
  return (rating.positive / rating.totalReviews) * 100;
};