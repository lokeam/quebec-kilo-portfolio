import { PageHeadline } from '@/shared/components/layout/page-headline';
import { PageMain } from '@/shared/components/layout/page-main';
import { PageGrid } from '@/shared/components/layout/page-grid';
import { Skeleton } from '@/shared/components/ui/skeleton';

export function HomePageSkeleton() {
  return (
    <PageMain>
      <PageHeadline>
        <Skeleton className="h-10 w-[320px]" />
        <div className="flex items-center space-x-2">
          <Skeleton className="h-10 w-[200px]" />
        </div>
      </PageHeadline>

      <PageGrid>
        {/* Statistics Cards Skeletons */}
        {Array.from({ length: 4 }).map((_, i) => (
          <div key={i} className="md:col-span-1">
            <Skeleton className="h-[120px] w-full" />
          </div>
        ))}

        {/* OnlineServicesCard Skeleton */}
        <div className="col-span-full lg:col-span-2">
          <Skeleton className="h-[400px] w-full" />
        </div>

        {/* StorageLocationsTabCard Skeleton */}
        <div className="col-span-full lg:col-span-2">
          <Skeleton className="h-[400px] w-full" />
        </div>

        {/* ItemsByPlatformCard Skeleton */}
        <div className="flex flex-col">
          <Skeleton className="h-[400px] w-full" />
        </div>

        {/* WishlistDealsCard Skeleton */}
        <div className="flex flex-col gap-6 min-h-[300px] h-full">
          <div className="flex flex-col gap-6 min-h-[300px] h-full">
            <Skeleton className="flex-1 w-full" />
            <Skeleton className="flex-1 w-full" />
          </div>
        </div>

        {/* MonthlySpendingCard Skeleton */}
        <div className="col-span-full lg:col-span-2">
          <Skeleton className="h-[400px] w-full" />
        </div>
      </PageGrid>
    </PageMain>
  );
}