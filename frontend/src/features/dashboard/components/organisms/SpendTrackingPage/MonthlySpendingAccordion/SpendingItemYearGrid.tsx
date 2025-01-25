import { Fragment, memo } from 'react';
import type { YearlySpending } from '@/features/dashboard/lib/types/spend-tracking/base';

interface SpendingItemYearGridProps {
  data: YearlySpending[];
  className?: string;
}

export const SpendingItemYearGrid = memo(function SpendingItemYearGrid({
  data,
  className = "space-y-4 sm:space-y-0 sm:grid sm:grid-cols-3 sm:gap-4"
}: SpendingItemYearGridProps) {
  return (
    <div className={className}>
      {data.map((yearData, index) => (
        <Fragment key={`${yearData.year}-${index}`}>
          <div className="text-gray-400">{yearData.year}</div>
        </Fragment>
      ))}
      {data.map((yearData) => (
        <Fragment key={yearData.year}>
          <div className="font-semibold">${yearData.amount}</div>
        </Fragment>
      ))}
    </div>
  );
});
