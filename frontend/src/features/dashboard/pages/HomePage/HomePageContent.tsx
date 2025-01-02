import { PageHeadline } from '@/shared/components/layout/page-headline';
import { PageMain } from '@/shared/components/layout/page-main';
import { PageGrid } from '@/shared/components/layout/page-grid';
import { Button } from '@/shared/components/ui/button';
import { Skeleton } from '@/shared/components/ui/skeleton';
import { SingleStatisticsCard } from '@/features/dashboard/organisms/SingleStatisticsCard/SingleStatisticsCard';

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
        <SingleStatisticsCard
          title="Games"
          value={50}
          lastUpdated="six months ago"
          icon="games"
        />
        <SingleStatisticsCard
          title="Monthly Online Services Costs"
          value={120}
          lastUpdated="last month"
          icon="onlineServices"
        />
        <SingleStatisticsCard
          title="Digital Storage Locations"
          value={3}
          lastUpdated="last month"
          icon="coin"
        />
        <SingleStatisticsCard
          title="Physical Storage Locations"
          value={5}
          lastUpdated="three months ago"
          icon="package"
        />

        {/* Larger Card Skeletons */}
        <Skeleton className="h-[400px] w-full" /> {/* OnlineServices */}
        <Skeleton className="h-[400px] w-full" /> {/* StorageLocations */}
        <Skeleton className="h-[400px] w-full" /> {/* ItemsByPlatform */}
        <Skeleton className="h-[400px] w-full" /> {/* MonthlySpending */}
      </PageGrid>
    </PageMain>
  );
}