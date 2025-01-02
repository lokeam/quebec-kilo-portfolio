import { PageHeadline } from '@/shared/components/layout/page-headline';
import { PageMain } from '@/shared/components/layout/page-main';
import { PageGrid } from '@/shared/components/layout/page-grid';
import { Button } from '@/shared/components/ui/button';
import { Skeleton } from '@/shared/components/ui/skeleton';
import { SingleStatisticsCard } from '@/features/dashboard/organisms/SingleStatisticsCard/SingleStatisticsCard';
import { OnlineServicesCard } from '@/features/dashboard/organisms/OnlineServicesCard/OnlineServicesCard';
import { StorageLocationsTabCard } from '@/features/dashboard/organisms/StorageLocationsTabCard/StorageLocationsTabCard';

// Mock Data
import { onlineServicesData } from '@/features/dashboard/organisms/OnlineServicesCard/onlineServicesCard.mockdata';
import { storageLocationsData } from '@/features/dashboard/organisms/StorageLocationsTabCard/storageLocationsTabCard.mockdata';

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

        {/* Statistics Cards */}
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

        {/* Larger Cards */}
        <OnlineServicesCard
          totalAnnual={onlineServicesData.totalAnnual}
          renewsThisMonth={onlineServicesData.renewsThisMonth}
          totalServices={onlineServicesData.totalServices}
          services={onlineServicesData.services}
        />

        <StorageLocationsTabCard
          totalDigitalLocations={storageLocationsData.totalDigitalLocations}
          totalPhysicalLocations={storageLocationsData.totalPhysicalLocations}
          digitalStorageServices={storageLocationsData.digitalStorage}
          physicalStorageLocations={storageLocationsData.physicalStorage}
        />

        <Skeleton className="h-[400px] w-full" /> {/* ItemsByPlatform */}
        <Skeleton className="h-[400px] w-full" /> {/* MonthlySpending */}
      </PageGrid>
    </PageMain>
  );
}
