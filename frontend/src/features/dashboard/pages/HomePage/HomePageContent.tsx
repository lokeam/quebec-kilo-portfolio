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
import { useEnableIntroToasts } from '@/features/dashboard/hooks/intro-toasts/useEnableIntroToasts';
import { useAuth0 } from '@auth0/auth0-react';

// API Layer hooks
import { useGetDashboardBFFResponse } from '@/core/api/queries/dashboard.queries';

// Skeleton UI
import { HomePageSkeleton } from './HomePageSkeleton';

// Sentry tracking
// import { useSentryTracking } from '@/shared/hooks/useSentryTracking';


export function HomePageContent() {
  const { data: dashboardData, isLoading } = useGetDashboardBFFResponse();
  const { wantsIntroToasts } = useEnableIntroToasts();
  const { user, getIdTokenClaims } = useAuth0();
  // const { trackAction, trackError, trackUserInteraction } = useSentryTracking();

  // Show intro toast for adding games to library
  useShowConditionalIntroToasts(1);

  // Debug function to check Auth0 claims
  const debugAuth0Claims = async () => {
    try {
      const claims = await getIdTokenClaims();
      console.log('🔍 Debug Auth0 Claims:', claims);
      console.log('🔍 User app_metadata:', user?.app_metadata);
      console.log('🔍 Current wantsIntroToasts value:', wantsIntroToasts);
    } catch (error) {
      console.error('Failed to get claims:', error);
    }
  };

  // Debug function to clear localStorage for testing
  const clearToastState = () => {
    localStorage.removeItem('shownIntroToasts');
    console.log('🧹 Cleared shownIntroToasts from localStorage');
    // Reload the page to trigger the toast again
    window.location.reload();
  };

  // Debug function to test optimistic intro toasts
  const testOptimisticIntroToasts = () => {
    console.log('🔍 Current optimistic states:', {
      onboardingComplete: window.__ONBOARDING_OPTIMISTIC_COMPLETE__,
      wantsIntroToasts: window.__WANTS_INTRO_TOASTS__
    });

    // Test setting optimistic state
    window.__WANTS_INTRO_TOASTS__ = true;
    localStorage.setItem('__WANTS_INTRO_TOASTS__', 'true');
    console.log('✅ Set optimistic wantsIntroToasts to true');

    // Reload to test
    window.location.reload();
  };

  // Debug function to reset optimistic state
  const resetOptimisticState = () => {
    delete window.__ONBOARDING_OPTIMISTIC_COMPLETE__;
    delete window.__WANTS_INTRO_TOASTS__;
    localStorage.removeItem('__ONBOARDING_OPTIMISTIC_COMPLETE__');
    localStorage.removeItem('__WANTS_INTRO_TOASTS__');
    console.log('🧹 Reset optimistic state');
    window.location.reload();
  };

  // Debug function to show current state and clear everything
  const debugCurrentState = () => {
    console.log('🔍 Current Debug State:', {
      optimisticOnboardingComplete: window.__ONBOARDING_OPTIMISTIC_COMPLETE__,
      optimisticWantsIntroToasts: window.__WANTS_INTRO_TOASTS__,
      localStorageOnboardingComplete: localStorage.getItem('__ONBOARDING_OPTIMISTIC_COMPLETE__'),
      localStorageWantsIntroToasts: localStorage.getItem('__WANTS_INTRO_TOASTS__'),
      localStorageShownToasts: localStorage.getItem('shownIntroToasts'),
      userAppMetadata: user?.app_metadata,
      wantsIntroToasts: wantsIntroToasts
    });
  };

  // Debug function to test "No toasts" scenario
  const testNoToasts = () => {
    // Clear any existing optimistic state
    delete window.__ONBOARDING_OPTIMISTIC_COMPLETE__;
    delete window.__WANTS_INTRO_TOASTS__;
    localStorage.removeItem('__ONBOARDING_OPTIMISTIC_COMPLETE__');
    localStorage.removeItem('__WANTS_INTRO_TOASTS__');
    localStorage.removeItem('shownIntroToasts');

    // Set optimistic state to "no toasts"
    window.__ONBOARDING_OPTIMISTIC_COMPLETE__ = true;
    window.__WANTS_INTRO_TOASTS__ = false;
    localStorage.setItem('__ONBOARDING_OPTIMISTIC_COMPLETE__', 'true');
    localStorage.setItem('__WANTS_INTRO_TOASTS__', 'false');

    console.log('✅ Set optimistic state to "no toasts"');
    window.location.reload();
  };

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

  console.log('dashboard bff response', dashboardData);

  return (
    <PageMain>
      <PageHeadline>
        <h1 className='text-2xl font-bold tracking-tight'>Dashboard</h1>
        <div className='flex items-center space-x-2'>
          {/* Debug button - remove this after testing */}
          <button
            onClick={debugAuth0Claims}
            className="px-3 py-1 text-sm bg-blue-500 text-white rounded hover:bg-blue-600"
          >
            Debug Auth0 Claims
          </button>
          {/* Debug button to clear toast state */}
          <button
            onClick={clearToastState}
            className="px-3 py-1 text-sm bg-red-500 text-white rounded hover:bg-red-600"
          >
            Clear Toast State
          </button>
          {/* Debug button to test optimistic intro toasts */}
          <button
            onClick={testOptimisticIntroToasts}
            className="px-3 py-1 text-sm bg-green-500 text-white rounded hover:bg-green-600"
          >
            Test Optimistic Toasts
          </button>
          {/* Debug button to reset optimistic state */}
          <button
            onClick={resetOptimisticState}
            className="px-3 py-1 text-sm bg-yellow-500 text-white rounded hover:bg-yellow-600"
          >
            Reset Optimistic State
          </button>
          {/* Debug button to show current state */}
          <button
            onClick={debugCurrentState}
            className="px-3 py-1 text-sm bg-purple-500 text-white rounded hover:bg-purple-600"
          >
            Debug Current State
          </button>
          {/* Debug button to test "No toasts" scenario */}
          <button
            onClick={testNoToasts}
            className="px-3 py-1 text-sm bg-gray-500 text-white rounded hover:bg-gray-600"
          >
            Test No Toasts
          </button>
          {/* <Button>Download Dashboard Summary</Button> */}
        </div>
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

        {/* <WishListDealsCard
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
