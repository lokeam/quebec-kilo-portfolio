// Template
import { PageHeadline } from '@/shared/components/layout/page-headline';
import { PageMain } from '@/shared/components/layout/page-main';
import { PageGrid } from '@/shared/components/layout/page-grid';
// import { NoResultsFound } from '@/features/dashboard/components/molecules/NoResultsFound';

// Components
import { MonthlySpendingAccordion } from '@/features/dashboard/components/organisms/SpendTrackingPage/MonthlySpendingAccordion/MonthlySpendingAccordion.tsx';
import { TotalMonthlySpendingCard } from '@/features/dashboard/components/organisms/SpendTrackingPage/TotalMonthlySpendingCard/TotalMonthlySpendingCard.tsx';
import { TotalAnnualSpendingCard } from '@/features/dashboard/components/organisms/SpendTrackingPage/TotalAnnualSpendingCard/TotalAnnualSpendingCard.tsx';
import { SpendTrackingPageSkeleton } from '@/features/dashboard/pages/SpendTrackingPage/SpendTrackingPageSkeleton';

// Queries
import { useGetSpendTrackingPageBFFResponse } from '@/core/api/queries/spendTracking.queries';

// Mock data for the page
import { spendTrackingPageMockData } from './SpendTrackingPage.mockdata';

export function SpendTrackingPageContent() {

  const {
    data: bffResponse,
    isLoading,
    error
  } = useGetSpendTrackingPageBFFResponse();

  console.log('🔍 DEBUG: SpendTrackingPage -  bffResponse', bffResponse);


  if (isLoading) {
    return <SpendTrackingPageSkeleton />;
  }

  if (error) {
    console.error('🔍 DEBUG: SpendTrackingPage -  error', error);
  }

  return (
    <PageMain>
      <PageHeadline>
        <div className='flex items-center'>
          <h1 className='text-2xl font-bold tracking-tight'>Spend Tracking</h1>
        </div>
      </PageHeadline>

      <PageGrid>
        <TotalMonthlySpendingCard totalMonthlySpending={bffResponse!.totalMonthlySpending} />
        <TotalAnnualSpendingCard totalAnnualSpending={bffResponse!.totalAnnualSpending} />
        <MonthlySpendingAccordion
          thisMonth={spendTrackingPageMockData.currentTotalThisMonth}
          future={spendTrackingPageMockData.recurringNextMonth}
          oneTimeTotal={spendTrackingPageMockData.yearlyTotals?.oneTimeTotal}
        />
      </PageGrid>
    </PageMain>
  );
}
