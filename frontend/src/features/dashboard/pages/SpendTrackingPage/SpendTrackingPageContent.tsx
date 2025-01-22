


// Template
import { PageHeadline } from '@/shared/components/layout/page-headline';
import { PageMain } from '@/shared/components/layout/page-main';
import { PageGrid } from '@/shared/components/layout/page-grid';
// import { NoResultsFound } from '@/features/dashboard/components/molecules/NoResultsFound';

// Components
//import { MonthlySpendCard } from '@/features/dashboard/components/organisms/SpendTrackingPage/MonthlySpendCard';
import { MonthlyRecurringPaymentCard } from '@/features/dashboard/components/organisms/SpendTrackingPage/MonthlyRecurringPaymentCard/MonthlyRecurringPaymentCard';
import { RecurringExpensesAccordion } from '@/features/dashboard/components/organisms/SpendTrackingPage/RecurringExpensesAccordion/RecurringExpensesAccordion';

import { spendTrackingPageMockData } from './SpendTrackingPage.mockdata.ts';

const sampleData = [
  {
    name: "hulu",
    month: "Jan",
    day: "21st",
    title: "Hulu",
    billingCycle: "Monthly",
    amount: "11.99",
  },
  {
    name: "netflix",
    month: "Jan",
    day: "28th",
    title: "Netflix",
    billingCycle: "Monthly",
    amount: "12.99",
  },
  {
    name: "spotify",
    month: "Feb",
    day: "1st",
    title: "Spotify",
    billingCycle: "Monthly",
    amount: "9.99",
  },
]


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
        <RecurringExpensesAccordion
          thisMonth={spendTrackingPageMockData.recurringThisMonth}
          future={spendTrackingPageMockData.recurringNextMonth}
        />

      </PageGrid>
    </PageMain>
  )
}
