import { memo } from 'react';

// Components
import { MemoizedDashboardBadge } from '@/features/dashboard/components/molecules/DashboardBadge/DashboardBadge';
import { SpendingItemYearGrid } from './SpendingItemYearGrid';
import { SpendingItemPaymentDetails } from './SpendingItemPaymentDetails';

// Shadcn UI Components
import { Card, CardContent, CardHeader } from '@/shared/components/ui/card';

// Hooks + Utils
import { LogoOrIcon } from '@/features/dashboard/components/organisms/MediaStoragePage/MediaStoragePageAccordion/LogoOrIcon';
import { useFormattedDate } from '@/features/dashboard/lib/hooks/useFormattedDate';
import { useSpendingData } from '@/features/dashboard/lib/hooks/useSpendingData';

// Icons
import { CreditCard } from 'lucide-react';

// Types
import type { SubscriptionSpend } from '@/features/dashboard/lib/types/spend-tracking/subscription';
import type { OneTimeSpend } from '@/features/dashboard/lib/types/spend-tracking/purchases';
import type { YearlySpending } from '@/features/dashboard/lib/types/spend-tracking/base';

// Guards
import { isSubscriptionSpend } from '@/features/dashboard/lib/types/spend-tracking/guards';

// Constants
import { type PurchasedMediaCategory } from '@/features/dashboard/lib/types/spend-tracking/media';

interface MonthlySpendingItemDetailsProps {
  item: SubscriptionSpend | OneTimeSpend;
  oneTimeTotal: YearlySpending[];
};

export const MonthlySpendingItemDetails = memo(function MonthlySpendingItemDetails({
  item,
  oneTimeTotal = [],
}: MonthlySpendingItemDetailsProps) {
  const { spendingData, title, isSubscription } = useSpendingData(item, oneTimeTotal);

  const dateDisplay = useFormattedDate(
    item.spendTransactionType as PurchasedMediaCategory,
    isSubscriptionSpend(item) ? item.nextBillingDate : '',
    !isSubscriptionSpend(item) ? item.purchaseDate : ''
  );

  return (
    <Card className="bg-[#0A0A0A] text-white border-none mb-4">
      <CardHeader className="space-y-1.5">
        <div className="flex flex-row items-center justify-between space-y-4">

          <div className="flex flex-col">
            <div className="flex flex-row gap-4">
              <MemoizedDashboardBadge
                variant="outline"
                className="bg-purple-900/50 text-purple-300 border-purple-700 w-auto"
                data-testid="media-type-badge"
              >
                {item.mediaType.charAt(0).toUpperCase() + item.mediaType.slice(1)}
              </MemoizedDashboardBadge>
            </div>

            {/* Provider Logo / Item Icon*/}
            <div className="h-14 w-14 flex items-center justify-center my-2">
              <LogoOrIcon name={item.provider?.id ?? ''} mediaType={item.mediaType} />
            </div>
            <h2 className="text-xl font-semibold">{item.title}</h2>
          </div>

          <SpendingItemPaymentDetails
            amount={item.amount}
            date={dateDisplay}
            isSubscription={isSubscription}
          />
        </div>
      </CardHeader>
      <CardContent className="space-y-8">
        {
          isSubscription && (
            <div>
              <h3 className="text-lg font-semibold mb-4">Subscription details</h3>
              <div className="flex xs:flex-col flex-row gap-4">
                <MemoizedDashboardBadge
                  variant="outline"
                  className="bg-blue-900/50 text-blue-300 border-blue-700"
                >
                  {(item as SubscriptionSpend).billingCycle}
                </MemoizedDashboardBadge>
                <MemoizedDashboardBadge
                  variant="outline"
                  className="bg-green-900/50 text-green-300 border-green-700"
                >
                  ${item.amount}
                </MemoizedDashboardBadge>
              </div>
            </div>
        )}

        {/* Yearly Spending */}
        <div>
          <h3 className="text-lg font-semibold mb-4">{title}</h3>
          <SpendingItemYearGrid data={spendingData} />
        </div>

        {/* Payment Method */}
        <div>
          <h3 className="text-lg font-semibold mb-4">Payment method</h3>
          <div className="flex items-center space-x-4">
            <div className="w-12 h-12 bg-gray-800 rounded-lg flex items-center justify-center">
              <CreditCard className="w-6 h-6 text-gray-400" />
            </div>
            <div>
              <div className="font-semibold">{(item as SubscriptionSpend).paymentMethod}</div>
              <div className="text-sm text-gray-400">${item.amount}</div>
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  );
});
