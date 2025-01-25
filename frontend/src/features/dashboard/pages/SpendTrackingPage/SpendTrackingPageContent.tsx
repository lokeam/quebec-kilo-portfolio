
// Template
import { PageHeadline } from '@/shared/components/layout/page-headline';
import { PageMain } from '@/shared/components/layout/page-main';
import { PageGrid } from '@/shared/components/layout/page-grid';
// import { NoResultsFound } from '@/features/dashboard/components/molecules/NoResultsFound';

// Components
import { MonthlySpendingAccordion } from '@/features/dashboard/components/organisms/SpendTrackingPage/MonthlySpendingAccordion/MonthlySpendingAccordion.tsx';
import { TotalMonthlySpendingCard } from '@/features/dashboard/components/organisms/SpendTrackingPage/TotalMonthlySpendingCard/TotalMonthlySpendingCard.tsx';
import { TotalAnnualSpendingCard } from '@/features/dashboard/components/organisms/SpendTrackingPage/TotalAnnualSpendingCard/TotalAnnualSpendingCard.tsx';

// Mock Data
import { spendTrackingPageMockData } from './SpendTrackingPage.mockdata.ts';


export function SpendTrackingPageContent() {
  console.log(`SpendTrackingPageContent: `, spendTrackingPageMockData);


  return (
    <PageMain>
      <PageHeadline>
        <div className='flex items-center'>
          <h1 className='text-2xl font-bold tracking-tight'>Spend Tracking</h1>
        </div>
      </PageHeadline>

      <PageGrid>
        <TotalMonthlySpendingCard />
        <TotalAnnualSpendingCard />
        <MonthlySpendingAccordion
          thisMonth={spendTrackingPageMockData.currentTotalThisMonth}
          future={spendTrackingPageMockData.recurringNextMonth}
          oneTimeTotal={spendTrackingPageMockData.yearlyTotals?.oneTimeTotal}
        />

      </PageGrid>
    </PageMain>
  );
}
