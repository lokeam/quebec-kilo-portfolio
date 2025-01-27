import { memo } from 'react';
import { calculateRatingPercentage } from '@/features/dashboard/lib/types/wishlist/ratings';

interface RatingProps {
  positive: number;
  negative: number;
  totalReviews: number;
}

export const RatingSection = memo(({ positive, negative, totalReviews }: RatingProps) => {
  return (
    <div className="space-y-1">
      <div className="h-2 w-32 bg-gray-700 rounded-full overflow-hidden">
        <div
          className="h-full bg-blue-500 rounded-full"
          style={{ width: `${calculateRatingPercentage({ positive, negative, totalReviews })}%` }}
        />
      </div>
      <div className="flex items-center gap-1 text-xs text-gray-400">
        <span className="text-blue-500">{positive}%</span>
        <span className="text-red-500">{negative}%</span>
        <span className="ml-1">
          {totalReviews.toLocaleString()} User Reviews
        </span>
      </div>
    </div>
  );
}, (prevProps, nextProps) => {
  return (
    prevProps.positive === nextProps.positive &&
    prevProps.negative === nextProps.negative &&
    prevProps.totalReviews === nextProps.totalReviews
  );
});
