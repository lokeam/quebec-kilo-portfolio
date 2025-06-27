import { formatCurrency } from '@/features/dashboard/lib/utils/formatCurrency';
import { Fragment, memo } from 'react';

interface YearlySpendingData {
  year: number;
  amount: number;
}

interface SpendingItemYearGridProps {
  data: YearlySpendingData[];
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
          <div className="font-semibold">{formatCurrency(yearData.amount)}</div>
        </Fragment>
      ))}
    </div>
  );
});
