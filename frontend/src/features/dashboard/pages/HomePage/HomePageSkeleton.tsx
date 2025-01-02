import { PageHeadline } from '@/shared/components/layout/page-headline';
import { PageMain } from '@/shared/components/layout/page-main';
import { PageGrid } from '@/shared/components/layout/page-grid';
import { Skeleton } from '@/shared/components/ui/skeleton';

export function HomePageSkeleton() {
  return (
    <PageMain>
      <PageHeadline>
        <Skeleton className="h-8 w-[250px]" />
        <div className="flex items-center space-x-2">
          <Skeleton className="h-10 w-[200px]" />
        </div>
      </PageHeadline>

      <PageGrid>
        {/* Statistics Cards Skeletons */}
        {Array.from({ length: 4 }).map((_, i) => (
          <Skeleton key={i} className="h-[120px] w-full" />
        ))}

        {/* Larger Card Skeletons */}
        <Skeleton className="h-[300px] w-full" /> {/* OnlineServices */}
        <Skeleton className="h-[400px] w-full" /> {/* StorageLocations */}
        <Skeleton className="h-[400px] w-full" /> {/* ItemsByPlatform */}
        <Skeleton className="h-[400px] w-full" /> {/* MonthlySpending */}
      </PageGrid>
    </PageMain>
  );
}