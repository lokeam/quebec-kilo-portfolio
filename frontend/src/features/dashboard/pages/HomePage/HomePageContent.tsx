
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
import { MonthlySpendingCard } from '@/features/dashboard/components/organisms/HomePage/MonthlySpendingCard/MonthlySpendingCard';

// Hooks
import { useShowConditionalIntroToasts } from '@/features/dashboard/hooks/intro-toasts/useShowConditionalIntroToasts';

// API Layer hooks
import { useGetDashboardBFFResponse } from '@/core/api/queries/dashboard.queries';

// Skeleton UI
import { HomePageSkeleton } from './HomePageSkeleton';

// Sentry tracking
import { useSentryTracking } from '@/shared/hooks/useSentryTracking';


export function HomePageContent() {
  const { data: dashboardData, isLoading } = useGetDashboardBFFResponse();
  const { trackAction, trackError, trackUserInteraction } = useSentryTracking();

  // Show intro toast for adding games to library
  useShowConditionalIntroToasts(1);


  if (isLoading) {
    return <HomePageSkeleton />;
  }

  console.log('dashboard bff response', dashboardData);

  return (
    <PageMain>
      <PageHeadline>
        <h1 className='text-2xl font-bold tracking-tight'>Home Page Dashboard</h1>
                <div className='flex items-center space-x-2'>
            <Button>Download Dashboard Summary</Button>
                        <Button
              variant="destructive"
              onClick={() => {
                const startTime = performance.now();

                // Track the test action
                trackAction('test_sentry_tunnel', {
                  purpose: 'testing_sentry_implementation',
                  timestamp: new Date().toISOString()
                });

                // Track user interaction with performance
                const duration = performance.now() - startTime;
                trackUserInteraction('test_button_click', duration, {
                  buttonType: 'test',
                  purpose: 'sentry_testing'
                });

                // Throw test error
                const testError = new Error("Sentry Test Error - Testing Phase 4 implementation");
                trackError(testError, {
                  testType: 'manual_test',
                  userInitiated: true,
                  phase: 'phase_4'
                });

                throw testError;
              }}
            >
              Test Sentry Phase 4
            </Button>
          </div>
      </PageHeadline>

      <PageGrid>

        {/* Statistics Cards */}
        {/* NOTE:
          We're already explictly data existence and loading state.
          React Query guarantees that the data grabbed from the API exists when loading is false.
          The non-null assertion employed below is the standard pattern recommended by the React Query docs.
        */}
        <SingleStatisticsCard
          stats={dashboardData!.gameStats}
        />
        <SingleStatisticsCard
          stats={dashboardData!.subscriptionStats}
        />
        <SingleStatisticsCard
          stats={dashboardData!.digitalLocationStats}
        />
        <SingleStatisticsCard
          stats={dashboardData!.physicalLocationStats}
        />

        {/* Larger Cards */}
        <OnlineServicesCard
          subscriptionTotal={dashboardData!.subscriptionTotal}
          digitalLocations={dashboardData!.digitalLocations}
        />

        <StorageLocationsTabCard
          digitalLocations={dashboardData!.digitalLocations}
          sublocations={dashboardData!.sublocations}
        />

        <ItemsByPlatformCard
          platformList={dashboardData!.platformList}
          newItemsThisMonth={dashboardData!.newItemsThisMonth}
        />

        {/* <WishListDealsCard
          starredItem={wishlistDealsCardMockData.starredItem}
          starredItemCurrentPrice={wishlistDealsCardMockData.starredItemCurrentPrice}
          itemsOnSale={wishlistDealsCardMockData.itemsOnSale}
          cheapestSaleItemPercentage={wishlistDealsCardMockData.cheapestSaleItemPercentage}
        /> */}

        <MonthlySpendingCard
          mediaTypeDomains={dashboardData!.mediaTypeDomains}
          monthlyExpenditures={dashboardData!.monthlyExpenditures}
        />
      </PageGrid>
    </PageMain>
  );
}
