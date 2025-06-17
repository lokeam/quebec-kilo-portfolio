
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
//import { onlineServicesData } from '@/features/dashboard/components/organisms/HomePage/OnlineServicesCard/onlineServicesCard.mockdata';
//import { storageLocationsData } from '@/features/dashboard/components/organisms/HomePage/StorageLocationsTabCard/storageLocationsTabCard.mockdata';
//import { itemsByPlatformData } from '@/features/dashboard/components/organisms/HomePage/ItemsByPlatformCard/itemsByPlatformCard.mock.data';
import { wishlistDealsCardMockData } from '@/features/dashboard/components/organisms/HomePage/WishlistDealsCard/wishlistDealsCard.mockdata';
//import { monthlySpendingData } from '@/features/dashboard/components/organisms/HomePage/MonthlySpendingCard/monthlySpendingCard.mockdata';

// Skeleton UI
import { HomePageSkeleton } from './HomePageSkeleton';

// Page mock data
import { homePageMockData } from './Homepage.mockdata';

export function HomePageContent() {
  // Replace with query hook
  const isLoading = false;

  if (isLoading) {
    return <HomePageSkeleton />;
  }

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
          stats={homePageMockData.gameStats}
        />
        <SingleStatisticsCard
          stats={homePageMockData.subscriptionStats}
        />
        <SingleStatisticsCard
          stats={homePageMockData.digitalLocationStats}
        />
        <SingleStatisticsCard
          stats={homePageMockData.physicalLocationStats}
        />

        {/* Larger Cards */}
        <OnlineServicesCard
          subscriptionTotal={homePageMockData.subscriptionTotal}
          subscriptionRecurringNextMonth={homePageMockData.subscriptionRecurringNextMonth}
          digitalLocations={homePageMockData.digitalLocations}
        />

        <StorageLocationsTabCard
          digitalLocations={homePageMockData.digitalLocations}
          sublocations={homePageMockData.sublocations}
        />

        <ItemsByPlatformCard
          totalItemCount={homePageMockData.totalGames}
          platformList={homePageMockData.platformList}
          newItemsThisMonth={homePageMockData.newItemsThisMonth}
        />

        <WishListDealsCard
          starredItem={wishlistDealsCardMockData.starredItem}
          starredItemCurrentPrice={wishlistDealsCardMockData.starredItemCurrentPrice}
          itemsOnSale={wishlistDealsCardMockData.itemsOnSale}
          cheapestSaleItemPercentage={wishlistDealsCardMockData.cheapestSaleItemPercentage}
        />

        <MonthlySpendingCard
          mediaTypeDomains={homePageMockData.mediaTypeDomains}
          monthlyExpenditures={homePageMockData.monthlyExpenditures}
        />
      </PageGrid>
    </PageMain>
  );
}
