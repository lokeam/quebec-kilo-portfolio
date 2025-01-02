import { PageHeadline } from '@/shared/components/layout/page-headline';
import { PageMain } from '@/shared/components/layout/page-main';
import { PageGrid } from '@/shared/components/layout/page-grid';
import { Button } from '@/shared/components/ui/button';
import { Skeleton } from '@/shared/components/ui/skeleton';

export function HomePageContent() {
  return (
    <PageMain>
      <PageHeadline>
        <h1 className='text-2xl font-bold tracking-tight'>Home Page Dashboard</h1>
        <div className='flex items-center space-x-2'>
            <Button>Download Dashboard Summary</Button>
          </div>
      </PageHeadline>

      <PageGrid>
        {/* Statistics Cards Skeletons */}
        {Array.from({ length: 4 }).map((_, i) => (
          <Skeleton key={i} className="h-[120px] w-full" />
        ))}

        {/* Larger Card Skeletons */}
        <Skeleton className="h-[400px] w-full" /> {/* OnlineServices */}
        <Skeleton className="h-[400px] w-full" /> {/* StorageLocations */}
        <Skeleton className="h-[400px] w-full" /> {/* ItemsByPlatform */}
        <Skeleton className="h-[400px] w-full" /> {/* MonthlySpending */}
      </PageGrid>
    </PageMain>
  );
}
