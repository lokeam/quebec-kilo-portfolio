import { PageHeadline } from '@/shared/components/layout/page-headline';
import { PageMain } from '@/shared/components/layout/page-main';
import { PageGrid } from '@/shared/components/layout/page-grid';
import { Skeleton } from '@/shared/components/ui/skeleton';

export function SpendTrackingPageSkeleton() {
  return (
    <PageMain>
      <PageHeadline>
        <Skeleton className="h-10 w-[320px]" />
        <div className="flex items-center space-x-2">
          <Skeleton className="h-10 w-[200px]" />
        </div>
      </PageHeadline>

      <PageGrid>
        {/* TotalMonthlySpendingCard Skeleton */}
        <div className="col-span-full lg:col-span-2">
          <Skeleton className="h-[400px] w-full" />
        </div>

        {/* TotalAnnualSpendingCard Skeleton */}
        <div className="col-span-full lg:col-span-2">
          <Skeleton className="h-[400px] w-full" />
        </div>

        {/* MonthlySpendingAccordion Skeleton */}
        {Array.from({ length: 6 }).map((_, i) => (
            <div key={i} className="col-span-full">
              <Skeleton className="h-[100px] w-full" />
            </div>
        ))}
      </PageGrid>
    </PageMain>
  );
}