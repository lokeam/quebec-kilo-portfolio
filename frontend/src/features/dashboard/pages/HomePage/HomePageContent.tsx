
// ShadCN UI Components
import { PageHeadline } from '@/shared/components/layout/page-headline';
import { PageMain } from '@/shared/components/layout/page-main';
import { PageGrid } from '@/shared/components/layout/page-grid';
import { Button } from '@/shared/components/ui/button';

// Components
import { ItemsByPlatformCard } from '@/features/dashboard/components/organisms/HomePage/ItemsByPlatformCard/ItemsByPlatformCard';
import { SingleStatisticsCard } from '@/features/dashboard/components/organisms/HomePage/SingleStatisticsCard/SingleStatisticsCard';
import { OnlineServicesCard } from '@/features/dashboard/components/organisms/HomePage/OnlineServicesCard/OnlineServicesCard';
import { StorageLocationsTabCard } from '@/features/dashboard/components/organisms/HomePage/StorageLocationsTabCard/StorageLocationsTabCard';
import { WishListDealsCard } from '@/features/dashboard/components/organisms/HomePage/WishlistDealsCard/WishListDealsCard';
import { MonthlySpendingCard } from '@/features/dashboard/components/organisms/HomePage/MonthlySpendingCard/MonthlySpendingCard';

// Mock Data
import { onlineServicesData } from '@/features/dashboard/components/organisms/HomePage/OnlineServicesCard/onlineServicesCard.mockdata';
import { storageLocationsData } from '@/features/dashboard/components/organisms/HomePage/StorageLocationsTabCard/storageLocationsTabCard.mockdata';
import { itemsByPlatformData } from '@/features/dashboard/components/organisms/HomePage/ItemsByPlatformCard/itemsByPlatformCard.mock.data';
import { wishlistDealsCardMockData } from '@/features/dashboard/components/organisms/HomePage/WishlistDealsCard/wishlistDealsCard.mockdata';
import { monthlySpendingData } from '@/features/dashboard/components/organisms/HomePage/MonthlySpendingCard/monthlySpendingCard.mockdata';

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
          value={72}
          lastUpdated="six months ago"
          icon="games"
          size="sm"
        />
        <SingleStatisticsCard
          title="Monthly Online Services Costs"
          value={120}
          lastUpdated="last month"
          icon="coin"
          size="sm"
        />
        <SingleStatisticsCard
          title="Digital Storage Locations"
          value={3}
          lastUpdated="last month"
          icon="onlineServices"
          size="sm"
        />
        <SingleStatisticsCard
          title="Physical Storage Locations"
          value={5}
          lastUpdated="three months ago"
          icon="package"
          size="sm"
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

        <ItemsByPlatformCard
          domain={itemsByPlatformData.domain}
          totalItemCount={itemsByPlatformData.totalItemCount}
          platformList={itemsByPlatformData.platformList}
          newItemCount={itemsByPlatformData.newItemCount}
        />

        <WishListDealsCard
          starredItem={wishlistDealsCardMockData.starredItem}
          starredItemCurrentPrice={wishlistDealsCardMockData.starredItemCurrentPrice}
          itemsOnSale={wishlistDealsCardMockData.itemsOnSale}
          cheapestSaleItemPercentage={wishlistDealsCardMockData.cheapestSaleItemPercentage}
        />

        <MonthlySpendingCard
          domains={monthlySpendingData.domains}
          spendingByMonth={monthlySpendingData.spendingByMonth}
        />
      </PageGrid>
    </PageMain>
  );
}
