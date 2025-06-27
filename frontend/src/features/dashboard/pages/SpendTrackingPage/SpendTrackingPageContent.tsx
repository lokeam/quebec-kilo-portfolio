import { useState, useCallback } from 'react';

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
import { SpendTrackingForm } from '@/features/dashboard/components/organisms/SpendTrackingPage/SpendTrackingForm/SpendTrackingForm';

import { DrawerContainer } from '@/features/dashboard/components/templates/DrawerContainer';

// Queries
import { useGetSpendTrackingPageBFFResponse } from '@/core/api/queries/spendTracking.queries';

// Mock data for the page
import { spendTrackingPageMockData } from './SpendTrackingPage.mockdata';

export function SpendTrackingPageContent() {
  const [addOneTimePurchaseOpen, setAddOneTimePurchaseOpen] = useState<boolean>(false);

  /* QUERIES */
  const {
    data: bffResponse,
    isLoading,
    error
  } = useGetSpendTrackingPageBFFResponse();

  console.log('🔍 DEBUG: SpendTrackingPage -  bffResponse', bffResponse);
  /* QUERIES */


  /* HANDLERS */
  const handleFormSuccess = useCallback(() => {
    setAddOneTimePurchaseOpen(false);
  }, []);



  /* RENDER LOGIC */
  if (isLoading) {
    return <SpendTrackingPageSkeleton />;
  }

  if (error) {
    console.error('🔍 DEBUG: SpendTrackingPage -  error', error);
  }
  /* RENDER LOGIC */

  /* DEBUG - DELETE AFTER UI NORMALIZATION */
  const showMockAccordion = false;
  const renderAccordion = () => {
    if (showMockAccordion) {
      return (
        <MonthlySpendingAccordion
          thisMonth={spendTrackingPageMockData.currentTotalThisMonth}
          future={spendTrackingPageMockData.recurringNextMonth}
          oneTimeTotal={spendTrackingPageMockData.yearlyTotals?.oneTimeTotal}
        />
      )
    } else {
      return (
        <MonthlySpendingAccordion
          thisMonth={bffResponse!.currentTotalThisMonth}
          future={bffResponse!.recurringNextMonth}
          oneTimeTotal={bffResponse!.yearlyTotals?.oneTimeTotal}
        />
      )
    }
  };
  /* DEBUG - DELETE AFTER UI NORMALIZATION */

  return (
    <PageMain>
      <PageHeadline>
        <div className='flex items-center'>
          <h1 className='text-2xl font-bold tracking-tight'>Spend Tracking</h1>
        </div>

        <div className="flex items-center space-x-2">
          <DrawerContainer
              open={addOneTimePurchaseOpen}
              onOpenChange={setAddOneTimePurchaseOpen}
              triggerAddLocation="Add a One-Time Purchase"
              title="One-Time Purchase"
              description="Tell us about the one-time purchase you want to add"
              triggerBtnIcon="spend"
            >
              <SpendTrackingForm
                onSuccess={handleFormSuccess}
                buttonText="Add One-Time Purchase"
              />
            </DrawerContainer>
        </div>
      </PageHeadline>

      <PageGrid>
        <TotalMonthlySpendingCard totalMonthlySpending={bffResponse!.totalMonthlySpending} />
        <TotalAnnualSpendingCard totalAnnualSpending={bffResponse!.totalAnnualSpending} />

        {renderAccordion()}

      </PageGrid>
    </PageMain>
  );
}
