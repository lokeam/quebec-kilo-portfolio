// ShadCN UI Components
import { PageHeadline } from '@/shared/components/layout/page-headline';
import { PageMain } from '@/shared/components/layout/page-main';
import { PageGrid } from '@/shared/components/layout/page-grid';

// Components
import { ItemsByPlatformCard } from '@/features/dashboard/components/organisms/HomePage/ItemsByPlatformCard/ItemsByPlatformCard';
import { SingleStatisticsCard } from '@/features/dashboard/components/organisms/HomePage/SingleStatisticsCard/SingleStatisticsCard';
import { OnlineServicesCard } from '@/features/dashboard/components/organisms/HomePage/OnlineServicesCard/OnlineServicesCard';
import { StorageLocationsTabCard } from '@/features/dashboard/components/organisms/HomePage/StorageLocationsTabCard/StorageLocationsTabCard';
import { MonthlySpendingCard } from '@/features/dashboard/components/organisms/HomePage/MonthlySpendingCard/MonthlySpendingCard';

// Hooks
import { useShowConditionalIntroToasts } from '@/features/dashboard/hooks/intro-toasts/useShowConditionalIntroToasts';

// API Layer hooks
import { useGetDashboardBFFResponse } from '@/core/api/queries/dashboard.queries';

// Skeleton UI
import { HomePageSkeleton } from './HomePageSkeleton';

// Sentry tracking
import { useSentryTracking } from '@/shared/hooks/useSentryTracking';
import { useEffect } from 'react';


export function HomePageContent() {
  const { data: dashboardData, isLoading, error } = useGetDashboardBFFResponse();
  const { trackError } = useSentryTracking();

  // Show intro toast for adding games to library
  useShowConditionalIntroToasts(1);

  // Track errors
  useEffect(() => {
    if (error) {
      trackError(error as Error, {
        component: 'HomePageContent',
        action: 'load_dashboard_data',
      });
    }
  }, [error, trackError]);

  if (isLoading) {
    return <HomePageSkeleton />;
  }

  // Add null check for dashboardData
  if (!dashboardData) {
    return (
      <PageMain>
        <PageHeadline>
          <h1 className='text-2xl font-bold tracking-tight'>Dashboard</h1>
        </PageHeadline>
        <div className="text-muted-foreground text-center py-8">
          Unable to load dashboard data
        </div>
      </PageMain>
    );
  }
  // console.log('dashboard bff response', dashboardData);

  return (
    <PageMain>
      <PageHeadline>
        <h1 className='text-2xl font-bold tracking-tight'>Dashboard</h1>
      </PageHeadline>

      <PageGrid>

        {/* Statistics Cards */}
        {/* NOTE:
          We're now properly checking for data existence before using it.
          This prevents crashes when data is null or undefined.
        */}
        <SingleStatisticsCard
          stats={dashboardData.gameStats}
        />
        <SingleStatisticsCard
          stats={dashboardData.subscriptionStats}
        />
        <SingleStatisticsCard
          stats={dashboardData.digitalLocationStats}
        />
        <SingleStatisticsCard
          stats={dashboardData.physicalLocationStats}
        />

        {/* Larger Cards */}
        <OnlineServicesCard
          subscriptionTotal={dashboardData.subscriptionTotal}
          digitalLocations={dashboardData.digitalLocations}
        />

        <StorageLocationsTabCard
          digitalLocations={dashboardData.digitalLocations}
          sublocations={dashboardData.sublocations}
        />

        <ItemsByPlatformCard
          platformList={dashboardData.platformList}
          newItemsThisMonth={dashboardData.newItemsThisMonth}
        />

        {/*
        // Wishlist Card for future use
        <WishListDealsCard
          starredItem={wishlistDealsCardMockData.starredItem}
          starredItemCurrentPrice={wishlistDealsCardMockData.starredItemCurrentPrice}
          itemsOnSale={wishlistDealsCardMockData.itemsOnSale}
          cheapestSaleItemPercentage={wishlistDealsCardMockData.cheapestSaleItemPercentage}
        /> */}

        <MonthlySpendingCard
          mediaTypeDomains={dashboardData.mediaTypeDomains}
          monthlyExpenditures={dashboardData.monthlyExpenditures}
        />

      </PageGrid>
    </PageMain>
  );
}
